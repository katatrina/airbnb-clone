package handler

import "github.com/katatrina/airbnb-clone/services/listing/internal/model"

type CreateListingRequest struct {
	Title         string `json:"title" validate:"required,min=10,max=200" normalize:"trim,singlespace"`
	Description   string `json:"description" validate:"omitempty,max=2000" normalize:"trim,singlespace"`
	PricePerNight int64  `json:"pricePerNight" validate:"required,gte=1"`
	ProvinceCode  int32  `json:"provinceCode" validate:"required"`
	DistrictCode  int32  `json:"districtCode" validate:"required"`
	WardCode      int32  `json:"wardCode" validate:"required"`
	AddressDetail string `json:"addressDetail" validate:"required,min=10,max=500" normalize:"trim,singlespace"`
}

type UpdateListingBasicInfoRequest struct {
	Title         *string `json:"title" validate:"omitnil,min=10,max=200" normalize:"trim,singlespace"`
	Description   *string `json:"description" validate:"omitnil,max=2000" normalize:"trim,singlespace"`
	PricePerNight *int64  `json:"pricePerNight" validate:"omitnil,gte=1"`
}

type UpdateListingAddressRequest struct {
	ProvinceCode  *int32  `json:"provinceCode" validate:"required_with=DistrictCode WardCode"`
	DistrictCode  *int32  `json:"districtCode" validate:"required_with=ProvinceCode WardCode"`
	WardCode      *int32  `json:"wardCode" validate:"required_with=ProvinceCode DistrictCode"`
	AddressDetail *string `json:"addressDetail" validate:"omitnil,min=10,max=500" normalize:"trim,singlespace"`
}

type ListingResponse struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	PricePerNight int64  `json:"pricePerNight"`
	Currency      string `json:"currency"`
	ProvinceCode  int32  `json:"provinceCode"`
	ProvinceName  string `json:"provinceName"`
	DistrictCode  int32  `json:"districtCode"`
	DistrictName  string `json:"districtName"`
	WardCode      int32  `json:"wardCode"`
	WardName      string `json:"wardName"`
	AddressDetail string `json:"addressDetail"`
	Status        string `json:"status"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
}

type ProvinceResponse struct {
	Code     int32  `json:"code"`
	FullName string `json:"fullName"`
}

type DistrictResponse struct {
	Code         int32  `json:"code"`
	FullName     string `json:"fullName"`
	ProvinceCode int32  `json:"provinceCode"`
}

type WardResponse struct {
	Code         int32  `json:"code"`
	FullName     string `json:"fullName"`
	DistrictCode int32  `json:"districtCode"`
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
	for i := range listings {
		l := &listings[i]

		resp[i] = ListingResponse{
			ID:            l.ID,
			Title:         l.Title,
			Description:   l.Description,
			PricePerNight: l.PricePerNight,
			Currency:      string(l.Currency),
			ProvinceCode:  l.ProvinceCode,
			ProvinceName:  l.ProvinceName,
			DistrictCode:  l.DistrictCode,
			DistrictName:  l.DistrictName,
			WardCode:      l.WardCode,
			WardName:      l.WardName,
			AddressDetail: l.AddressDetail,
			Status:        string(l.Status),
			CreatedAt:     l.CreatedAt.Unix(),
			UpdatedAt:     l.UpdatedAt.Unix(),
		}
	}
	return resp
}
