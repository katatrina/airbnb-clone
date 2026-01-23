package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/pkg/response"
)

func (h *ListingHandler) Health(c *gin.Context) {
	response.OK(c, gin.H{"status": "ok"}, "Listing Service Operational")
}
