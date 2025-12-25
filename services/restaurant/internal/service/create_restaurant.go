package service

import (
	"context"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/katatrina/food-delivery/services/restaurant/internal/model"
)

type CreateRestaurantCmd struct {
	Name    string
	Address string
	Phone   string
	Email   string
}

func (c *CreateRestaurantCmd) Validate() error {
	length := utf8.RuneCountInString(c.Name)

	if length == 0 {
		return model.ErrRestaurantNameRequired
	}

	if length < MinRestaurantNameLength {
		return model.ErrRestaurantNameTooShort
	}

	if length > MaxRestaurantNameLength {
		return model.ErrRestaurantNameTooLong
	}

	if c.Address == "" {
		return model.ErrRestaurantAddressRequired
	}

	return nil
}

func (s *RestaurantService) CreateRestaurant(ctx context.Context, cmd CreateRestaurantCmd) (*model.Restaurant, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	restaurantID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	restaurant := model.Restaurant{
		ID:        restaurantID.String(),
		Name:      cmd.Name,
		Address:   cmd.Address,
		Phone:     cmd.Phone,
		Email:     cmd.Email,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.store.CreateRestaurant(ctx, &restaurant); err != nil {
		return nil, err
	}

	return &restaurant, nil
}
