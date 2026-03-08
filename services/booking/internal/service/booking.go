package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/katatrina/airbnb-clone/services/booking/internal/model"
)

func (s *BookingService) CreateBooking(ctx context.Context, arg model.CreateBookingParams) (*model.Booking, error) {
	if !arg.CheckOutDate.After(arg.CheckInDate) {
		return nil, model.ErrInvalidDateRange
	}

	today := time.Now().Truncate(24 * time.Hour)
	if arg.CheckInDate.Before(today) {
		return nil, model.ErrCheckInPast
	}

	listing, err := s.listingClient.GetActiveListingByID(ctx, arg.ListingID)
	if err != nil {
		return nil, err
	}

	if listing.HostID == arg.GuestID {
		return nil, model.ErrSelfBooking
	}

	totalNights := int(arg.CheckOutDate.Sub(arg.CheckInDate).Hours() / 24)
	totalPrice := listing.PricePerNight * int64(totalNights)

	bookingID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate booking ID: %w", err)
	}

	now := time.Now()
	booking := model.Booking{
		ID:            bookingID.String(),
		ListingID:     arg.ListingID,
		GuestID:       arg.GuestID,
		HostID:        listing.HostID,
		CheckInDate:   arg.CheckInDate,
		CheckOutDate:  arg.CheckOutDate,
		TotalNights:   totalNights,
		PricePerNight: listing.PricePerNight,
		TotalPrice:    totalPrice,
		Currency:      listing.Currency,
		Status:        model.BookingStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	createdBooking, err := s.bookingRepo.Create(ctx, booking)
	if err != nil {
		return nil, err
	}

	return createdBooking, nil
}
