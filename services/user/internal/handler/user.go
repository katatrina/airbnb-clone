package handler

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/services/user/internal/constant"
	"github.com/katatrina/airbnb-clone/services/user/internal/model"
)

// GetMe returns the current authenticated user's profile.
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.MustGet(constant.UserIDKey).(string)

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			// This shouldn't happen in normal flow because:
			// 1. User logged in (JWT was issued)
			// 2. JWT contains valid user ID
			// 3. User should exist
			//
			// But it can happen if:
			// - User was deleted after login
			// - Database was reset
			// - Token was issued for a non-existent user (bug)
			response.NotFound(c, "User not found")
			return
		default:
			log.Printf("[ERROR] GetMe failed for user %s: %s", userID, err)
			response.InternalError(c)
			return
		}
	}

	response.OK(c, UserResponse{
		ID:            user.ID,
		DisplayName:   user.DisplayName,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		LastLoginAt: func() *int64 {
			if user.LastLoginAt == nil {
				return nil
			}
			lastLoginAt := user.LastLoginAt.Unix()
			return &lastLoginAt
		}(),
		CreatedAt: user.CreatedAt.Unix(),
	})
}
