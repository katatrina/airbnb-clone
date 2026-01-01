package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListProvinces(c *gin.Context) {
	provinces, err := h.listingRepo.ListProvinces(c.Request.Context())
	if err != nil {
		log.Printf("failed to list provinces: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, provinces)
}
