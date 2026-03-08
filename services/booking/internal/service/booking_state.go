package service

import (
	"context"

	"github.com/katatrina/airbnb-clone/services/booking/internal/model"
)

func (s *BookingService) ConfirmBooking(
	ctx context.Context,
	bookingID, userID string,
) (*model.Booking, error) {
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if booking.HostID != userID {
		return nil, model.ErrNotBookingHost
	}

	if booking.Status != model.BookingStatusPending {
		return nil, model.ErrBookingNotPending
	}

	return s.bookingRepo.UpdateStatus(ctx, bookingID, model.BookingStatusConfirmed)
}

func (s *BookingService) RejectBooking(
	ctx context.Context,
	bookingID, userID string,
) (*model.Booking, error) {
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if booking.HostID != userID {
		return nil, model.ErrNotBookingHost
	}

	if booking.Status != model.BookingStatusPending {
		return nil, model.ErrBookingNotPending
	}

	return s.bookingRepo.UpdateStatus(ctx, bookingID, model.BookingStatusRejected)
}

func (s *BookingService) CancelBooking(
	ctx context.Context,
	bookingID, userID string,
) (*model.Booking, error) {
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if booking.GuestID != userID {
		return nil, model.ErrNotBookingGuest
	}

	if booking.Status != model.BookingStatusPending {
		return nil, model.ErrBookingNotPending
	}

	return s.bookingRepo.UpdateStatus(ctx, bookingID, model.BookingStatusCancelled)
}
