package local

import (
	"context"
	domain "daedalus-engine-go/cmd/internal/domain/entities"
	"daedalus-engine-go/cmd/internal/domain/repository"
	"database/sql"
)

type PostgresMessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) repository.MessageRepository {
	return &PostgresMessageRepository{db: db}
}

// Create implements repository.MessageRepository.
func (p *PostgresMessageRepository) Create(ctx context.Context, m *domain.Message) error {
	query := `
		INSERT INTO daedalus.messages(
			id, conversation_id, user_id, sender, content, created_at
		) VALUES ($1,$2,$3,$4,$5, NOW())
	`
	_, err := p.db.ExecContext(ctx, query, m.ID, m.ConversationId, m.UserId, m.Sender, m.Content)

	return err
}

// ListByConversation implements repository.MessageRepository.
func (p *PostgresMessageRepository) ListByConversation(ctx context.Context, conversationID string) ([]*domain.Message, error) {
	query := `
		SELECT id, conversation_id, user_id, sender, content,created_at FROM daedalus.messages	
	`

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []*domain.Message

	for rows.Next() {
		var m domain.Message

		err := rows.Scan(
			&m.ID,
			&m.ConversationId,
			&m.UserId,
			&m.Sender,
			&m.Content,
			&m.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
