package handlers

import (
	authUseCase "daedalus-engine-go/cmd/internal/features/auth/application/uses_cases"
	authRequest "daedalus-engine-go/cmd/internal/features/auth/domain/dtos/requests"
	"net/http"

	baseErrors "daedalus-engine-go/cmd/internal/pkg/errors"
	httpErrors "daedalus-engine-go/cmd/internal/pkg/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthController struct {
	authUseCase *authUseCase.AuthUseCase
	log         *zap.Logger
}

func NewAuthController(authUseCase *authUseCase.AuthUseCase, log *zap.Logger) *AuthController {
	return &AuthController{
		authUseCase: authUseCase,
		log:         log,
	}
}

// Login endpoint
func (ctrl *AuthController) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email", binding:"required,email"`
		Password string `json:"password", binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.log.Warn("Request login invalido", zap.Error(err))
		validationErr := baseErrors.NewValidationError("AUTH_001", "Request con formato inválido")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	// Obtenemos info del cliente por Zero Trust
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	loginReq := &authRequest.LoginRequest{
		Email:     req.Email,
		Password:  req.Password,
		IP:        clientIP,
		UserAgent: userAgent,
	}

	tokenPair, err := ctrl.authUseCase.Login(c.Request.Context(), loginReq)

	if err != nil {
		ctrl.log.Error("Login falló", zap.Error(err))
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	// Set refresh token como httpOnly cookie (Security best practice)
	c.SetCookie("refresh_token", tokenPair.RefreshToken, 7*24*3600, "/", "", true, true)

	// Retornamos solo el access token en la respuesta
	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenPair.AccessToken,
		"expires_in":   tokenPair.ExpiresIn,
		"token_type":   "Bearer",
	})
}

// Refresh Token endpoint
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		ctrl.log.Warn("Falta refresh token cookie")
		validationErr := baseErrors.NewValidationError("AUTH_002", "Se requiere token refresh")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	// Obtenemos info del cliente por Zero Trust
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	refreshReq := &authRequest.RefreshTokenRequest{
		RefreshToken: refreshToken,
		IP:           clientIP,
		UserAgent:    userAgent,
	}

	tokenPair, err := ctrl.authUseCase.RefreshToken(c.Request.Context(), refreshReq)
	if err != nil {
		ctrl.log.Error("Fállo en el Token refresh", zap.Error(err))
		// Limpiamos cookie invalido
		c.SetCookie("refresh_token", "", -1, "/", "", true, true)
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.SetCookie(
		"refresh_token",
		tokenPair.RefreshToken,
		7*24*3600,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenPair.AccessToken,
		"expires_in":   tokenPair.ExpiresIn,
		"token_type":   "Bearer",
	})
}

// Logout endpoint
func (ctrl *AuthController) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Logged out exitoso"})
		return
	}

	if err := ctrl.authUseCase.Logout(c.Request.Context(), refreshToken); err != nil {
		ctrl.log.Error("Logout failed", zap.Error(err))
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out exitoso",
	})
}

// Logout All endpoint
func (ctrl *AuthController) LogoutAll(c *gin.Context) {
	userID, exist := c.Get("user_id")
	if !exist {
		validationErr := baseErrors.NewValidationError("AUTH_003", "Usuario no autenticado")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	if err := ctrl.authUseCase.LogoutAll(c.Request.Context(), userID.(string)); err != nil {
		ctrl.log.Error("Logout all falló", zap.Error(err))
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	// Limpiamos el actual refresh_token tookie
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Todas las sesiones terminadas",
	})
}
