package routes

import (
	"daedalus-engine-go/internal/core/infraestructure/controllers"
	"daedalus-engine-go/internal/core/infraestructure/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterTenantRoutes(rg *gin.RouterGroup, ctr *controllers.TenantController) {
	tenant := rg.Group("/tenants")

	// SOLO ADMIN puede crear tenants - se agrega middleware, separados por ,
	tenant.POST("", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.CrearTenant)

	// TODO: Más adelante:
	// obtener tenant por id
	tenant.GET(":id", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.ObtenerTenantPorId)
	// obtener tenant por slug
	tenant.GET("slug/:slug", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.ObtenerTenantPorSlug)
	// actualizar tenant por id
	tenant.PUT(":id", middlewares.AuthMiddleware(), middlewares.RequiredGlobalAdmin(), ctr.ActualizarTenant)
}
