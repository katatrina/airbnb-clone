// Package repository handles all database operations.
// This layer is responsible for:
// 1. Executing SQL queries
// 2. Translating database-specific errors into domain errors (sentinel errors)
//
// The key insight: repository NEVER returns raw database errors like pgx.ErrNoRows
// or pgconn.PgError. Instead, it translates them into model.ErrXxx that the
// service and handler layers can understand without knowing about PostgreSQL.
package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

// CreateUser inserts a new user into the database.
// Returns model.ErrEmailAlreadyExists if email violates unique constraint.
func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, display_name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query, user.ID, user.DisplayName, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		// PostgreSQL error code 23505 = unique_violation
		// This happens when email already exists in database
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.ErrEmailAlreadyExists
		}

		// For any other database error, we return it as-is
		// The handler will treat unknown errors as internal server errors
		return err
	}

	return nil
}

// FindUserByEmail finds a user by their email address.
// Returns model.ErrUserNotFound if no user exists with the given email.
// This is used for login - we check if email exists, then compare password.
func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, display_name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	//var user model.User
	//err := r.db.QueryRow(ctx, query, email).Scan(
	//	&user.ID,
	//	&user.DisplayName,
	//	&user.Email,
	//	&user.PasswordHash,
	//	&user.CreatedAt,
	//	&user.UpdatedAt,
	//)
	rows, _ := r.db.Query(ctx, query, email)
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])

	if err != nil {
		// pgx.ErrNoRows means no user found with this email
		// We translate this to our domain error
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, display_name, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	rows, _ := r.db.Query(ctx, query, id)
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
