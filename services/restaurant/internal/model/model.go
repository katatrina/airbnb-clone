package model

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

// MenuItem represents a menu item belonging to a restaurant.
type MenuItem struct {
	ID           string
	RestaurantID string
	Name         string
	Description  string
	Price        float64
	Category     string
	IsAvailable  bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
