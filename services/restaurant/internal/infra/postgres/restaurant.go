package postgres

import (
	"context"

	"github.com/katatrina/food-delivery/services/restaurant/internal/infra/postgres/sqlc"
	"github.com/katatrina/food-delivery/services/restaurant/internal/model"
)

func (s *Store) CreateRestaurant(ctx context.Context, restaurant *model.Restaurant) error {
	arg := sqlc.CreateRestaurantParams{
		ID:        restaurant.ID,
		Name:      restaurant.Name,
		Address:   restaurant.Address,
		Phone:     restaurant.Phone,
		Email:     restaurant.Email,
		IsActive:  restaurant.IsActive,
		CreatedAt: restaurant.CreatedAt,
		UpdatedAt: restaurant.UpdatedAt,
	}

	return s.db.CreateRestaurant(ctx, arg)
}
