package local

import (
	"context"
	domain "daedalus-engine-go/cmd/internal/features/conversation/domain/entities"
	"daedalus-engine-go/cmd/internal/features/conversation/domain/repository"
	"database/sql"
)

type PostgresConversationRepository struct {
	db *sql.DB
}

// Delete implements [repository.ConversationRepository].
func (p *PostgresConversationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM daedalus.conversations WHERE id = $1`
	_, err := p.db.ExecContext(ctx, query, id)
	return err
}

// GetByUserId implements [repository.ConversationRepository].
func (p *PostgresConversationRepository) GetByUserId(ctx context.Context, userId, tenantId string, page int, pageSize int) ([]*domain.Conversation, int, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT id, tenant_id, user_id, title, created_at, updated_at
		FROM daedalus.conversations WHERE user_id = $1 AND tenant_id = $2
		ORDER BY updated_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := p.db.QueryContext(ctx, query, userId, tenantId, pageSize, offset)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}

		conversations = append(conversations, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// Query para obtener total count
	countQuery := `
		SELECT COUNT(*)
		FROM daedalus.conversations
		WHERE user_id = $1 AND tenant_id = $2
	`

	var total int
	err = p.db.QueryRowContext(ctx, countQuery, userId, tenantId).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return conversations, total, nil

}

// Update implements [repository.ConversationRepository].
func (p *PostgresConversationRepository) Update(ctx context.Context, c *domain.Conversation) error {
	query := `
		UPDATE daedalus.conversations SET title = $2, updated_at = $3 WHERE id =$1
	`

	_, err := p.db.ExecContext(ctx, query, c.ID, c.Title, c.UpdatedAt)
	return err
}

// Create implements repository.ConversationRepository.
func (p *PostgresConversationRepository) Create(ctx context.Context, c *domain.Conversation) error {
	query := `
		INSERT INTO daedalus.conversations(
			id, tenant_id, user_id, title, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := p.db.ExecContext(ctx, query, c.ID, c.TenantId, c.UserId, c.Title, c.CreatedAt, c.UpdatedAt)

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

	if err == sql.ErrNoRows {
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
