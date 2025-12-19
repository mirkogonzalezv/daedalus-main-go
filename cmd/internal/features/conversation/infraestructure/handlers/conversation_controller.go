package handlers

import (
	conversationUsesCases "daedalus-engine-go/cmd/internal/features/conversation/application/uses_cases"
	conversationRequestsDTO "daedalus-engine-go/cmd/internal/features/conversation/domain/dtos/requests"
	"net/http"
	"strconv"

	baseErrors "daedalus-engine-go/cmd/internal/pkg/errors"
	httpErrors "daedalus-engine-go/cmd/internal/pkg/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ConversationController struct {
	uc  *conversationUsesCases.ConversationUseCase
	log *zap.Logger
}

func NewConversationController(uc *conversationUsesCases.ConversationUseCase, log *zap.Logger) *ConversationController {
	return &ConversationController{
		uc:  uc,
		log: log,
	}
}

func (ctr *ConversationController) CreateConversation(c *gin.Context) {
	ctr.log.Info("Iniciando creación de conversación")

	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")

	var req conversationRequestsDTO.CreateConversationRequestDTO
	if err := c.BindJSON(&req); err != nil {
		ctr.log.Error("Error parsing request", zap.Error(err))
		validationErr := baseErrors.NewValidationError("REQUEST_001", "Formato de request inválido")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	// Ejecutamos el caso de uso
	conversation, err := ctr.uc.CreateConversation(c.Request.Context(), &req, userID, tenantID)
	if err != nil {
		ctr.log.Error("Error creando conversación", zap.Error(err))
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusCreated, conversation)
}

func (ctr *ConversationController) ObtenerConversaciones(c *gin.Context) {
	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	conversations, err := ctr.uc.GetConversationByUser(c.Request.Context(), userID, tenantID, page, pageSize)
	if err != nil {
		ctr.log.Error("Error obteniendo conversaciones", zap.Error(err))
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusOK, conversations)
}

func (ctr *ConversationController) ObtenerConversacionPorId(c *gin.Context) {
	conversationID := c.Param("id")
	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")

	if conversationID == "" {
		ctr.log.Error("ID de conversación es requerido")
		validationErr := baseErrors.NewValidationError("REQUEST_002", "ID de conversación es requerido")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	conversation, err := ctr.uc.GetConversationByID(c.Request.Context(), conversationID, userID, tenantID)

	if err != nil {
		ctr.log.Error("Conversación no encontrada", zap.Error(err))
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusOK, conversation)
}

func (ctr *ConversationController) ActualizarConversacion(c *gin.Context) {
	conversationID := c.Param("id")
	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")

	if conversationID == "" {
		ctr.log.Error("ID de conversación es requerido")
		validationErr := baseErrors.NewValidationError("REQUEST_002", "ID de conversación es requerido")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	var req conversationRequestsDTO.UpdateConversationRequestDTO
	if err := c.BindJSON(&req); err != nil {
		ctr.log.Error("Error parsing request", zap.Error(err))
		validationErr := baseErrors.NewValidationError("REQUEST_001", "Formato de request inválido")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	conversation, err := ctr.uc.UpdateConversation(c.Request.Context(), conversationID, &req, userID, tenantID)

	if err != nil {
		ctr.log.Error("Error actualizando conversación", zap.Error(err))
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusOK, conversation)
}

func (ctr *ConversationController) EliminarConversacion(c *gin.Context) {
	conversationID := c.Param("id")
	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")

	if conversationID == "" {
		ctr.log.Error("ID de conversación es requerido")
		validationErr := baseErrors.NewValidationError("REQUEST_002", "ID de conversación es requerido")
		statusCode, errorResponse := httpErrors.MapErrorToHttp(validationErr)
		c.JSON(statusCode, errorResponse)
		return
	}

	err := ctr.uc.DeleteConversation(c.Request.Context(), conversationID, userID, tenantID)
	if err != nil {
		ctr.log.Error("Error eliminando conversación", zap.Error(err))
		statusCode, errorResponse := httpErrors.MapErrorToHttp(err)
		c.JSON(statusCode, errorResponse)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
