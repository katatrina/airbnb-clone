package model

import "errors"

var (
	ErrListingNotFound              = errors.New("listing not found")
	ErrListingOwnerMismatch         = errors.New("listing does not belong to user")
	ErrActiveListingCannotBeUpdated = errors.New("active listing cannot be updated")

	ErrProvinceCodeNotFound     = errors.New("province code not found")
	ErrDistrictCodeNotFound     = errors.New("district code not found")
	ErrWardCodeNotFound         = errors.New("ward code not found")
	ErrDistrictProvinceMismatch = errors.New("district does not belong to the selected province")
	ErrWardDistrictMismatch     = errors.New("ward does not belong to the selected district")
)
