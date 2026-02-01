package token

import "time"

type TokenMaker interface {
	CreateToken(userID string) (string, time.Time, error)

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
