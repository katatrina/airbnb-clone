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

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (id, display_name, email, password_hash, email_verified, last_login_at, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, display_name, email, password_hash, email_verified, last_login_at, created_at, updated_at, deleted_at
	`

	rows, _ := r.db.Query(ctx, query, user.ID, user.DisplayName, user.Email, user.PasswordHash, user.EmailVerified, user.LastLoginAt, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
	createdUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
				return nil, model.ErrEmailAlreadyExists
			}
		}
		return nil, err
	}

	return &createdUser, nil
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, display_name, email, password_hash, email_verified, last_login_at, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 and deleted_at IS NULL
	`

	rows, _ := r.db.Query(ctx, query, email)
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
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

func (r *UserRepository) UpdateUserLastLogin(ctx context.Context, id string, lastLoginAt time.Time) error {
	query := `
		UPDATE users
		SET last_login_at = $1
		WHERE id = $2 AND deleted_at IS NULL
    `

	_, err := r.db.Exec(ctx, query, lastLoginAt, id)
	// No need to check for affected rows here
	return err
}
