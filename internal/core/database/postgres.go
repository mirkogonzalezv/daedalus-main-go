package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
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
