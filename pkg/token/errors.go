package token

import "errors"

var (
	// ErrTokenInvalid is returned when the token is malformed, has an invalid
	// signature, or cannot be parsed for any reason.
	ErrTokenInvalid = errors.New("token is invalid")

	// ErrTokenExpired is returned when the token was valid but has passed
	// its expiration time. This is separate from ErrTokenInvalid because
	// you might want to handle expired tokens differently (e.g., offer refresh).
	ErrTokenExpired = errors.New("token has expired")
)
