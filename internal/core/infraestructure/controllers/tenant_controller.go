package controllers

import (
	logger "daedalus-engine-go/internal/common/logger"
	usecases "daedalus-engine-go/internal/core/application/use_cases"
	domainErrors "daedalus-engine-go/internal/core/domain/errors"
	httpErrors "daedalus-engine-go/internal/core/infraestructure/http"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TenantController struct {
	uc *usecases.TenantUseCase
}

func NewTenantController(uc *usecases.TenantUseCase) *TenantController {
	return &TenantController{uc: uc}
}

func (ctr *TenantController) CrearTenant(c *gin.Context) {
	log := logger.L()

	log.Info("Iniciando creación de tenant")

	// Parseamos el body
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("Error parsing request", zap.Error(err))
		validationErr := domainErrors.NewValidationError("REQUEST_001", "Formato de request inválido")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	// Ejecutamos caso de uso
	tenant, err := ctr.uc.CreateTenant(c.Request.Context(), req.Name, req.Slug, req.Plan)
	if err != nil {
		log.Error("Error creando tenant", zap.Error(err))
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	// Respuesta
	c.JSON(http.StatusCreated, gin.H{
		"id":     tenant.ID,
		"name":   tenant.Name,
		"slug":   tenant.Slug,
		"plan":   tenant.SubscriptionPlan,
		"status": tenant.Status,
	})
}

func (ctr *TenantController) ObtenerTenantPorId(c *gin.Context) {
	log := logger.L()

	id, _ := c.Params.Get("id")

	// Ejecutamos el caso de uso
	tenant, err := ctr.uc.ObtenerTenantPorId(c, id)

	if err != nil {
		log.Error("Tenant no encontrado")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tenant": tenant,
	})
}

// Obtener tenant por slug
func (ctr *TenantController) ObtenerTenantPorSlug(c *gin.Context) {

}
