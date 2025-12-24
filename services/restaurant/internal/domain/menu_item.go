package domain

import "time"

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
