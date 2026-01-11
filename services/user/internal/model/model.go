package model

import "time"

type User struct {
	ID            string     `repository:"id"`
	DisplayName   string     `repository:"display_name"`
	Email         string     `repository:"email"`
	PasswordHash  string     `repository:"password_hash"`
	EmailVerified bool       `repository:"email_verified"`
	LastLoginAt   *time.Time `repository:"last_login_at"`
	CreatedAt     time.Time  `repository:"created_at"`
	UpdatedAt     time.Time  `repository:"updated_at"`
	DeletedAt     *time.Time `repository:"deleted_at"`
}
