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
//	    if errors.Is(err, token.ErrTokenExpired) {
//	        // Token was valid but has expired - maybe offer refresh
//	    }
//	    if errors.Is(err, token.ErrTokenInvalid) {
//	        // Token is malformed or signature is wrong
//	    }
//	}
var (
	// ErrTokenInvalid is returned when the token is malformed, has an invalid
	// signature, or cannot be parsed for any reason.
	ErrTokenInvalid = errors.New("token is invalid")

	// ErrTokenExpired is returned when the token was valid but has passed
	// its expiration time. This is separate from ErrTokenInvalid because
	// you might want to handle expired tokens differently (e.g., offer refresh).
	ErrTokenExpired = errors.New("token has expired")
)
