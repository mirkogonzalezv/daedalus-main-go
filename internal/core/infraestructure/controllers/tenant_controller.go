package controllers

import (
	"daedalus-engine-go/internal/api/dtos/requests"
	"daedalus-engine-go/internal/api/mappers"
	usecases "daedalus-engine-go/internal/core/application/use_cases"
	domainErrors "daedalus-engine-go/internal/core/domain/errors"
	httpErrors "daedalus-engine-go/internal/core/infraestructure/http"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TenantController struct {
	uc  *usecases.TenantUseCase
	log *zap.Logger
}

func NewTenantController(uc *usecases.TenantUseCase, log *zap.Logger) *TenantController {
	return &TenantController{uc: uc, log: log}
}

func (ctr *TenantController) CrearTenant(c *gin.Context) {
	ctr.log.Info("Iniciando creación de tenant")

	// Parseamos el body
	var req requests.CreateTenantRequest
	if err := c.BindJSON(&req); err != nil {
		ctr.log.Error("Error parsing request", zap.Error(err))
		validationErr := domainErrors.NewValidationError("REQUEST_001", "Formato de request inválido")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	// Ejecutamos caso de uso
	tenant, err := ctr.uc.CreateTenant(c.Request.Context(), req.Name, req.Slug, req.Plan)
	if err != nil {
		ctr.log.Error("Error creando tenant", zap.Error(err))
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	// Respuesta
	response := mappers.ToCreateTenantResponse(tenant)
	c.JSON(http.StatusCreated, response)
}

func (ctr *TenantController) ObtenerTenantPorId(c *gin.Context) {

	id := c.Param("id")

	// Ejecutamos el caso de uso
	tenant, err := ctr.uc.ObtenerTenantPorId(c, id)

	if err != nil {
		ctr.log.Error("Tenant no encontrado")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	response := mappers.ToGetTenantResponse(tenant)
	c.JSON(http.StatusOK, response)
}
