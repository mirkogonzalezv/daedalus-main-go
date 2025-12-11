package database

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

func (p *PostgresDB) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	return p.pool.Query(ctx, query, args...)
}

func (p *PostgresDB) Exec(ctx context.Context, query string, args ...any) (CommandTag, error) {
	return p.pool.Exec(ctx, query, args...)
}

func (p *PostgresDB) Close() {
	p.pool.Close()
}

func (p *PostgresDB) GetStdlibDB() (*sql.DB, error) {
	config := p.pool.Config()

	// Creamos conexión stdlib desde la configuración de pgx
	connStr := stdlib.RegisterConnConfig(config.ConnConfig)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
