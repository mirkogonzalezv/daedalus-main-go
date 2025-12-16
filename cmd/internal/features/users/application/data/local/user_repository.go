package local

import (
	"context"
	domain "daedalus-engine-go/cmd/internal/features/users/domain/entities"
	repository "daedalus-engine-go/cmd/internal/features/users/domain/repository"
	"database/sql"
	"errors"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO daedalus.users(
			id, tenant_id, name, email, password_hash, role, status, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
	`

	_, err := r.db.ExecContext(ctx, query, u.ID, u.TenantID, u.Name, u.Email, u.PasswordHash, u.Role, u.Status)
	return err
}

func (r *PostgresUserRepository) GetById(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, name, email, password_hash, role, status, created_at, updated_at FROM daedalus.users WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var u domain.User

	err := row.Scan(
		&u.ID,
		&u.TenantID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *PostgresUserRepository) GetByEmailAndTenant(ctx context.Context, tenantID string, email string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, name, email, password_hash, role, status, created_at, updated_at FROM daedalus.users WHERE tenant_id = $1 AND email =$2
	`

	row := r.db.QueryRowContext(ctx, query, tenantID, email)

	var u domain.User

	err := row.Scan(
		&u.ID,
		&u.TenantID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *PostgresUserRepository) Update(ctx context.Context, u *domain.User) error {
	query := `
		UPDATE daedalus.users SET name = $1, email= $2, role = $3, status = $4, updated_at = NOW() WHERE id = $5 
	`

	result, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.Role, u.Status, u.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err == nil && rows == 0 {
		return errors.New("user not found")
	}

	return err
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM daedalus.users WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

/*
# Funciones ROOT
*/
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, name, email, password_hash, role, status, created_at, updated_at FROM daedalus.users WHERE email = $1
	`

	// QueryRowContext:  es utilizado para devolver solo 1 linea
	row := r.db.QueryRowContext(ctx, query, email)

	var u domain.User

	err := row.Scan(
		&u.ID,
		&u.TenantID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *PostgresUserRepository) GetRootUsers(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT id, tenant_id, name, email, password_hash, role, status, created_at, updated_at FROM daedalus.users
	`

	// QueryContext: se utiliza cuando esperamos multiples filas de resultado.
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close() //Cerramos rows despues de terminar de usarlo

	var users []*domain.User

	for rows.Next() {
		var u domain.User

		// Se escane cada fila de la estructura de usuario
		err := rows.Scan(
			&u.ID,
			&u.TenantID,
			&u.Name,
			&u.Email,
			&u.PasswordHash,
			&u.Role,
			&u.Status,
			&u.CreatedAt,
			&u.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PostgresUserRepository) GetUsersByRole(ctx context.Context, role string) ([]*domain.User, error) {
	query := `
		SELECT id, tenant_id, name, email, password_hash, role, status, created_at, updated_at FROM daedalus.users WHERE role = $1
	`

	rows, err := r.db.QueryContext(ctx, query, role)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		var u domain.User

		err := rows.Scan(
			&u.ID,
			&u.TenantID,
			&u.Name,
			&u.Email,
			&u.PasswordHash,
			&u.Role,
			&u.Status,
			&u.CreatedAt,
			&u.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
