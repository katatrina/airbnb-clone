package model

import (
	"time"
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
	ProvinceCode  string          `db:"province_code"`
	ProvinceName  string          `db:"province_name"`
	WardCode      string          `db:"ward_code"`
	WardName      string          `db:"ward_name"`
	AddressDetail string          `db:"address_detail"`
	Status        ListingStatus   `db:"status"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
	DeletedAt     *time.Time      `db:"deleted_at"`
}

type Province struct {
	Code      string    `db:"code"`
	FullName  string    `db:"full_name"`
	CreatedAt time.Time `db:"created_at"`
}

type Ward struct {
	Code         string    `db:"code" `
	FullName     string    `db:"full_name" `
	ProvinceCode string    `db:"province_code"`
	CreatedAt    time.Time `db:"created_at"`
}
