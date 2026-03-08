// services/booking/internal/model/errors.go

package model

import "errors"

var (
	ErrBookingNotFound   = errors.New("booking not found")
	ErrBookingNotPending = errors.New("booking must be in pending status")
	ErrNotBookingGuest   = errors.New("user is not the guest of this booking")
	ErrNotBookingHost    = errors.New("user is not the host of this booking")
	ErrSelfBooking       = errors.New("host cannot book their own listing")
	ErrDatesUnavailable  = errors.New("selected dates are not available")
	ErrInvalidDateRange  = errors.New("check-out date must be after check-in date")
	ErrCheckInPast       = errors.New("check-in date cannot be in the past")

	ErrListingNotFound           = errors.New("listing not found or not active")
	ErrListingServiceUnavailable = errors.New("listing service is unavailable")
)
