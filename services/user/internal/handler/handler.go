package handler

import "github.com/katatrina/airbnb-clone/services/user/internal/service"

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
