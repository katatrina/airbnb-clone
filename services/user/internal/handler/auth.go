package handler

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type User struct {
	ID           string
	DisplayName  string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RegisterRequest struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (r RegisterRequest) Validate() error {
	if r.DisplayName == "" {
		return errors.New("displayName is required")
	}

	if r.Email == "" {
		return errors.New("email is required")
	}

	if !emailRegex.MatchString(r.Email) {
		return errors.New("wrong email format")
	}

	if utf8.RuneCountInString(r.Password) < 8 {
		return errors.New("password is too short (min 8 chars)")
	}

	return nil
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.NewV7()
	if err != nil {
		log.Printf("failed to generate user ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to generate user password hash: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	user := User{
		ID:           userID.String(),
		DisplayName:  req.DisplayName,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}
	err = h.db.QueryRow(c.Request.Context(), "INSERT INTO users (id, display_name, email, password_hash) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at", user.ID, user.DisplayName, user.Email, user.PasswordHash).
		Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}

		log.Printf("failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	resp := RegisterResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}
