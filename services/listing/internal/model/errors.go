package model

import "errors"

var (
	ErrListingNotFound      = errors.New("listing not found")
	ErrProvinceCodeNotFound = errors.New("province code not found")
	ErrWardCodeNotFound     = errors.New("ward code not found")
	ErrWardProvinceMismatch = errors.New("ward does not belong to the selected province")
)
