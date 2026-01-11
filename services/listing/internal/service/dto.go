package service

type CreateListingParams struct {
	HostID        string
	Title         string
	Description   string
	PricePerNight int64
	ProvinceCode  string
	WardCode      string
	AddressDetail string
}
