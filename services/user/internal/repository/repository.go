package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katatrina/airbnb-clone/services/user/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Close() {
	r.db.Close()
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, display_name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *
	`

	_, err := r.db.Exec(ctx, query, user.ID, user.DisplayName, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	return err
}
