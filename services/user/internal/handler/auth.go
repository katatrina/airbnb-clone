package handler

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/pkg/request"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/services/user/internal/model"
)

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := request.ShouldBindJSON(c, &req); err != nil {
		response.HandleJSONBindingError(c, err)
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), model.CreateUserParams{
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Password:    req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEmailAlreadyExists):
			response.Conflict(c, response.CodeEmailAlreadyExists, "Email already exists")
			return
		default:
			log.Printf("[ERROR] Register failed: %v", err)
			response.InternalServerError(c)
			return
		}
	}

	response.Created(c, UserResponse{
		ID:            user.ID,
		DisplayName:   user.DisplayName,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt.Unix(),
	}, "User registered successfully")
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := request.ShouldBindJSON(c, &req); err != nil {
		response.HandleJSONBindingError(c, err)
		return
	}

	result, err := h.userService.LoginUser(c.Request.Context(), model.LoginUserParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrIncorrectCredentials):
			response.Unauthorized(c, response.CodeCredentialsInvalid, "Incorrect email or password")
			return

		default:
			log.Printf("[ERROR] fail to login user: %v", err)
			response.InternalServerError(c)
			return
		}
	}

	response.OK(c, LoginResponse{
		AccessToken: result.AccessToken,
	}, "User login successfully")
}
