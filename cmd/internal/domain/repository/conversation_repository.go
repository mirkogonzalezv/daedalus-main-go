package repository

import (
	"context"
	domain "daedalus-engine-go/cmd/internal/domain/entities"
)

type ConversationRepository interface {
	Create(ctx context.Context, c *domain.Conversation) error
	GetById(ctx context.Context, id string) (*domain.Conversation, error)
	ListAllConversations(ctx context.Context) ([]*domain.Conversation, error)
}
