package repository

import (
	"context"
	domain "daedalus-engine-go/cmd/internal/features/subscription/domain/entities"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *domain.Subscription) error
	GetById(ctx context.Context, id string) (*domain.Subscription, error)
	UpdateStatus(ctx context.Context, status string) error
	Delete(ctx context.Context, id string) error
}
