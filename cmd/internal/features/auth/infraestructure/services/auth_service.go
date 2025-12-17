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
	now := time.Now()

	expiresAt := now.Add(s.refreshTTL)

	var tenantID *string
	if user.TenantID != nil {
		tenantID = user.TenantID
	}

	return &sessionDomain.Session{
		ID:         uuid.NewString(),
		TenantID:   *tenantID,
		UserID:     user.ID,
		RefreshJWT: refreshHash,
		IssueAt:    &now,
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
	tenantID, _ := claims["tenant_id"].(string)
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

	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"role":      user.Role,
		"email":     user.Email,
		"iat":       now.Unix(),
		"exp":       now.Add(s.accessTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
