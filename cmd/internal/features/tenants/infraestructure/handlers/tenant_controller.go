package controllers

import (
	usecases "daedalus-engine-go/cmd/internal/features/tenants/application/use_cases"
	"daedalus-engine-go/cmd/internal/features/tenants/domain/dtos/requests"
	"daedalus-engine-go/cmd/internal/features/tenants/domain/mappers"
	baseErrors "daedalus-engine-go/cmd/internal/pkg/errors"
	httpErrors "daedalus-engine-go/cmd/internal/pkg/http"
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
		validationErr := baseErrors.NewValidationError("REQUEST_001", "Formato de request inválido")
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

func (ctr *TenantController) ObtenerTenantPorSlug(c *gin.Context) {
	slug := c.Param("slug")

	tenant, err := ctr.uc.ObtenerTenantPorSlug(c, slug)

	if err != nil {
		ctr.log.Error("Tenant no encontrado")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	response := mappers.ToGetTenantResponse(tenant)
	c.JSON(http.StatusOK, response)
}

func (ctr *TenantController) ActualizarTenant(c *gin.Context) {
	id := c.Param("id")

	var req requests.UpdateTenantRequest

	if err := c.BindJSON(&req); err != nil {
		ctr.log.Error("Error parsing request", zap.Error(err))
		validationErr := baseErrors.NewValidationError("REQUEST_001", "Formato de request inválido")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	existingTenant, err := ctr.uc.ObtenerTenantPorId(c, id)

	if err != nil {
		ctr.log.Error("Tenant no encontrado")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	if req.Name != "" {
		existingTenant.Name = req.Name
	}

	if req.Slug != "" {
		existingTenant.Slug = req.Slug
	}

	if req.Plan != "" {
		existingTenant.SubscriptionPlan = req.Plan
	}

	err = ctr.uc.ActualizarTenant(c.Request.Context(), existingTenant)
	if err != nil {
		ctr.log.Error("Error actualizando tenant", zap.Error(err))
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	response := mappers.ToUpdateTenantResponse(existingTenant)
	c.JSON(http.StatusOK, response)
}
