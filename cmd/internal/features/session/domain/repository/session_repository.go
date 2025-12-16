package repository

import (
	"context"
	domain "daedalus-engine-go/cmd/internal/features/session/domain/entities"
)

type SessionRepository interface {
	Create(ctx context.Context, s *domain.Session) error
	GetByRefreshToken(ctx context.Context, refreshHash string) (*domain.Session, error)
	Revoke(ctx context.Context, id string) error
}
