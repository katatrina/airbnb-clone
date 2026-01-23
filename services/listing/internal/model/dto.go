package model

type CreateListingParams struct {
	HostID        string
	Title         string
	Description   string
	PricePerNight int64
	ProvinceCode  string
	DistrictCode  string
	WardCode      string
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
	ProvinceCode  *string
	ProvinceName  *string
	DistrictCode  *string
	DistrictName  *string
	WardCode      *string
	WardName      *string
	AddressDetail *string
}
