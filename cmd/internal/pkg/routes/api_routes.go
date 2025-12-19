package routes

import (
	"daedalus-engine-go/cmd/common/logger"
	"daedalus-engine-go/cmd/container"
	authRoutes "daedalus-engine-go/cmd/internal/features/auth/infraestructure/routes"
	conversationRoutes "daedalus-engine-go/cmd/internal/features/conversation/infraestructure/routes"
	tenantRoutes "daedalus-engine-go/cmd/internal/features/tenants/infraestructure/routes"
	userRoutes "daedalus-engine-go/cmd/internal/features/users/infraestructure/routes"
	authMiddleware "daedalus-engine-go/cmd/internal/pkg/middlewares/auth"
	ipMiddleware "daedalus-engine-go/cmd/internal/pkg/middlewares/ip"
	rateLimitMiddlewares "daedalus-engine-go/cmd/internal/pkg/middlewares/rate_limit"

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
	authGroup.Use(ipMiddleware.IPRateLimitMiddleware(r.container.RateLimitService, log))
	authRoutes.RegisterAuthRoutes(authGroup, r.container.AuthController, r.container.AuthService)

	// Rutas protegidas con auth + rate limiting por userID
	protectedGroup := api.Group("")
	protectedGroup.Use(authMiddleware.AuthMiddleware(r.container.AuthService))
	protectedGroup.Use(rateLimitMiddlewares.RateLimitMiddleware(r.container.RateLimitService, log))

	tenantRoutes.RegisterTenantRoutes(protectedGroup, r.container.TenantController, r.container.AuthService)
	userRoutes.RegisterUserRoutes(protectedGroup, r.container.UserController, r.container.AuthService)
	conversationRoutes.RegisterConversationRoutes(protectedGroup, r.container.ConversationController, r.container.AuthService)
	protectedGroup.POST("/auth/logout-all", authMiddleware.RequiredGlobalAdmin(), r.container.AuthController.LogoutAll)
}
