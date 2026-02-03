package service

import (
	"context"

	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
)

func (s *ListingService) PublishListing(ctx context.Context, listingID, hostID string) (*model.Listing, error) {
	listing, err := s.listingRepo.FindByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if listing.HostID != hostID {
		return nil, model.ErrListingOwnerMismatch
	}

	if listing.Status != model.ListingStatusDraft {
		return nil, model.ErrListingNotDraft
	}

	if err = listing.ValidateForPublish(); err != nil {
		return nil, err
	}

	updatedListing, err := s.listingRepo.UpdateStatus(ctx, listingID, model.ListingStatusActive)
	if err != nil {
		return nil, err
	}

	return updatedListing, nil
}

func (s *ListingService) DeactivateListingByID(ctx context.Context, listingID, hostID string) (*model.Listing, error) {
	listing, err := s.listingRepo.FindByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if listing.HostID != hostID {
		return nil, model.ErrListingOwnerMismatch
	}

	if listing.Status != model.ListingStatusActive {
		return nil, model.ErrListingNotActive
	}

	updatedListing, err := s.listingRepo.UpdateStatus(ctx, listingID, model.ListingStatusInactive)
	if err != nil {
		return nil, err
	}

	return updatedListing, nil
}

func (s *ListingService) ReactivateListingByID(ctx context.Context, listingID, hostID string) (*model.Listing, error) {
	listing, err := s.listingRepo.FindByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if listing.HostID != hostID {
		return nil, model.ErrListingOwnerMismatch
	}

	if listing.Status != model.ListingStatusInactive {
		return nil, model.ErrListingNotInactive
	}

	if err = listing.ValidateForPublish(); err != nil {
		return nil, err
	}

	updatedListing, err := s.listingRepo.UpdateStatus(ctx, listingID, model.ListingStatusActive)
	if err != nil {
		return nil, err
	}

	return updatedListing, nil
}
