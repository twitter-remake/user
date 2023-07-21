package clients

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgreSQLClient creates a new PostgreSQL pgxpool client
func NewPostgreSQLClient(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
