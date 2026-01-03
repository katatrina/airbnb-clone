package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/services/user/config"
	"github.com/katatrina/airbnb-clone/services/user/internal/service"
)

type UserHandler struct {
	userService *service.UserService
	cfg         *config.Config
}

func NewUserHandler(userService *service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{userService: userService, cfg: cfg}
}

func (h *UserHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
