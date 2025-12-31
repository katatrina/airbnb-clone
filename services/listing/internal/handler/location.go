package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Province struct {
	Code      string    `json:"code"`
	FullName  string    `json:"fullName"`
	CreatedAt time.Time `json:"createdAt"`
}

func (h *Handler) ListProvinces(c *gin.Context) {
	rows, err := h.db.Query(c.Request.Context(), "SELECT code, full_name, created_at FROM provinces")
	if err != nil {
		log.Printf("failed to list provinces: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer rows.Close()

	var provinces []Province
	for rows.Next() {
		var p Province
		err = rows.Scan(&p.Code, &p.FullName, &p.CreatedAt)
		if err != nil {
			log.Printf("row scan error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		provinces = append(provinces, p)
	}

	if err = rows.Err(); err != nil {
		log.Printf("row iteration error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, provinces)
}
