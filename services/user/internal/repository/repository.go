package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) (*UserRepository, error) {
	return &UserRepository{db: db}, nil
}

func (receiver) name() {

}
