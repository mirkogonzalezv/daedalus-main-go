package usescases

import (
	"context"
	conversationRequest "daedalus-engine-go/cmd/internal/features/conversation/domain/dtos/requests"
	conversationResponse "daedalus-engine-go/cmd/internal/features/conversation/domain/dtos/responses"
	conversationEntity "daedalus-engine-go/cmd/internal/features/conversation/domain/entities"
	conversationError "daedalus-engine-go/cmd/internal/features/conversation/domain/errors"
	conversationMapper "daedalus-engine-go/cmd/internal/features/conversation/domain/mappers"
	conversationRepository "daedalus-engine-go/cmd/internal/features/conversation/domain/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ConversationUseCase struct {
	conversationRepo conversationRepository.ConversationRepository
	log              *zap.Logger
}

func NewConversationUseCase(conversationRepo conversationRepository.ConversationRepository, log *zap.Logger) *ConversationUseCase {
	return &ConversationUseCase{
		conversationRepo: conversationRepo,
		log:              log,
	}
}

// CreateConversation - Zero Trust: validar tenant y user
func (uc *ConversationUseCase) CreateConversation(ctx context.Context, req *conversationRequest.CreateConversationRequestDTO, userID, tenantID string) (*conversationResponse.ConversationResponseDTO, error) {
	// Input validation (OWASP)
	if req.Title == "" {
		return nil, conversationError.TitleRequiredError()
	}

	// Zero Trust: validar que userID y tenantID no estén vacíos
	if userID == "" || tenantID == "" {
		uc.log.Warn("Intento de crear conversación sin user/tenant válido")
		return nil, conversationError.UnauthorizedError()
	}

	conversation := &conversationEntity.Conversation{
		ID:        uuid.NewString(),
		TenantId:  tenantID,
		UserId:    userID,
		Title:     req.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Persistir utilizando el repo
	if err := uc.conversationRepo.Create(ctx, conversation); err != nil {
		uc.log.Error("Error al crear conversación", zap.Error(err))
		return nil, conversationError.InternalServerError()
	}

	uc.log.Info("Conversación creada", zap.String("conversation_id", conversation.ID), zap.String("user_id", conversation.UserId))

	return conversationMapper.ToConversationResponse(conversation), nil
}

// GetConversationByUser - Solo conversaciones del usuario autenticado
func (uc *ConversationUseCase) GetConversationByUser(ctx context.Context, userID, tenantID string, page, pageSize int) (*conversationResponse.ConversationListResponseDTO, error) {
	if userID == "" || tenantID == "" {
		return nil, conversationError.UnauthorizedError()
	}

	// Validar paginación
	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	conversations, total, err := uc.conversationRepo.GetByUserId(ctx, userID, tenantID, page, pageSize)
	if err != nil {
		uc.log.Error("Error al obtener conversaciones", zap.Error(err))
		return nil, conversationError.InternalServerError()
	}

	return conversationMapper.ToConversationListResponse(conversations, total, page, pageSize), nil
}

// GetConversationByID - Solo si pertenece al usuario
func (uc *ConversationUseCase) GetConversationByID(ctx context.Context, conversationID, userID, tenantID string) (*conversationResponse.ConversationResponseDTO, error) {
	// Validamos parametros
	if conversationID == "" || userID == "" || tenantID == "" {
		return nil, conversationError.UnauthorizedError()
	}

	conversation, err := uc.conversationRepo.GetById(ctx, conversationID)
	if err != nil {
		uc.log.Error("Error al obtener conversación", zap.Error(err))
		return nil, conversationError.InternalServerError()
	}

	// Verificamos el dueño
	if conversation.UserId != userID || conversation.TenantId != tenantID {
		uc.log.Warn("Intento de acceso no autorizado a conversación", zap.String("conversation_id", conversationID),
			zap.String("user_id", userID))
		return nil, conversationError.ForbiddenError()
	}

	return conversationMapper.ToConversationResponse(conversation), nil
}

// UpdateConversation - Solo el propietario puede actualizar
func (uc *ConversationUseCase) UpdateConversation(ctx context.Context, conversationID string, req *conversationRequest.UpdateConversationRequestDTO, userID, tenantID string) (*conversationResponse.ConversationResponseDTO, error) {
	// Zero Trust: Siempre validar los inputs
	if conversationID == "" || userID == "" || tenantID == "" || req.Title == "" {
		return nil, conversationError.InvalidParametersError()
	}

	conversation, err := uc.conversationRepo.GetById(ctx, conversationID)

	if err != nil {
		uc.log.Error("Error al obtener conversación para actualizar", zap.Error(err))
		return nil, conversationError.InternalServerError()
	}

	if conversation == nil {
		uc.log.Error("Conversación no encontrada")
		return nil, conversationError.NotFoundError()
	}

	if conversation.UserId != userID || conversation.TenantId != tenantID {
		uc.log.Warn("Intento de actualizar no autorizado", zap.String("user_id", userID),
			zap.String("tenant_id", tenantID))
		return nil, conversationError.ForbiddenError()
	}

	// Actualizar
	conversation.Title = req.Title
	conversation.UpdatedAt = time.Now()

	if err := uc.conversationRepo.Update(ctx, conversation); err != nil {
		uc.log.Error("Error al actualizar conversación", zap.Error(err))
		return nil, conversationError.InternalServerError()
	}

	uc.log.Info("Conversación actualizada", zap.String("conversation_id", conversationID))

	return conversationMapper.ToConversationResponse(conversation), nil
}

func (uc *ConversationUseCase) DeleteConversation(ctx context.Context, conversationID, userID, tenantID string) error {
	// Zero Trust: Siempre validar los inputs
	if conversationID == "" || userID == "" || tenantID == "" {
		return conversationError.InvalidParametersError()
	}

	conversation, err := uc.conversationRepo.GetById(ctx, conversationID)

	if err != nil {
		uc.log.Error("Error al obtener conversación para actualizar", zap.Error(err))
		return conversationError.InternalServerError()
	}

	if conversation == nil {
		uc.log.Error("Conversación no encontrada")
		return conversationError.NotFoundError()
	}

	if conversation.UserId != userID || conversation.TenantId != tenantID {
		uc.log.Warn("Intento de eliminación no autorizada",
			zap.String("conversation_id", conversationID),
			zap.String("user_id", userID),
			zap.String("tenant_id", tenantID))
		return conversationError.ForbiddenError()
	}

	if err := uc.conversationRepo.Delete(ctx, conversationID); err != nil {
		uc.log.Error("Error al eliminar conversación", zap.Error(err))
		return conversationError.InternalServerError()
	}

	uc.log.Info("Conversación eliminada", zap.String("conversation_id", conversationID))
	return nil
}
