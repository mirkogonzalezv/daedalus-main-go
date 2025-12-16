package local

import (
	"context"
	domain "daedalus-engine-go/cmd/core/domain/entities"
	"daedalus-engine-go/cmd/core/domain/repository"
	"database/sql"
)

type PostgresConversationRepository struct {
	db *sql.DB
}

// Create implements repository.ConversationRepository.
func (p *PostgresConversationRepository) Create(ctx context.Context, c *domain.Conversation) error {
	query := `
		INSERT INTO daedalus.conversations(
			id, tenant_id, user_id, title, created_at, updated_at
		) VALUES ($1,$2,$3,$4,NOW(),NOW())
	`

	_, err := p.db.ExecContext(ctx, query, c.ID, c.TenantId, c.UserId, c.Title)

	return err
}

// GetById implements repository.ConversationRepository.
func (p *PostgresConversationRepository) GetById(ctx context.Context, id string) (*domain.Conversation, error) {
	query := `
		SELECT id, tenant_id, user_id, title, created_at, updated_at FROM daedalus.conversations WHERE id = $1
	`

	row := p.db.QueryRowContext(ctx, query, id)

	var c domain.Conversation

	err := row.Scan(
		&c.ID,
		&c.TenantId,
		&c.UserId,
		&c.Title,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != sql.ErrNoRows {
		return nil, nil
	}

	return &c, err

}

// ListByUser implements repository.ConversationRepository.
func (p *PostgresConversationRepository) ListAllConversations(ctx context.Context) ([]*domain.Conversation, error) {
	query := `
		SELECT id, tenant_id, user_id, title, created_at, updated_at
	`

	rows, err := p.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var conversations []*domain.Conversation

	for rows.Next() {
		var c domain.Conversation

		err := rows.Scan(
			&c.ID,
			&c.TenantId,
			&c.UserId,
			&c.Title,
			&c.CreatedAt,
			&c.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		conversations = append(conversations, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

func NewConversationRepository(db *sql.DB) repository.ConversationRepository {
	return &PostgresConversationRepository{db: db}
}
