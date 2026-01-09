package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/pkg/response"
)

func (h *UserHandler) Health(c *gin.Context) {
	response.OK(c, gin.H{"status": "ok"})
}
