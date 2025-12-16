package local

import (
	"context"
	domain "daedalus-engine-go/cmd/core/domain/entities"
	"daedalus-engine-go/cmd/core/domain/repository"
	"database/sql"
)

type PostgresSubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) repository.SubscriptionRepository {
	return &PostgresSubscriptionRepository{db: db}
}

// Create implements repository.SubscriptionRepository.
func (p *PostgresSubscriptionRepository) Create(ctx context.Context, s *domain.Subscription) error {
	query := `
		INSERT INTO daedalus.subscriptions(
			id, tenant_id, stripe_customer_id, stripe_subscription_id,period_start,period_end,status,created_at, updated_at
		)
			VALUES($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
	`
	_, err := p.db.ExecContext(ctx, query, s.ID, s.TenantID, s.StripeCustomerId, s.StripeSubscriptionID, s.PeriodStart, s.PeriodEnd, s.Status)
	return err
}

// Delete implements repository.SubscriptionRepository.
func (p *PostgresSubscriptionRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM daedalus.subscriptions WHERE id = $1
	`
	_, err := p.db.ExecContext(ctx, query, id)
	return err
}

// GetById implements repository.SubscriptionRepository.
func (p *PostgresSubscriptionRepository) GetById(ctx context.Context, id string) (*domain.Subscription, error) {
	query := `
		SELECT id, tenant_id, stripe_customer_id, stripe_subscription_id, period_start, period_end, status
		FROM daedalus.subscriptions WHERE id = $1
	`

	row := p.db.QueryRowContext(ctx, query, id)

	var s domain.Subscription

	err := row.Scan(
		&s.ID,
		&s.TenantID,
		&s.StripeCustomerId,
		&s.StripeSubscriptionID,
		&s.PeriodStart,
		&s.PeriodEnd,
		&s.Status,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != sql.ErrNoRows {
		return nil, nil
	}

	return &s, err
}

// GetByUser implements repository.SubscriptionRepository.
func (p *PostgresSubscriptionRepository) GetByUser(ctx context.Context, userID string, tenantID string) (*domain.Subscription, error) {
	panic("unimplemented")
}

// Update implements repository.SubscriptionRepository.
func (p *PostgresSubscriptionRepository) Update(ctx context.Context, s *domain.Subscription) error {
	panic("unimplemented")
}
