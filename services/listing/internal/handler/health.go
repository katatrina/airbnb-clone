package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ListingHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
