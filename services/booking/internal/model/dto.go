package model

import "time"

type CreateBookingParams struct {
	ListingID    string
	GuestID      string
	CheckInDate  time.Time
	CheckOutDate time.Time
}
