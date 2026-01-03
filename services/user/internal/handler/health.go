package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/services/user/config"
	"github.com/katatrina/airbnb-clone/services/user/internal/repository"
)

type UserHandler struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewUserHandler(userRepo *repository.UserRepository, cfg *config.Config) *UserHandler {
	return &UserHandler{userRepo: userRepo, cfg: cfg}
}

func (h *UserHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
