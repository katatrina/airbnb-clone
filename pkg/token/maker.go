// Package token provides functionality for creating and verifying authentication tokens.
// It defines an interface (TokenMaker) that allows different token implementations
// (JWT, Paseto, etc.) to be used interchangeably throughout the application.
//
// This package is intentionally placed in pkg/ (not internal/) because multiple
// services need to verify tokens. For example:
//   - user-service: creates tokens during login
//   - listing-service: verifies tokens to authenticate requests
//   - booking-service: verifies tokens to authenticate requests
package token

import "time"

// TokenMaker is an interface for managing authentication tokens.
type TokenMaker interface {
	// CreateToken generates a new token for the given user ID.
	// The token will be valid for the duration configured in the maker.
	// Returns the token string or an error if token creation fails.
	CreateToken(userID string) error

	// VerifyToken parses and validates a token string.
	// Returns the claims embedded in the token if valid.
	// Returns an error if the token is invalid, expired, or malformed.
	VerifyToken(tokenString string) (*Claims, error)
}

// Claims contains the payload data extracted from a valid token.
//
// We keep this struct minimal and focused. If you need more claims
// (like roles, permissions, etc.), add them here and update the
// TokenMaker implementations accordingly.
type Claims struct {
	// UserID is the unique identifier of the authenticated user.
	// This is extracted from the "sub" (subject) claim in JWT.
	UserID string

	// IssuedAt is when the token was created.
	// Useful for implementing token refresh logic.
	IssuedAt time.Time

	// ExpiresAt is when the token becomes invalid.
	// Services can use this to log out the user.
	ExpiresAt time.Time
}
