package controllers

import (
	requests "daedalus-engine-go/cmd/api/dtos/requests"
	"daedalus-engine-go/cmd/api/mappers"
	usecases "daedalus-engine-go/cmd/core/application/use_cases"
	domainErrors "daedalus-engine-go/cmd/core/domain/errors"
	httpErrors "daedalus-engine-go/cmd/core/infraestructure/http"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	uc  *usecases.UserUseCase
	log *zap.Logger
}

func NewUserController(uc *usecases.UserUseCase, log *zap.Logger) *UserController {
	return &UserController{uc: uc, log: log}
}

func (ctr *UserController) CrearUsuario(c *gin.Context) {
	ctr.log.Info("Iniciando creación del usuario")

	var req requests.CreateUserRequestDTO
	if err := c.BindJSON(&req); err != nil {
		ctr.log.Error("Error parsing request", zap.Error(err))
		validationErr := domainErrors.NewValidationError("REQUEST_001", "Formato de request inválido")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	// Ejecutamos el caso de uso
	usuario, err := ctr.uc.CreateUser(c.Request.Context(), req.Name, req.TenantID, req.Email, req.Password, req.Role)
	if err != nil {
		ctr.log.Error("Error creando usuario", zap.Error(err))
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	response := mappers.ToCreateUserResponse(usuario)
	c.JSON(http.StatusOK, response)
}

func (ctr *UserController) ObtenerUsuarioPorId(c *gin.Context) {
	id := c.Param("id")

	usuario, err := ctr.uc.ObtenerUsuarioPorId(c, id)

	if err != nil {
		ctr.log.Error("Usuario no encontrado")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	response := mappers.ToGetUserResponse(usuario)
	c.JSON(http.StatusOK, response)
}

func (ctr *UserController) ObtenerUsuarioPorEmailYTenant(c *gin.Context) {
	email := c.Query("email")        // devuelve "" si no existe
	tenantID := c.Query("tenant_id") // devuelve "" si no existe

	if email == "" {
		ctr.log.Error("Email es obligatorio")
		validationErr := domainErrors.NewValidationError("REQUEST_002", "Email es requerido para la busqueda")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	usuario, err := ctr.uc.ObtenerUsuarioPorEmailYTenant(c, tenantID, email)
	if err != nil {
		ctr.log.Error("Usuario no encontrado")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	response := mappers.ToGetUserResponse(usuario)
	c.JSON(http.StatusOK, response)
}

func (ctr *UserController) EliminarUsuario(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		ctr.log.Error("Id es requerido")
		validationErr := domainErrors.NewValidationError("REQUEST_003", "Id es requerido")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	err := ctr.uc.EliminarUsuario(c, id)
	if err != nil {
		ctr.log.Error("Problemas al eliminar usuario")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}
}

// # Funciones de ROOT
func (ctr *UserController) ObtenerUsuarioPorEmail(c *gin.Context) {
	email := c.Query("email")

	usuario, err := ctr.uc.ObtenerUsuarioPorEmail(c, email)
	if err != nil {
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	response := mappers.ToGetUserResponse(usuario)
	c.JSON(http.StatusOK, response)
}
