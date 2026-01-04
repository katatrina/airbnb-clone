package model

import "errors"

var (
	// ErrInvalidCredentials is returned when email/password combination is wrong.
	// We use a generic message to avoid leaking information about which field is incorrect.
	ErrInvalidCredentials = errors.New("invalid email or password")

	// ErrUserNotFound is returned when a user cannot be found by ID or other identifier.
	ErrUserNotFound = errors.New("user not found")

	// ErrEmailAlreadyExists is returned when attempting to register with an email
	// that is already in use. This maps to HTTP 409 Conflict.
	ErrEmailAlreadyExists = errors.New("email already exists")
)
