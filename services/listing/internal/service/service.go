package service

import (
	"context"

	"github.com/katatrina/airbnb-clone/pkg/token"
	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
	"github.com/katatrina/airbnb-clone/services/listing/internal/repository"
)

type ListingRepository interface {
	Create(ctx context.Context, listing model.Listing) (*model.Listing, error)
	FindByID(ctx context.Context, id string) (*model.Listing, error)
	Delete(ctx context.Context, id string) error

	ListByStatus(ctx context.Context, status model.ListingStatus, limit, offset int) ([]model.Listing, error)
	ListByHostID(ctx context.Context, hostID string) ([]model.Listing, error)
	CountByStatus(ctx context.Context, status model.ListingStatus) (int64, error)

	UpdateStatus(ctx context.Context, id string, status model.ListingStatus) (*model.Listing, error)
	UpdateBasicInfo(ctx context.Context, id string, arg model.UpdateListingBasicInfoParams) (*model.Listing, error)
	UpdateAddress(ctx context.Context, arg model.UpdateListingAddressParams) (*model.Listing, error)
}

type LocationRepository interface {
	FindProvinceByCode(ctx context.Context, code string) (*model.Province, error)
	FindDistrictByCode(ctx context.Context, code string) (*model.District, error)
	FindWardByCode(ctx context.Context, code string) (*model.Ward, error)

	ListProvinces(ctx context.Context) ([]model.Province, error)
	ListDistrictsByProvinceCode(ctx context.Context, provinceCode string) ([]model.District, error)
	ListWardsByDistrictCode(ctx context.Context, districtCode string) ([]model.Ward, error)
}

type ListingService struct {
	listingRepo  ListingRepository
	locationRepo LocationRepository
	tokenMaker   token.TokenMaker
}

func NewListingService(
	listingRepo *repository.ListingRepository,
	locationRepo *repository.LocationRepository,
	tokenMaker token.TokenMaker,
) *ListingService {
	return &ListingService{
		listingRepo:  listingRepo,
		locationRepo: locationRepo,
		tokenMaker:   tokenMaker,
	}
}
