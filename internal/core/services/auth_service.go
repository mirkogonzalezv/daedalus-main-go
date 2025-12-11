package services

import (
	"context"
	domain "daedalus-engine-go/internal/core/domain/entities"
)

type RegisterDTO struct {
	TenantSlug string
	Tenant     string

	Name     string
	Email    string
	Password string
	Role     string // owner/admin/user
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         *domain.User
	Tenant       *domain.Tenant
}

type AuthService interface {
	Register(ctx context.Context, dto RegisterDTO) (*AuthResult, error)
	Login(ctx context.Context, tenantSlug, email, password string) (*AuthResult, error)
	Refresh(ctx context.Context, refreshToken string) (*AuthResult, error)
	Logout(ctx context.Context, refreshToken string) error
}
