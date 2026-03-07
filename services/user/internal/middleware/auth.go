package middleware

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/pkg/token"
	"github.com/katatrina/airbnb-clone/services/user/internal/constant"
)

func AuthMiddleware(tokenMaker token.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, response.CodeAuthenticationRequired, "Authorization header is required")
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
				response.Unauthorized(c, response.CodeTokenExpired, "Token has expired")
			case errors.Is(err, token.ErrTokenInvalid):
				response.Unauthorized(c, response.CodeTokenInvalid, "Invalid token")
			default:
				log.Printf("Unexpected token verification error: %v", err)
				response.Unauthorized(c, response.CodeAuthenticationRequired, "Authentication failed")
			}
			c.Abort()
			return
		}

		c.Set(constant.UserIDKey, claims.UserID)

		c.Next()
	}
}
