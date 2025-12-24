package model

import "errors"

// Restaurant errors
var (
	ErrRestaurantNotFound        = errors.New("restaurant not found")
	ErrRestaurantNameRequired    = errors.New("restaurant name is required")
	ErrRestaurantNameTooShort    = errors.New("restaurant name is too short")
	ErrRestaurantNameTooLong     = errors.New("restaurant name is too long")
	ErrRestaurantAddressRequired = errors.New("restaurant address is required")
	ErrRestaurantAddressTooLong  = errors.New("restaurant address is too long")
	ErrRestaurantNotActive       = errors.New("restaurant is not active")
	ErrRestaurantAlreadyExists   = errors.New("restaurant already exists")
)

// MenuItem errors
var (
	ErrMenuItemNotFound           = errors.New("menu item not found")
	ErrMenuItemNameRequired       = errors.New("menu item name is required")
	ErrMenuItemNameTooShort       = errors.New("menu item name is too short")
	ErrMenuItemNameTooLong        = errors.New("menu item name is too long")
	ErrMenuItemPriceInvalid       = errors.New("menu item price must be greater than zero")
	ErrMenuItemDescriptionTooLong = errors.New("menu item description is too long")
	ErrMenuItemRestaurantRequired = errors.New("menu item must belong to a restaurant")
)
