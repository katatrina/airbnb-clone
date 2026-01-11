package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListingRepository struct {
	db *pgxpool.Pool
}

func NewListingRepository(db *pgxpool.Pool) *ListingRepository {
	return &ListingRepository{
		db: db,
	}
}
