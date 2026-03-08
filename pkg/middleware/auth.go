package middleware

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/pkg/token"
)

// AuthUser contains user info from token.
type AuthUser struct {
	ID string
}

const AuthUserKey = "authUser"

// GetAuthUser .
func GetAuthUser(c *gin.Context) *AuthUser {
	val, exists := c.Get(AuthUserKey)
	if !exists {
		return nil
	}
	return val.(*AuthUser)
}

// MustGetAuthUser .
func MustGetAuthUser(c *gin.Context) *AuthUser {
	val, exists := c.Get(AuthUserKey)
	if !exists {
		panic("MustGetAuthUser: AuthMiddleware not attached to this route")
	}
	return val.(*AuthUser)
}

// AuthMiddleware .
func AuthMiddleware(tokenMaker token.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, response.CodeAuthenticationRequired,
				"Authorization header is required")
			c.Abort()
			return
		}

		scheme, tokenString, found := strings.Cut(authHeader, " ")
		if !found || scheme != "Bearer" || tokenString == "" {
			response.Unauthorized(c, response.CodeAuthenticationRequired,
				"Authorization header format must be: Bearer {token}")
			c.Abort()
			return
		}

		claims, err := tokenMaker.VerifyToken(tokenString)
		if err != nil {
			switch {
			case errors.Is(err, token.ErrTokenExpired):
				response.Unauthorized(c, response.CodeTokenExpired,
					"Token has expired")
			case errors.Is(err, token.ErrTokenInvalid):
				response.Unauthorized(c, response.CodeTokenInvalid,
					"Invalid token")
			default:
				log.Printf("Unexpected token verification error: %v", err)
				response.Unauthorized(c, response.CodeAuthenticationRequired,
					"Authentication failed")
			}
			c.Abort()
			return
		}

		c.Set(AuthUserKey, &AuthUser{ID: claims.UserID})
		c.Next()
	}
}
