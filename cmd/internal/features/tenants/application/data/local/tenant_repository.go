package local

import (
	"context"
	domain "daedalus-engine-go/cmd/internal/features/tenants/domain/entities"
	repository "daedalus-engine-go/cmd/internal/features/tenants/domain/repository"
	"database/sql"
	"errors"
)

type PostgresTenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) repository.TenantRepository {
	return &PostgresTenantRepository{db: db}
}

func (r *PostgresTenantRepository) Create(ctx context.Context, t *domain.Tenant) error {
	// Query para insertar un tenants
	query := `
		INSERT INTO daedalus.tenants(
			id, name, slug, subscription_plan, status, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,NOW(),NOW())
	`
	// Forma para ejecutar una query, y luego pasar parametros para completar la query
	_, err := r.db.ExecContext(ctx, query, t.ID, t.Name, t.Slug, t.SubscriptionPlan, t.Status)
	return err
}

func (r *PostgresTenantRepository) GetById(ctx context.Context, id string) (*domain.Tenant, error) {
	query := `
		SELECT id, name, slug, subscription_plan, status, created_at, updated_at FROM daedalus.tenants WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var t domain.Tenant

	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Slug,
		&t.SubscriptionPlan,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &t, err
}

func (r *PostgresTenantRepository) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {

	query := `
		SELECT id, name, slug, subscription_plan, status, created_at, updated_at FROM daedalus.tenants WHERE slug = $1
	`

	row := r.db.QueryRowContext(ctx, query, slug)

	var t domain.Tenant

	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Slug,
		&t.SubscriptionPlan,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &t, err
}

func (r *PostgresTenantRepository) Update(ctx context.Context, t *domain.Tenant) error {

	query := `
		UPDATE daedalus.tenants
		SET name = $1, slug = $2, subscription_plan = $3, status = $4, updated_at = NOW() WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query, t.Name, t.Slug, t.SubscriptionPlan, t.Status, t.ID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err == nil && rows == 0 {
		return errors.New("tenant not found")
	}

	return err
}
