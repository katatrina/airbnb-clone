package model

type CreateListingParams struct {
	HostID        string
	Title         string
	Description   string
	PricePerNight int64
	ProvinceCode  int32
	DistrictCode  int32
	WardCode      int32
	AddressDetail string
}

type UpdateListingBasicInfoParams struct {
	Title         *string
	Description   *string
	PricePerNight *int64
}

type UpdateListingAddressParams struct {
	ListingID     string
	HostID        string
	ProvinceCode  *int32
	ProvinceName  *string
	DistrictCode  *int32
	DistrictName  *string
	WardCode      *int32
	WardName      *string
	AddressDetail *string
}
