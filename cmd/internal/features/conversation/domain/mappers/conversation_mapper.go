package mappers

import (
	conversationResponseDTO "daedalus-engine-go/cmd/internal/features/conversation/domain/dtos/responses"
	conversationEntity "daedalus-engine-go/cmd/internal/features/conversation/domain/entities"
)

func ToConversationResponse(conversation *conversationEntity.Conversation) *conversationResponseDTO.ConversationResponseDTO {
	return &conversationResponseDTO.ConversationResponseDTO{
		ID:        conversation.ID,
		TenantID:  conversation.TenantId,
		UserID:    conversation.UserId,
		Title:     conversation.Title,
		CreatedAt: conversation.CreatedAt,
		UpdatedAt: conversation.UpdatedAt,
	}
}

func ToConversationListResponse(conversations []*conversationEntity.Conversation, total, page, pageSize int) *conversationResponseDTO.ConversationListResponseDTO {
	conversationResponse := make([]conversationResponseDTO.ConversationResponseDTO, len(conversations))

	for i, conv := range conversations {
		conversationResponse[i] = *ToConversationResponse(conv)
	}

	return &conversationResponseDTO.ConversationListResponseDTO{
		Conversations: conversationResponse,
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
	}
}
