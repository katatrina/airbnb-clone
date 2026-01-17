package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
)

func (s *ListingService) CreateListing(ctx context.Context, arg CreateListingParams) (*model.Listing, error) {
	// Validate Province
	province, err := s.listingRepo.GetProvinceByCode(ctx, arg.ProvinceCode)
	if err != nil {
		return nil, err
	}

	// Validate District
	district, err := s.listingRepo.GetDistrictByCode(ctx, arg.DistrictCode)
	if err != nil {
		return nil, err
	}

	// Validate District belongs to Province
	if district.ProvinceCode != province.Code {
		return nil, model.ErrDistrictProvinceMismatch
	}

	// Validate Ward
	ward, err := s.listingRepo.GetWardByCode(ctx, arg.WardCode)
	if err != nil {
		return nil, err
	}

	// Validate Ward belongs to District
	if ward.DistrictCode != district.Code {
		return nil, model.ErrWardDistrictMismatch
	}

	listingID, _ := uuid.NewV7()
	now := time.Now()
	listing := model.Listing{
		ID:            listingID.String(),
		HostID:        arg.HostID,
		Title:         arg.Title,
		Description:   arg.Description,
		PricePerNight: arg.PricePerNight,
		Currency:      model.ListingCurrencyVND,
		ProvinceCode:  arg.ProvinceCode,
		ProvinceName:  province.FullName,
		DistrictCode:  arg.DistrictCode,
		DistrictName:  district.FullName,
		WardCode:      arg.WardCode,
		WardName:      ward.FullName,
		AddressDetail: arg.AddressDetail,
		Status:        model.ListingStatusDraft,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	err = s.listingRepo.CreateListing(ctx, listing)
	if err != nil {
		return nil, err
	}

	return &listing, nil
}

func (s *ListingService) GetActiveListingByID(ctx context.Context, listingID string) (*model.Listing, error) {
	listing, err := s.listingRepo.FindListingByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if listing.Status != model.ListingStatusActive {
		return nil, model.ErrListingNotFound
	}

	return listing, nil
}

func (s *ListingService) ListActiveListings(ctx context.Context, limit, offset int) ([]model.Listing, int64, error) {
	listings, err := s.listingRepo.ListListingsByStatus(ctx, model.ListingStatusActive, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.listingRepo.CountListingSByStatus(ctx, model.ListingStatusActive)
	if err != nil {
		return nil, 0, err
	}

	return listings, total, nil
}
