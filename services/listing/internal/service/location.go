package service

import (
	"context"

	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
)

func (s *ListingService) ListProvinces(ctx context.Context) ([]model.Province, error) {
	return s.listingRepo.ListProvinces(ctx)
}

func (s *ListingService) ListWardsByProvince(ctx context.Context, provinceCode string) ([]model.Ward, error) {
	return s.listingRepo.ListWardsByProvinceCode(ctx, provinceCode)
}
