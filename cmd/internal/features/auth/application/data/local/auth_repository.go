package local

import (
	"context"
	authRepo "daedalus-engine-go/cmd/internal/features/auth/domain/repository"
	domainSession "daedalus-engine-go/cmd/internal/features/session/domain/entities"
	"database/sql"
)

type PostgresAuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) authRepo.AuthRepository {
	return &PostgresAuthRepository{db: db}
}

// CreateSession implements repository.AuthRepository.
func (p *PostgresAuthRepository) CreateSession(ctx context.Context, session *domainSession.Session) error {
	query := `
		INSERT INTO daedalus.sessions(
			id, tenant_id, user_id, refresh_hash, issued_at, expires_at, user_agent, revoked
		)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := p.db.ExecContext(ctx,
		query,
		session.ID,
		session.TenantID,
		session.UserID,
		session.RefreshJWT,
		session.IssueAt,
		session.ExpiresAt,
		session.UserAgent,
		session.Revoked)

	return err
}

// GetSessionByID implements repository.AuthRepository.
func (p *PostgresAuthRepository) GetSessionByID(ctx context.Context, sessionID string) (*domainSession.Session, error) {
	query := `
		SELECT id, tenant_id, user_id, refresh_hash, issued_at, expires_at,
		user_agent, revoked FROM daedalus.sessions WHERE id = &1
	`

	row := p.db.QueryRowContext(ctx, query, sessionID)

	var session domainSession.Session
	err := row.Scan(
		&session.ID,
		&session.TenantID,
		&session.UserID,
		&session.RefreshJWT,
		&session.IssueAt,
		&session.ExpiresAt,
		&session.UserAgent,
		&session.Revoked,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &session, err
}

// GetSessionByRefreshHash implements repository.AuthRepository.
func (p *PostgresAuthRepository) GetSessionByRefreshHash(ctx context.Context, refreshHash string) (*domainSession.Session, error) {
	query := `
		SELECT id, tenant_id, user_id, refresh_hash, issued_at, expires_at,
		user_agent, revoked FROM daedalus.sessions WHERE refresh_hash = &1 AND revoked = false
	`

	row := p.db.QueryRowContext(ctx, query, refreshHash)

	var session domainSession.Session
	err := row.Scan(
		&session.ID,
		&session.TenantID,
		&session.UserID,
		&session.RefreshJWT,
		&session.IssueAt,
		&session.ExpiresAt,
		&session.UserAgent,
		&session.Revoked,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &session, err
}

// RevokeAllUserSessions implements repository.AuthRepository.
func (p *PostgresAuthRepository) RevokeAllUserSessions(ctx context.Context, userID string) error {
	query := `UPDATE daedalus.sessions SET revoked = true WHERE user_id =$1`
	_, err := p.db.ExecContext(ctx, query, userID)
	return err
}

// RevokeSession implements repository.AuthRepository.
func (p *PostgresAuthRepository) RevokeSession(ctx context.Context, sessionID string) error {
	query := `UPDATE daedalus.sessions SET revoked = true WHERE id =$1`
	_, err := p.db.ExecContext(ctx, query, sessionID)
	return err
}

// UpdateSessionLastUsed implements repository.AuthRepository.
func (p *PostgresAuthRepository) UpdateSessionLastUsed(ctx context.Context, sessionID string) error {
	query := `UPDATE daedalus.sessions SET issued_at = NOW() WHERE id =$1`
	_, err := p.db.ExecContext(ctx, query, sessionID)
	return err
}
