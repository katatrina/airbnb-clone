package handler

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/services/user/internal/model"
	"github.com/katatrina/airbnb-clone/services/user/internal/service"
)

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// TODO: Customize raw validation error message
		response.BadRequest(c, response.CodeValidationFailed, err.Error())
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), service.CreateUserParams{
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Password:    req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEmailAlreadyExists):
			response.Conflict(c, response.CodeEmailExists, "Email already exists")
			return
		default:
			// Unknown error = internal server error
			// IMPORTANT: Log the actual error for debugging
			// but don't expose it to the client (security risk)
			log.Printf("[ERROR] Register failed: %v", err)
			response.InternalError(c)
			return
		}
	}

	response.Created(c, UserResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt.Unix(),
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// TODO: Customize raw validation error message
		response.BadRequest(c, response.CodeValidationFailed, err.Error())
		return
	}

	result, err := h.userService.Login(c.Request.Context(), service.LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidCredentials):
			response.Unauthorized(c, "Invalid email or password")
			return

		default:
			log.Printf("[ERROR] Login failed: %v", err)
			response.InternalError(c)
			return
		}
	}

	response.OK(c, LoginResponse{
		AccessToken: result.AccessToken,
	})
}
