package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/katatrina/airbnb-clone/services/user/internal/constant"
)

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.MustGet(constant.UserIDKey).(string)

	var user User
	err := h.db.QueryRow(c.Request.Context(), "SELECT id, email, display_name, created_at, updated_at FROM users WHERE id=$1", userID).
		Scan(&user.ID, &user.Email, &user.DisplayName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("userID %s not found", userID)})
			return
		}

		log.Printf("failed to get user profile by id: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
