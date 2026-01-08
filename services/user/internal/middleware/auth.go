// Package middleware contains HTTP middleware functions.
// Middleware runs before/after handlers and can:
// - Authenticate requests
// - Add logging
// - Handle CORS
// - Rate limit
// etc.
package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/pkg/token"
	"github.com/katatrina/airbnb-clone/services/user/internal/constant"
)

// AuthMiddleware validates tokens and extracts user ID.
// It uses TokenMaker interface - doesn't know if it's JWT or Paseto
//
// Benefits of using TokenMaker:
// 1. Middleware doesn't need jwt library import
// 2. Can swap JWT for Paseto without changing middleware
// 3. Easier to test - just mock TokenMaker
//
// Usage:
//
//	tokenMaker := token.NewJWTMaker(secret, expiry)
//	router.GET("/users/me", middleware.AuthMiddleware(tokenMaker), handler.GetMe)
func AuthMiddleware(tokenMaker token.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, response.CodeUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		// Parse "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, response.CodeUnauthorized, "Authorization header format must be: Bearer {token}")
			c.Abort()
			return
		}
		tokenString := parts[1]

		// Verify token using TokenMaker
		// All the JWT complexity is hidden inside tokenMaker.VerifyToken()
		// Middleware doesn't know or care if it's JWT, Paseto, or something else
		claims, err := tokenMaker.VerifyToken(tokenString)
		if err != nil {
			// we can check for specific token errors if we want different handling
			switch {
			case errors.Is(err, token.ErrTokenExpired):
				// Could return a different message or status code for expired tokens
				// e.g, to prompt client to refresh
				response.Unauthorized(c, response.CodeUnauthorized, "Token has expired")
			case errors.Is(err, token.ErrTokenInvalid):
				response.Unauthorized(c, response.CodeUnauthorized, "Invalid token")
			default:
				response.Unauthorized(c, response.CodeUnauthorized, "Authentication failed")
			}
			c.Abort()
			return
		}

		// Set user ID in request context for downstream handlers
		c.Set(constant.UserIDKey, claims.UserID)

		// Continue to the next handler
		c.Next()
	}
}
