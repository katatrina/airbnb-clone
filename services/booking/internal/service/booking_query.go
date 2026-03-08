// services/booking/internal/service/booking_query.go

package service

import (
	"context"

	"github.com/katatrina/airbnb-clone/services/booking/internal/model"
)

func (s *BookingService) GetBookingByID(
	ctx context.Context,
	bookingID, userID string,
) (*model.Booking, error) {
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if booking.GuestID != userID && booking.HostID != userID {
		return nil, model.ErrBookingNotFound
	}

	return booking, nil
}

func (s *BookingService) ListGuestBookings(
	ctx context.Context,
	guestID string,
) ([]model.Booking, error) {
	return s.bookingRepo.ListByGuestID(ctx, guestID)
}

func (s *BookingService) ListHostBookings(
	ctx context.Context,
	hostID string,
) ([]model.Booking, error) {
	return s.bookingRepo.ListByHostID(ctx, hostID)
}
