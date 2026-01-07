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
	ID            string        `db:"id" json:"id"`
	HostID        string        `db:"host_id" json:"hostId"`
	Title         string        `db:"title" json:"title"`
	Description   *string       `db:"description" json:"description"`
	PricePerNight int64         `db:"price_per_night" json:"pricePerNight"`
	Currency      string        `db:"currency" json:"currency"`
	ProvinceCode  string        `db:"province_code" json:"provinceCode"`
	ProvinceName  string        `db:"province_name" json:"provinceName"`
	WardCode      string        `db:"ward_code" json:"wardCode"`
	WardName      string        `db:"ward_name" json:"wardName"`
	AddressDetail string        `db:"address_detail" json:"addressDetail"`
	Status        ListingStatus `db:"status" json:"status"`
	CreatedAt     time.Time     `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time     `db:"updated_at" json:"-"`
	DeletedAt     *time.Time    `db:"deleted_at" json:"-"`
}

type Province struct {
	Code      string    `db:"code"`
	FullName  string    `db:"full_name"`
	CreatedAt time.Time `db:"created_at"`
}

type Ward struct {
	Code         string    `db:"code" json:"code"`
	FullName     string    `db:"full_name" json:"fullName"`
	ProvinceCode string    `db:"province_code" json:"provinceCode"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
}
