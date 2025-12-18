package usescases

import (
	"context"
	authRequest "daedalus-engine-go/cmd/internal/features/auth/domain/dtos/requests"
	authDomain "daedalus-engine-go/cmd/internal/features/auth/domain/entities"
	authRepository "daedalus-engine-go/cmd/internal/features/auth/domain/repository"
	authService "daedalus-engine-go/cmd/internal/features/auth/infraestructure/services"
	userRepository "daedalus-engine-go/cmd/internal/features/users/domain/repository"
	"time"

	authErrors "daedalus-engine-go/cmd/internal/features/auth/domain/errors"
	userErrors "daedalus-engine-go/cmd/internal/features/users/domain/errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	authRepo    authRepository.AuthRepository
	userRepo    userRepository.UserRepository
	authService *authService.AuthService
	log         *zap.Logger
}

func NewAuthUseCase(authRepo authRepository.AuthRepository, userRepo userRepository.UserRepository, authService *authService.AuthService, log *zap.Logger) *AuthUseCase {
	return &AuthUseCase{
		authRepo:    authRepo,
		userRepo:    userRepo,
		authService: authService,
		log:         log,
	}
}

// Login - Zero Trust
func (uc *AuthUseCase) Login(ctx context.Context, req *authRequest.LoginRequest) (*authDomain.TokenPair, error) {
	// Input validaciones para OWASP
	if req.Email == "" || req.Password == "" {
		return nil, authErrors.ErrInvalidCredentials()
	}

	// Obtener usuario por email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		uc.log.Warn("Login utilizado con un email inválido", zap.String("email", req.Email), zap.String("ip", req.IP))
		return nil, authErrors.ErrInvalidCredentials()
	}

	// Validar password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		uc.log.Warn("Login con password invalida", zap.String("user_id", user.ID), zap.String("ip", req.IP))
		return nil, authErrors.ErrInvalidCredentials()
	}

	// Check del estado del usuario (Zero Trust)
	if user.Status != "active" {
		uc.log.Warn("Login con usuario inactivo", zap.String("user_id", user.ID))
		return nil, authErrors.ErrUserInactive()
	}

	// Generar tokens
	tokenPair, refreshHash, err := uc.authService.GenerateTokenPair(user)
	if err != nil {
		uc.log.Error("Falló en generar el tokens", zap.Error(err))
		return nil, authErrors.ErrInvalidCredentials()
	}

	// Create session
	session := uc.authService.CreateSession(user, refreshHash, req.IP, req.UserAgent)

	if session == nil {
		return nil, authErrors.ErrAuthInternalServerError()
	}

	if err := uc.authRepo.CreateSession(ctx, session); err != nil {
		uc.log.Error("Failed al crear sesión", zap.Error(err))
		return nil, authErrors.ErrAuthInternalServerError()
	}

	uc.log.Info("Login exitoso", zap.String("user_id", user.ID), zap.String("ip", req.IP))

	return tokenPair, nil
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, req *authRequest.RefreshTokenRequest) (*authDomain.TokenPair, error) {
	// Obtenemos sesión por refresh token
	session, err := uc.authRepo.GetSessionByRefreshHash(ctx, req.RefreshToken)
	if err != nil || session == nil {
		uc.log.Warn("Refresh token invalido", zap.String("session_id", session.ID))
		return nil, authErrors.ErrTokenExpired()
	}

	// Validamos si la sesión a expirado
	if session.ExpiresAt.IsZero() && time.Now().After(session.ExpiresAt) {
		uc.log.Warn("El token refresh ha expirado", zap.String("session_id", session.ID))
		return nil, authErrors.ErrTokenExpired()
	}

	if session.Revoked {
		uc.log.Warn("Refresh token revocado", zap.String("session_id", session.ID))
		return nil, authErrors.ErrTokenRevoked()
	}

	// Obtenemos el usuario
	user, err := uc.userRepo.GetById(ctx, session.UserID)

	if err != nil || user == nil {
		uc.log.Error("Usuario no encontrado para la sesión", zap.String("session_user_id", session.UserID))
		return nil, userErrors.ErrUserNotFound()
	}

	// confirmamos si el usuario esta activo
	if user.Status != "active" {
		uc.log.Warn("Intento de actualización para un usuario inactivo", zap.String("user_id", user.ID))
		return nil, userErrors.ErrUserInactive()
	}

	// Revocamos sesión antigua (revocamos el token)
	if err := uc.authRepo.RevokeSession(ctx, session.ID); err != nil {
		uc.log.Error("Error al revocar la sesión antigua", zap.Error(err))
	}

	// Generamos nuevo token
	tokenPair, refreshHash, err := uc.authService.GenerateTokenPair(user)
	if err != nil {
		uc.log.Error("Fállo al generar nuevo token", zap.Error(err))
		return nil, authErrors.ErrAuthInternalServerError()
	}

	// Creamos nueva sesión
	newSession := uc.authService.CreateSession(user, refreshHash, req.IP, req.UserAgent)
	if err := uc.authRepo.CreateSession(ctx, newSession); err != nil {
		uc.log.Error("Fállo al crear nueva sesión")
		return nil, authErrors.ErrAuthInternalServerError()
	}
	uc.log.Info("Token refreshed", zap.String("user_id", user.ID), zap.String("ip", req.IP))

	return tokenPair, nil
}

// Logout - Revoke session
func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	session, err := uc.authRepo.GetSessionByRefreshHash(ctx, refreshToken)
	if err != nil || session == nil {
		return nil // Ya es invalido, considerar logout
	}

	if err := uc.authRepo.RevokeSession(ctx, session.ID); err != nil {
		uc.log.Error("Error al revocar sesión", zap.Error(err))
		return authErrors.ErrAuthInternalServerError()
	}

	uc.log.Info("Usuario logout", zap.String("user_id", session.UserID))
	return nil
}

// LogoutAll - Revoke all user sessions
func (uc *AuthUseCase) LogoutAll(ctx context.Context, userID string) error {
	if err := uc.authRepo.RevokeAllUserSessions(ctx, userID); err != nil {
		uc.log.Error("Error al revocar todas las sesiones", zap.Error(err))
		return authErrors.ErrAuthInternalServerError()
	}
	uc.log.Info("Todas las sesiones revocadas", zap.String("user_id", userID))
	return nil
}
