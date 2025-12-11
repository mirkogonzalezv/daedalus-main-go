package repository

import (
	"context"
	domain "daedalus-engine-go/internal/core/domain/entities"
)

type SessionRepository interface {
	Crete(ctx context.Context, s *domain.Session) error
	GetByRefreshToken(ctx context.Context, refreshHash string) (*domain.Session, error)
	Revoke(ctx context.Context, id string) error
}
