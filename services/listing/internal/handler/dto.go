package handler

import "github.com/katatrina/airbnb-clone/services/listing/internal/model"

type CreateListingRequest struct {
	Title         string `json:"title" validate:"required,min=10,max=200" normalize:"trim,singlespace"`
	Description   string `json:"description" validate:"omitempty,max=2000" normalize:"trim,singlespace"`
	PricePerNight int64  `json:"pricePerNight" validate:"required,gte=1"`
	ProvinceCode  string `json:"provinceCode" validate:"required" normalize:"trim"`
	DistrictCode  string `json:"districtCode" validate:"required" normalize:"trim"`
	WardCode      string `json:"wardCode" validate:"required" normalize:"trim"`
	AddressDetail string `json:"addressDetail" validate:"required,min=10,max=500" normalize:"trim,singlespace"`
}

type ListingResponse struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	PricePerNight int64  `json:"pricePerNight"`
	Currency      string `json:"currency"`
	ProvinceCode  string `json:"provinceCode"`
	ProvinceName  string `json:"provinceName"`
	DistrictCode  string `json:"districtCode"`
	DistrictName  string `json:"districtName"`
	WardCode      string `json:"wardCode"`
	WardName      string `json:"wardName"`
	AddressDetail string `json:"addressDetail"`
	Status        string `json:"status"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
}

type ProvinceResponse struct {
	Code     string `json:"code"`
	FullName string `json:"fullName"`
}

type DistrictResponse struct {
	Code         string `json:"code"`
	FullName     string `json:"fullName"`
	ProvinceCode string `json:"provinceCode"`
}

type WardResponse struct {
	Code         string `json:"code"`
	FullName     string `json:"fullName"`
	DistrictCode string `json:"districtCode"`
}

func NewListingResponse(listing *model.Listing) *ListingResponse {
	return &ListingResponse{
		ID:            listing.ID,
		Title:         listing.Title,
		Description:   listing.Description,
		PricePerNight: listing.PricePerNight,
		Currency:      string(listing.Currency),
		ProvinceCode:  listing.ProvinceCode,
		ProvinceName:  listing.ProvinceName,
		DistrictCode:  listing.DistrictCode,
		DistrictName:  listing.DistrictName,
		WardCode:      listing.WardCode,
		WardName:      listing.WardName,
		AddressDetail: listing.AddressDetail,
		Status:        string(listing.Status),
		CreatedAt:     listing.CreatedAt.Unix(),
		UpdatedAt:     listing.UpdatedAt.Unix(),
	}
}

func NewListingsResponse(listings []model.Listing) []ListingResponse {
	resp := make([]ListingResponse, len(listings))
	for i, listing := range listings {
		resp[i] = *NewListingResponse(&listing)
	}
	return resp
}
