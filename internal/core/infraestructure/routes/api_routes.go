package routes

import (
	"daedalus-engine-go/internal/core/infraestructure/controllers"

	"github.com/gin-gonic/gin"
)

type APIRouter struct {
	TenantController *controllers.TenantController
}

// Constructor e inyección de dependencia
func NewAPIRouter(
	tenantController *controllers.TenantController,
) *APIRouter {
	return &APIRouter{
		TenantController: tenantController,
	}
}

func (r *APIRouter) RegisterRouter(router *gin.Engine) {
	api := router.Group("/api")

	// Rutas de Tenant
	RegisterTenantRoutes(api, r.TenantController)

	// Aquí van las otras rutas:
}
