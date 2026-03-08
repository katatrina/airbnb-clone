package handler

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/pkg/middleware"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/services/user/internal/model"
)

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.MustGetAuthUser(c).ID

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
			response.NotFound(c, response.CodeUserNotFound, "User not found")
			return
		default:
			log.Printf("[ERROR] failed to get user profile: %s", err)
			response.InternalServerError(c)
			return
		}
	}

	response.OK(c, NewUserResponse(user), "")
}
