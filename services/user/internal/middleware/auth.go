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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, response.CodeAuthenticationRequired, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, response.CodeAuthenticationRequired, "Authorization header format must be: Bearer {token}")
			c.Abort()
			return
		}
		tokenString := parts[1]

		claims, err := tokenMaker.VerifyToken(tokenString)
		if err != nil {
			switch {
			case errors.Is(err, token.ErrTokenExpired):
				response.Unauthorized(c, response.CodeTokenExpired, "Token has expired")
			case errors.Is(err, token.ErrTokenInvalid):
				response.Unauthorized(c, response.CodeTokenInvalid, "Invalid token")
			default:
				response.Unauthorized(c, response.CodeAuthenticationRequired, "Authentication failed")
			}
			c.Abort()
			return
		}

		c.Set(constant.UserIDKey, claims.UserID)

		c.Next()
	}
}
