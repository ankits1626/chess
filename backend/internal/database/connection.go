package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps pgxpool for database operations.
type DB struct {
	*Queries
	pool *pgxpool.Pool
}

// NewDB creates database connection pool.
func NewDB(ctx context.Context, dsn string) (*DB, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		Queries: New(pool),
		pool:    pool,
	}, nil
}

// Close closes database connection pool.
func (db *DB) Close() {
	db.pool.Close()
}

// Pool returns underlying connection pool.
func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}
