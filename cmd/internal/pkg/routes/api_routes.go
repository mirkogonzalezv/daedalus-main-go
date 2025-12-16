package routes

import (
	"daedalus-engine-go/cmd/container"
	tenantRoutes "daedalus-engine-go/cmd/internal/features/tenants/infraestructure/routes"
	userRoutes "daedalus-engine-go/cmd/internal/features/users/infraestructure/routes"

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
	api := router.Group("/api")

	// Rutas de Tenant
	tenantRoutes.RegisterTenantRoutes(api, r.container.TenantController)
	// Aquí van las otras rutas:
	userRoutes.RegisterUserRoutes(api, r.container.UserController)
}
