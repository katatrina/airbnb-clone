package service

import (
	"context"

	"github.com/katatrina/food-delivery/services/restaurant/internal/model"
)

const (
	MinRestaurantNameLength = 2
	MaxRestaurantNameLength = 100
)

type Store interface {
	CreateRestaurant(ctx context.Context, restaurant *model.Restaurant) error
	// GetRestaurantByID(ctx context.Context, id string) (*model.Restaurant, error)
}

type RestaurantService struct {
	store Store
}

func NewRestaurantService(store Store) *RestaurantService {
	return &RestaurantService{store: store}
}
