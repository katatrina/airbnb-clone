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

func (h *Handler) ListWards(c *gin.Context) {
	provinceCode := c.Query("provinceCode")
	if provinceCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provinceCode is required"})
		return
	}

	wards, err := h.listingRepo.ListWards(c.Request.Context(), provinceCode)
	if err != nil {
		log.Printf("failed to get wards by province code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if len(wards) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "province code not found"})
		return
	}

	c.JSON(http.StatusOK, wards)
}
