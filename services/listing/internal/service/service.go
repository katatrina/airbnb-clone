package service

import (
	"github.com/katatrina/airbnb-clone/pkg/token"
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
