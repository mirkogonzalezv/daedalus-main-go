package repository

import (
	"context"
	sessionDomain "daedalus-engine-go/cmd/internal/features/session/domain/entities"
)

type AuthRepository interface {
	// Gestión de la sesión
	CreateSession(ctx context.Context, session *sessionDomain.Session) error
	GetSessionByRefreshHash(ctx context.Context, refreshHash string) (*sessionDomain.Session, error)
	GetSessionByID(ctx context.Context, sessionID string) (*sessionDomain.Session, error)
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
	UpdateSessionLastUsed(ctx context.Context, sessionID string) error
}
