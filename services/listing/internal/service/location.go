package service

import (
	"context"

	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
)

func (s *ListingService) ListProvinces(ctx context.Context) ([]model.Province, error) {
	return s.listingRepo.ListProvinces(ctx)
}

func (s *ListingService) ListDistrictsByProvince(ctx context.Context, provinceCode string) ([]model.District, error) {
	_, err := s.listingRepo.GetProvinceByCode(ctx, provinceCode)
	if err != nil {
		return nil, err
	}

	return s.listingRepo.ListDistrictsByProvinceCode(ctx, provinceCode)
}

func (s *ListingService) ListWardsByDistrict(ctx context.Context, districtCode string) ([]model.Ward, error) {
	_, err := s.listingRepo.GetDistrictByCode(ctx, districtCode)
	if err != nil {
		return nil, err
	}

	return s.listingRepo.ListWardsByDistrictCode(ctx, districtCode)
}
