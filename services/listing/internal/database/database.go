package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err = dbPool.Ping(ctx); err != nil {
		return nil, err
	}

	return dbPool, nil
}
