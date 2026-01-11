package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/katatrina/airbnb-clone/pkg/token"
	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
	"github.com/katatrina/airbnb-clone/services/listing/internal/repository"
)

type ListingService struct {
	listingRepo *repository.ListingRepository
	tokenMaker  token.TokenMaker
}

func NewListingService(listingRepo *repository.ListingRepository, tokenMaker token.TokenMaker) *ListingService {
	return &ListingService{
		listingRepo: listingRepo,
		tokenMaker:  tokenMaker,
	}
}

func (s *ListingService) CreateListing(ctx context.Context, arg CreateListingParams) (*model.Listing, error) {
	province, err := s.listingRepo.GetProvinceByCode(ctx, arg.ProvinceCode)
	if err != nil {
		return nil, err
	}

	ward, err := s.listingRepo.GetWardByCode(ctx, arg.WardCode)
	if err != nil {
		return nil, err
	}

	if ward.ProvinceCode != province.Code {
		return nil, model.ErrWardProvinceMismatch
	}

	listingID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

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
