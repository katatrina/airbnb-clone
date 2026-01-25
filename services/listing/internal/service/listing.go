package service

import (
	"context"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
)

func (s *ListingService) CreateListing(ctx context.Context, arg model.CreateListingParams) (*model.Listing, error) {
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

func (s *ListingService) PublishListing(ctx context.Context, listingID, hostID string) (*model.Listing, error) {
	listing, err := s.listingRepo.FindListingByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if listing.HostID != hostID {
		return nil, model.ErrListingOwnerMismatch
	}

	if listing.Status != model.ListingStatusDraft {
		return nil, model.ErrListingNotDraft
	}

	// TODO: We can add more validation rules here later

	if utf8.RuneCountInString(listing.Description) < 50 {
		return nil, model.ErrListingIncomplete
	}

	err = s.listingRepo.UpdateListingStatus(ctx, listingID, model.ListingStatusActive)
	if err != nil {
		return nil, err
	}
	listing.Status = model.ListingStatusActive
	listing.UpdatedAt = time.Now()

	return listing, nil
}

func (s *ListingService) UpdateListingBasicInfo(ctx context.Context, listingID, hostID string, arg model.UpdateListingBasicInfoParams) (*model.Listing, error) {
	listing, err := s.listingRepo.FindListingByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if listing.HostID != hostID {
		return nil, model.ErrListingOwnerMismatch
	}

	if listing.Status == model.ListingStatusActive {
		return nil, model.ErrActiveListingCannotBeUpdated
	}

	updatedListing, err := s.listingRepo.UpdateListingBasicInfo(ctx, listingID, arg)
	if err != nil {
		return nil, err
	}

	if updatedListing == nil {
		return listing, nil
	}

	return updatedListing, nil
}

func (s *ListingService) UpdateListingAddress(ctx context.Context, arg model.UpdateListingAddressParams) (*model.Listing, error) {
	listing, err := s.listingRepo.FindListingByID(ctx, arg.ListingID)
	if err != nil {
		return nil, err
	}

	if listing.HostID != arg.HostID {
		return nil, model.ErrListingOwnerMismatch
	}

	if listing.Status == model.ListingStatusActive {
		return nil, model.ErrActiveListingCannotBeUpdated
	}

	if arg.ProvinceCode != nil && arg.DistrictCode != nil && arg.WardCode != nil {
		province, err := s.listingRepo.FindProvinceByCode(ctx, *arg.ProvinceCode)
		if err != nil {
			return nil, err
		}

		district, err := s.listingRepo.FindDistrictByCode(ctx, *arg.DistrictCode)
		if err != nil {
			return nil, err
		}

		if district.ProvinceCode != province.Code {
			return nil, model.ErrDistrictProvinceMismatch
		}

		ward, err := s.listingRepo.FindWardByCode(ctx, *arg.WardCode)
		if err != nil {
			return nil, err
		}

		if ward.DistrictCode != district.Code {
			return nil, model.ErrWardDistrictMismatch
		}

		arg.ProvinceName = &province.FullName
		arg.DistrictName = &district.FullName
		arg.WardName = &ward.FullName
	}

	updatedListing, err := s.listingRepo.UpdateListingAddress(ctx, arg)
	if err != nil {
		return nil, err
	}

	if updatedListing == nil {
		return listing, nil
	}

	return updatedListing, nil
}

func (s *ListingService) ListHostListings(ctx context.Context, hostID string) ([]model.Listing, error) {
	return s.listingRepo.ListHostListings(ctx, hostID)
}

func (s *ListingService) GetHostListingByID(ctx context.Context, listingID, hostID string) (*model.Listing, error) {
	listing, err := s.listingRepo.FindListingByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if listing.HostID != hostID {
		return nil, model.ErrListingOwnerMismatch
	}

	return listing, nil
}

func (s *ListingService) DeleteListingByID(ctx context.Context, listingID, hostID string) error {
	listing, err := s.listingRepo.FindListingByID(ctx, listingID)
	if err != nil {
		return err
	}

	if listing.HostID != hostID {
		return model.ErrListingOwnerMismatch
	}

	// Allow deleting listing from any status

	// TODO: Check active booking(s) related to this listing (future in booking service)

	return s.listingRepo.DeleteListingByID(ctx, listingID)
}

func (s *ListingService) DeactivateListingByID(ctx context.Context, listingID, hostID string) (*model.Listing, error) {
	listing, err := s.listingRepo.FindListingByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if listing.HostID != hostID {
		return nil, model.ErrListingOwnerMismatch
	}

	if listing.Status != model.ListingStatusActive {
		return nil, model.ErrListingNotActive
	}

	newStatus := model.ListingStatusInactive
	err = s.listingRepo.UpdateListingStatus(ctx, listingID, newStatus)
	if err != nil {
		return nil, err
	}
	listing.Status = newStatus

	return listing, nil
}

func (s *ListingService) ReactivateListingByID(ctx context.Context, listingID, hostID string) (*model.Listing, error) {
	listing, err := s.listingRepo.FindListingByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if listing.HostID != hostID {
		return nil, model.ErrListingOwnerMismatch
	}

	if listing.Status != model.ListingStatusInactive {
		return nil, model.ErrListingNotInactive
	}

	// We can skip validate completeness here (as we do when publishing listing)

	newStatus := model.ListingStatusActive
	err = s.listingRepo.UpdateListingStatus(ctx, listingID, newStatus)
	if err != nil {
		return nil, err
	}
	listing.Status = newStatus

	return listing, nil
}
