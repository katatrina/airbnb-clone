package model

import (
	"strings"
	"time"
	"unicode/utf8"
)

type (
	ListingStatus   string
	ListingCurrency string
)

const (
	ListingStatusDraft    ListingStatus = "draft"
	ListingStatusActive   ListingStatus = "active"
	ListingStatusInactive ListingStatus = "inactive"

	ListingCurrencyVND ListingCurrency = "VND"
)

type Listing struct {
	ID            string          `db:"id"`
	HostID        string          `db:"host_id"`
	Title         string          `db:"title"`
	Description   string          `db:"description"`
	PricePerNight int64           `db:"price_per_night"`
	Currency      ListingCurrency `db:"currency"`
	ProvinceCode  int32           `db:"province_code"`
	ProvinceName  string          `db:"province_name"`
	DistrictCode  int32           `db:"district_code"`
	DistrictName  string          `db:"district_name"`
	WardCode      int32           `db:"ward_code"`
	WardName      string          `db:"ward_name"`
	AddressDetail string          `db:"address_detail"`
	Status        ListingStatus   `db:"status"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
	DeletedAt     *time.Time      `db:"deleted_at"`
}

func (l *Listing) ValidateForPublish() error {
	var missing []string

	if strings.TrimSpace(l.Title) == "" || utf8.RuneCountInString(l.Title) < 10 {
		missing = append(missing, "title")
	}

	if l.PricePerNight <= 0 {
		missing = append(missing, "pricePerNight")
	}

	if utf8.RuneCountInString(l.Description) < 50 {
		missing = append(missing, "description")
	}

	if l.ProvinceCode == 0 || l.DistrictCode == 0 || l.WardCode == 0 {
		missing = append(missing, "address")
	}

	if strings.TrimSpace(l.AddressDetail) == "" || utf8.RuneCountInString(l.AddressDetail) < 10 {
		missing = append(missing, "addressDetail")
	}

	if len(missing) > 0 {
		return &IncompleteListingError{MissingFields: missing}
	}

	return nil
}
