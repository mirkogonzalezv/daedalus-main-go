package local

import (
	"context"
	domain "daedalus-engine-go/cmd/core/domain/entities"
	"daedalus-engine-go/cmd/core/domain/repository"
	"database/sql"
	"errors"
)

type PostgresSessionRepository struct {
	db *sql.DB
}

// Create implements repository.SessionRepository.
func (p *PostgresSessionRepository) Create(ctx context.Context, s *domain.Session) error {
	query := `
		INSERT INTO daedalus.sessions(id, tenant_id, user_id, refresh_hash, issued_at, expires_at, ip, user_agent, revoked)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`
	_, err := p.db.ExecContext(ctx, query, s.ID, s.TenantID, s.UserID, s.RefreshJWT, s.IssueAt, s.ExpiresAt, s.IP, s.UserAgent, s.Revoked)
	return err
}

// GetByRefreshToken implements repository.SessionRepository.
func (p *PostgresSessionRepository) GetByRefreshToken(ctx context.Context, refreshHash string) (*domain.Session, error) {
	query := `
		SELECT id, tenant_id, user_id, refresh_hash, issued_at, expires_at, ip, user_agent, revoked FROM daedalus.sessions WHERE refresh_hash = $1
	`
	row := p.db.QueryRowContext(ctx, query, refreshHash)

	var s domain.Session

	err := row.Scan(
		&s.ID,
		&s.TenantID,
		&s.UserID,
		&s.RefreshJWT,
		&s.IssueAt,
		&s.ExpiresAt,
		&s.IP,
		&s.UserAgent,
		&s.Revoked,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &s, err
}

// Revoke implements repository.SessionRepository.
func (p *PostgresSessionRepository) Revoke(ctx context.Context, id string) error {
	query := `
		UPDATE daedalus.sessions SET revoked = false WHERE id = $1
	`

	result, err := p.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err == nil && rows == 0 {
		return errors.New("sessions not found")
	}

	return err
}

func NewSessionRepository(db *sql.DB) repository.SessionRepository {
	return &PostgresSessionRepository{db: db}
}
