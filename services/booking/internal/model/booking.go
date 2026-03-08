package model

import "time"

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusRejected  BookingStatus = "rejected"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)

type Booking struct {
	ID            string        `db:"id"`
	ListingID     string        `db:"listing_id"`
	GuestID       string        `db:"guest_id"`
	HostID        string        `db:"host_id"`
	CheckInDate   time.Time     `db:"check_in_date"`
	CheckOutDate  time.Time     `db:"check_out_date"`
	TotalNights   int           `db:"total_nights"`
	PricePerNight int64         `db:"price_per_night"`
	TotalPrice    int64         `db:"total_price"`
	Currency      string        `db:"currency"`
	Status        BookingStatus `db:"status"`
	CreatedAt     time.Time     `db:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at"`
	DeletedAt     *time.Time    `db:"deleted_at"`
}
