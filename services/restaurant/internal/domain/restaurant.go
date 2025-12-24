package domain

import (
	"time"
)

// Restaurant represents a restaurant in the system.
type Restaurant struct {
	ID        string
	Name      string
	Address   string
	Phone     string
	Email     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
