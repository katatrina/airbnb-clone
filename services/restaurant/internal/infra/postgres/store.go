package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katatrina/food-delivery/services/restaurant/internal/infra/postgres/sqlc"
)

type Store struct {
	db *sqlc.Queries
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{db: sqlc.New(pool)}
}
