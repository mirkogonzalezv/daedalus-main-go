package repository

import (
	"context"
	domain "daedalus-engine-go/cmd/core/domain/entities"
)

type MessageRepository interface {
	Create(ctx context.Context, m *domain.Message) error
	ListByConversation(ctx context.Context, conversationID string) ([]*domain.Message, error)
}
