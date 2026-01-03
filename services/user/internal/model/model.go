package model

import "time"

type User struct {
	ID           string    `db:"id"`
	DisplayName  string    `db:"display_name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
