package repository

import (
	"context"
	domain "daedalus-engine-go/cmd/internal/domain/entities"
)

type TenantRepository interface {
	Create(ctx context.Context, t *domain.Tenant) error
	GetById(ctx context.Context, id string) (*domain.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error)
	Update(ctx context.Context, t *domain.Tenant) error
}
