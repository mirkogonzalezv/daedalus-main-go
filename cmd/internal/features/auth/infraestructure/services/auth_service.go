package services

import (
	"crypto/rand"
	authDomain "daedalus-engine-go/cmd/internal/features/auth/domain/entities"
	sessionDomain "daedalus-engine-go/cmd/internal/features/session/domain/entities"
	userDomain "daedalus-engine-go/cmd/internal/features/users/domain/entities"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthService struct {
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
	log        *zap.Logger
}

func NewAuthService(jwtSecret string, accessTTL, refreshTTL time.Duration, log *zap.Logger) *AuthService {
	return &AuthService{
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL, // 7 días
		log:        log,
	}
}

func (s *AuthService) GenerateTokenPair(user *userDomain.User) (*authDomain.TokenPair, string, error) {
	// Gerenamos refresh token hash
	refreshHash, err := s.generateSecureHash()
	if err != nil {
		return nil, "", err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, "", err
	}

	return &authDomain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshHash,
		ExpiresIn:    int64(s.accessTTL.Seconds()),
	}, refreshHash, nil
}

func (s *AuthService) CreateSession(user *userDomain.User, refreshHash, ip, userAgent string) *sessionDomain.Session {

	// Zero Trust: Validamos que solo los usuarios globales pueden tener sesiones sin tenant
	if user.TenantID == nil && user.Role != "root" && user.Role != "system_admin" {
		s.log.Error("Intento no autorizado de crear una sesión global",
			zap.String("user_id", user.ID),
			zap.String("role", user.Role),
			zap.String("email", user.Email),
		)

		return nil
	}

	// audit log para sesiones globales
	if user.TenantID == nil {
		s.log.Info("Global session created",
			zap.String("user_id", user.ID),
			zap.String("role", user.Role),
			zap.String("ip", ip))
	}

	now := time.Now()
	expiresAt := now.Add(s.refreshTTL)

	return &sessionDomain.Session{
		ID:         uuid.NewString(),
		TenantID:   user.TenantID,
		UserID:     user.ID,
		RefreshJWT: refreshHash,
		IssuedAt:   &now,
		ExpiresAt:  expiresAt,
		IP:         &ip,
		UserAgent:  &userAgent,
		Revoked:    false,
	}
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*authDomain.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validando firma method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	// Obtenemos los claims seguros
	userID, _ := claims["user_id"].(string)

	// Manejo seguro de tenant-id nullable
	var tenantID string
	if tid, ok := claims["tenant_id"]; ok && tid != nil {
		tenantID, _ = tid.(string)
	}

	role, _ := claims["role"].(string)
	email, _ := claims["email"].(string)
	iat, _ := claims["iat"].(float64)
	exp, _ := claims["exp"].(float64)

	return &authDomain.Claims{
		UserID:    userID,
		TenantID:  tenantID,
		Role:      role,
		Email:     email,
		IssuedAt:  int64(iat),
		ExpiresAt: int64(exp),
	}, nil
}

// Funciones privadas
func (s *AuthService) generateSecureHash() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) generateAccessToken(user *userDomain.User) (string, error) {
	now := time.Now()

	// Manejo seguro del tenant_id null para JWT
	var tenantIDClaim interface{}
	if user.TenantID != nil {
		tenantIDClaim = *user.TenantID
	} else {
		tenantIDClaim = nil
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": tenantIDClaim,
		"role":      user.Role,
		"email":     user.Email,
		"iat":       now.Unix(),
		"exp":       now.Add(s.accessTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
