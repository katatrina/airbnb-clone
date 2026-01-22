// Package repository handles all database operations.
// This layer is responsible for:
// 1. Executing SQL queries
// 2. Translating database-specific errors into domain errors (sentinel errors)
//
// The key insight: repository NEVER returns raw database errors like pgx.ErrNoRows
// or pgconn.PgError (except internal database error). Instead, it translates them into model.ErrXxx that the
// service and handler layers can understand without knowing about PostgreSQL.
package repository

import (
	"context"
	"errors"
	"time"

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

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, display_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, query, user.ID, user.DisplayName, user.Email, user.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "users_email_key" {
			return model.ErrEmailAlreadyExists
		}

		return err
	}

	return nil
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, display_name, email, password_hash, email_verified, last_login_at, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 and deleted_at IS NULL
	`

	//var user model.User
	//err := r.repository.QueryRow(ctx, query, email).Scan(
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

func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, display_name, email, password_hash, email_verified, last_login_at, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 and deleted_at IS NULL
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

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string, lastLoginAt *time.Time) error {
	query := `
		UPDATE users
		SET last_login_at = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
    `

	_, err := r.db.Exec(ctx, query, lastLoginAt, id)
	// No need to check for affected rows here
	return err
}
