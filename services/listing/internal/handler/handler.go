package handler

import (
	"github.com/katatrina/airbnb-clone/services/listing/internal/service"
)

type ListingHandler struct {
	listingService *service.ListingService
}

func NewListingHandler(listingService *service.ListingService) *ListingHandler {
	return &ListingHandler{
		listingService: listingService,
	}
}
