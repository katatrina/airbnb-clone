package service

import (
	"context"

	"github.com/katatrina/airbnb-clone/services/booking/internal/model"
	"github.com/katatrina/airbnb-clone/services/booking/internal/repository"
)

type ListingInfo struct {
	ID            string
	HostID        string
	PricePerNight int64
	Currency      string
	Status        string
}

type ListingClient interface {
	GetActiveListingByID(ctx context.Context, id string) (*ListingInfo, error)
}

type BookingRepository interface {
	Create(ctx context.Context, booking model.Booking) (*model.Booking, error)
	FindByID(ctx context.Context, id string) (*model.Booking, error)
	UpdateStatus(ctx context.Context, id string, status model.BookingStatus) (*model.Booking, error)
	ListByGuestID(ctx context.Context, guestID string) ([]model.Booking, error)
	ListByHostID(ctx context.Context, hostID string) ([]model.Booking, error)
}

type BookingService struct {
	bookingRepo   BookingRepository
	listingClient ListingClient
}

func NewBookingService(
	bookingRepo *repository.BookingRepository,
	listingClient ListingClient,
) *BookingService {
	return &BookingService{
		bookingRepo,
		listingClient,
	}
}
