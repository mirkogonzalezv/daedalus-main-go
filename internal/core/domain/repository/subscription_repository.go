package repository

import (
	"context"
	domain "daedalus-engine-go/internal/core/domain/entities"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *domain.Subscription) error
	GetById(ctx context.Context, id string) (*domain.Subscription, error)
	GetByUser(ctx context.Context, userID string, tenantID string) (*domain.Subscription, error)
	Update(ctx context.Context, s *domain.Subscription) error
	Delete(ctx context.Context, id string) error
}
