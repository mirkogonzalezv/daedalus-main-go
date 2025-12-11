package repository

import (
	"context"
	domain "daedalus-engine-go/internal/core/domain/entities"
)

type ConversationRepository interface {
	Create(ctx context.Context, c *domain.Conversation) error
	GetById(ctx context.Context, id string) (*domain.Conversation, error)
	ListByUser(ctx context.Context, userID string, tenantID string) ([]*domain.Conversation, error)
}
