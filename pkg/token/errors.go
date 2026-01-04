package token

import "errors"

// Sentinel errors for token operations.
// These errors are used by TokenMaker implementations and can be checked
// using errors.Is() in handlers/middleware.
//
// Example usage in middleware:
//
//	claims, err := tokenMaker.VerifyToken(tokenString)
//	if err != nil {
//	    if errors.Is(err, token.ErrExpiredToken) {
//	        // Token was valid but has expired - maybe offer refresh
//	    }
//	    if errors.Is(err, token.ErrInvalidToken) {
//	        // Token is malformed or signature is wrong
//	    }
//	}
var (
	// ErrInvalidToken is returned when the token is malformed, has an invalid
	// signature, or cannot be parsed for any reason.
	ErrInvalidToken = errors.New("token is invalid")

	// ErrExpiredToken is returned when the token was valid but has passed
	// its expiration time. This is separate from ErrInvalidToken because
	// you might want to handle expired tokens differently (e.g., offer refresh).
	ErrExpiredToken = errors.New("token has expired")
)
