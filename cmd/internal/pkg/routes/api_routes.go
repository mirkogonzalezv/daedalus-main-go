package routes

import (
	"daedalus-engine-go/cmd/common/logger"
	"daedalus-engine-go/cmd/container"
	authRoutes "daedalus-engine-go/cmd/internal/features/auth/infraestructure/routes"
	tenantRoutes "daedalus-engine-go/cmd/internal/features/tenants/infraestructure/routes"
	userRoutes "daedalus-engine-go/cmd/internal/features/users/infraestructure/routes"
	"daedalus-engine-go/cmd/internal/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

type APIRouter struct {
	container *container.Container
}

// Constructor e inyección de dependencia
func NewAPIRouter(cont *container.Container) *APIRouter {
	return &APIRouter{
		container: cont,
	}
}

func (r *APIRouter) RegisterRouter(router *gin.Engine) {
	log := logger.L()
	api := router.Group("/api")

	// Rutas de auth con rate limit por IP (auth_routes ya maneja el prefijo /auth)
	authGroup := api.Group("")
	authGroup.Use(middlewares.IPRateLimitMiddleware(r.container.RateLimitService, log))
	authRoutes.RegisterAuthRoutes(authGroup, r.container.AuthController, r.container.AuthService)

	// Rutas protegidas con auth + rate limiting por userID
	protectedGroup := api.Group("")
	protectedGroup.Use(middlewares.AuthMiddleware(r.container.AuthService))
	protectedGroup.Use(middlewares.RateLimitMiddleware(r.container.RateLimitService, log))

	tenantRoutes.RegisterTenantRoutes(protectedGroup, r.container.TenantController, r.container.AuthService)
	userRoutes.RegisterUserRoutes(protectedGroup, r.container.UserController, r.container.AuthService)

	protectedGroup.POST("/auth/logout-all", middlewares.RequiredGlobalAdmin(), r.container.AuthController.LogoutAll)
}
