// Package model contains domain models and sentinel errors.
// Sentinel errors are pre-defined errors that can be compared using errors.Is().
// This is idiomatic Go - simple, explicit, and easy to understand.
package model

import "time"

// User represents a user in the system.
// This is the domain entity - it maps directly to the database table.
// Notice we use `db` tags for database mapping, NOT `json` tags.
// JSON serialization is handled by DTOs in the handler layer.
type User struct {
	ID           string    `db:"id"`
	DisplayName  string    `db:"display_name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
