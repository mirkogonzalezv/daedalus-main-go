package repository

import (
	"context"
	domain "daedalus-engine-go/cmd/internal/features/message/domain/entities"
)

type MessageRepository interface {
	Create(ctx context.Context, m *domain.Message) error
	ListByConversation(ctx context.Context, conversationID string) ([]*domain.Message, error)
}
