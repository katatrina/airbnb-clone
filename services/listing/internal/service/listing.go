package service

import (
	"context"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
)

func (s *ListingService) CreateListing(ctx context.Context, arg CreateListingParams) (*model.Listing, error) {
	province, err := s.listingRepo.FindProvinceByCode(ctx, arg.ProvinceCode)
	if err != nil {
		return nil, err
	}

	district, err := s.listingRepo.FindDistrictByCode(ctx, arg.DistrictCode)
	if err != nil {
		return nil, err
	}

	if district.ProvinceCode != province.Code {
		return nil, model.ErrDistrictProvinceMismatch
	}

	ward, err := s.listingRepo.FindWardByCode(ctx, arg.WardCode)
	if err != nil {
		return nil, err
	}

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

func (s *ListingService) GetActiveListingByID(ctx context.Context, id string) (*model.Listing, error) {
	listing, err := s.listingRepo.FindListingByID(ctx, id)
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

func (s *ListingService) PublishListing(ctx context.Context, listingID, userID string) (*model.Listing, error) {
	listing, err := s.listingRepo.FindListingByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if listing.HostID != userID {
		return nil, model.ErrListingOwnerMismatch
	}

	if listing.Status != model.ListingStatusDraft {
		return nil, model.ErrListingNotDraft
	}

	if utf8.RuneCountInString(listing.Description) < 50 {
		return nil, model.ErrListingIncomplete
	}

	err = s.listingRepo.UpdateListingStatus(ctx, listingID, model.ListingStatusActive)
	if err != nil {
		return nil, err
	}
	listing.Status = model.ListingStatusActive

	return listing, nil
}
