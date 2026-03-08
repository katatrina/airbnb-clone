package handler

import "github.com/katatrina/airbnb-clone/services/booking/internal/model"

type CreateBookingRequest struct {
	ListingID    string `json:"listingId" validate:"required" normalize:"trim"`
	CheckInDate  string `json:"checkInDate" validate:"required"`
	CheckOutDate string `json:"checkOutDate" validate:"required"`
}

type BookingResponse struct {
	ID            string `json:"id"`
	CheckInDate   string `json:"checkInDate"`
	CheckOutDate  string `json:"checkOutDate"`
	TotalNights   int    `json:"totalNights"`
	PricePerNight int64  `json:"pricePerNight"`
	TotalPrice    int64  `json:"totalPrice"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
}

func NewBookingResponse(b *model.Booking) *BookingResponse {
	return &BookingResponse{
		ID:            b.ID,
		CheckInDate:   b.CheckInDate.Format("2006-01-02"),
		CheckOutDate:  b.CheckOutDate.Format("2006-01-02"),
		TotalNights:   b.TotalNights,
		PricePerNight: b.PricePerNight,
		TotalPrice:    b.TotalPrice,
		Currency:      b.Currency,
		Status:        string(b.Status),
		CreatedAt:     b.CreatedAt.Unix(),
		UpdatedAt:     b.UpdatedAt.Unix(),
	}
}

func NewBookingsResponse(bookings []model.Booking) []BookingResponse {
	resp := make([]BookingResponse, len(bookings))
	for i := range bookings {
		b := &bookings[i]
		resp[i] = *NewBookingResponse(b)
	}
	return resp
}
