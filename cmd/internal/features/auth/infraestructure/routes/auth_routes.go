package routes

import (
	authControllers "daedalus-engine-go/cmd/internal/features/auth/infraestructure/handlers"
	authService "daedalus-engine-go/cmd/internal/features/auth/infraestructure/services"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, ctrl *authControllers.AuthController, authService *authService.AuthService) {
	auth := rg.Group("/auth")

	// Endpoints publicos (no requiere autenticación)
	auth.POST("/login", ctrl.Login)
	auth.POST("/refresh", ctrl.RefreshToken)
	auth.POST("/logout", ctrl.Logout)
}
