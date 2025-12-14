package repository

import (
	"context"
	domain "daedalus-engine-go/cmd/core/domain/entities"
)

// Contratos que define todo lo que puede hacer con la entidad de Usuario
type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetById(ctx context.Context, id string) (*domain.User, error)
	GetByEmailAndTenant(ctx context.Context, tenantID string, email string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id string) error
}
