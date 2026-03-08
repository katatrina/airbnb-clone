package handler

import "github.com/katatrina/airbnb-clone/services/booking/internal/service"

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}
