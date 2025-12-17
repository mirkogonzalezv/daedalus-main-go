package routes

import (
	authService "daedalus-engine-go/cmd/internal/features/auth/infraestructure/services"
	controllers "daedalus-engine-go/cmd/internal/features/tenants/infraestructure/handlers"
	middlewares "daedalus-engine-go/cmd/internal/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterTenantRoutes(rg *gin.RouterGroup, ctr *controllers.TenantController, authService *authService.AuthService) {
	tenant := rg.Group("/tenants")

	// SOLO ADMIN puede crear tenants - se agrega middleware, separados por ,
	tenant.POST("", middlewares.AuthMiddleware(authService), middlewares.RequiredGlobalAdmin(), ctr.CrearTenant)

	// TODO: Más adelante:
	// obtener tenant por id
	tenant.GET(":id", middlewares.AuthMiddleware(authService), middlewares.RequiredGlobalAdmin(), ctr.ObtenerTenantPorId)
	// obtener tenant por slug
	tenant.GET("slug/:slug", middlewares.AuthMiddleware(authService), middlewares.RequiredGlobalAdmin(), ctr.ObtenerTenantPorSlug)
	// actualizar tenant por id
	tenant.PUT(":id", middlewares.AuthMiddleware(authService), middlewares.RequiredGlobalAdmin(), ctr.ActualizarTenant)
}
