package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrListingNotFound      = errors.New("listing not found")
	ErrListingOwnerMismatch = errors.New("listing does not belong to user")

	ErrListingNotDraft    = errors.New("listing must be in draft status")
	ErrListingNotActive   = errors.New("listing must be in active status")
	ErrListingNotInactive = errors.New("listing must be in inactive status")

	ErrActiveListingCannotBeUpdated = errors.New("active listing cannot be updated")
	ErrListingHasActiveBookings     = errors.New("listing has active bookings")

	ErrProvinceCodeNotFound     = errors.New("province code not found")
	ErrDistrictCodeNotFound     = errors.New("district code not found")
	ErrWardCodeNotFound         = errors.New("ward code not found")
	ErrDistrictProvinceMismatch = errors.New("district does not belong to the selected province")
	ErrWardDistrictMismatch     = errors.New("ward does not belong to the selected district")
)

type IncompleteListingError struct {
	MissingFields []string
}

func (e *IncompleteListingError) Error() string {
	return fmt.Sprintf("Listing incomplete: %s", strings.Join(e.MissingFields, ", "))
}
