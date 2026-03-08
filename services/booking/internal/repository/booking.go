package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/katatrina/airbnb-clone/services/booking/internal/model"
)

func (r *BookingRepository) Create(
	ctx context.Context,
	booking model.Booking,
) (*model.Booking, error) {
	query := `
        INSERT INTO bookings (
            id, listing_id, guest_id, host_id,
            check_in_date, check_out_date, total_nights,
            price_per_night, total_price, currency,
            status, created_at, updated_at, deleted_at
        ) VALUES (
            $1, $2, $3, $4,
            $5, $6, $7,
            $8, $9, $10,
            $11, $12, $13, $14
        )
        RETURNING
            id, listing_id, guest_id, host_id,
            check_in_date, check_out_date, total_nights,
            price_per_night, total_price, currency,
            status, created_at, updated_at, deleted_at
    `

	rows, _ := r.db.Query(ctx, query,
		booking.ID, booking.ListingID, booking.GuestID, booking.HostID,
		booking.CheckInDate, booking.CheckOutDate, booking.TotalNights,
		booking.PricePerNight, booking.TotalPrice, booking.Currency,
		booking.Status, booking.CreatedAt, booking.UpdatedAt, booking.DeletedAt,
	)

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Booking])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23P01" &&
				pgErr.ConstraintName == "no_overlapping_bookings" {
				return nil, model.ErrDatesUnavailable
			}
		}
		return nil, err
	}

	return &created, nil
}

func (r *BookingRepository) FindByID(
	ctx context.Context,
	id string,
) (*model.Booking, error) {
	query := `
        SELECT
            id, listing_id, guest_id, host_id,
            check_in_date, check_out_date, total_nights,
            price_per_night, total_price, currency,
            status, created_at, updated_at, deleted_at
        FROM bookings
        WHERE id = $1 AND deleted_at IS NULL
    `

	rows, _ := r.db.Query(ctx, query, id)
	booking, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Booking])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrBookingNotFound
		}
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status model.BookingStatus,
) (*model.Booking, error) {
	query := `
        UPDATE bookings
        SET status = $1, updated_at = NOW()
        WHERE id = $2 AND deleted_at IS NULL
        RETURNING
            id, listing_id, guest_id, host_id,
            check_in_date, check_out_date, total_nights,
            price_per_night, total_price, currency,
            status, created_at, updated_at, deleted_at
    `

	rows, _ := r.db.Query(ctx, query, status, id)
	booking, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Booking])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrBookingNotFound
		}
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepository) ListByGuestID(
	ctx context.Context,
	guestID string,
) ([]model.Booking, error) {
	query := `
        SELECT
            id, listing_id, guest_id, host_id,
            check_in_date, check_out_date, total_nights,
            price_per_night, total_price, currency,
            status, created_at, updated_at, deleted_at
        FROM bookings
        WHERE guest_id = $1 AND deleted_at IS NULL
        ORDER BY created_at DESC
    `

	rows, _ := r.db.Query(ctx, query, guestID)
	bookings, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Booking])
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *BookingRepository) ListByHostID(
	ctx context.Context,
	hostID string,
) ([]model.Booking, error) {
	query := `
        SELECT
            id, listing_id, guest_id, host_id,
            check_in_date, check_out_date, total_nights,
            price_per_night, total_price, currency,
            status, created_at, updated_at, deleted_at
        FROM bookings
        WHERE host_id = $1 AND deleted_at IS NULL
        ORDER BY created_at DESC
    `

	rows, _ := r.db.Query(ctx, query, hostID)
	bookings, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Booking])
	if err != nil {
		return nil, err
	}

	return bookings, nil
}
