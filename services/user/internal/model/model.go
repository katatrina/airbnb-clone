package model

import "time"

type User struct {
	ID           string    `db:"id" json:"id"`
	DisplayName  string    `db:"display_name" json:"displayName"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
}
