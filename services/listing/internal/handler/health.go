package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/services/listing/config"
	"github.com/katatrina/airbnb-clone/services/listing/internal/db"
)

type Handler struct {
	listingRepo *db.ListingRepository
	cfg         *config.Config
}

func NewHandler(listingRepo *db.ListingRepository, cfg *config.Config) *Handler {
	return &Handler{
		listingRepo: listingRepo,
		cfg:         cfg,
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
