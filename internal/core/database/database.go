package database

import (
	"context"
)

// Interfaz para luego hacer las implementaciones respectivas
type Database interface {
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	Exec(ctx context.Context, query string, args ...any) (CommandTag, error)
	Close()
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close()
}

type CommandTag interface{}
