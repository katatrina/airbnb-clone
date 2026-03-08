package repository

import "github.com/jackc/pgx/v5/pgxpool"

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db}
}
