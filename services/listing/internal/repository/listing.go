package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
)

func (r *ListingRepository) CreateListing(ctx context.Context, listing model.Listing) error {
	query := `
		INSERT INTO listings (
			id, host_id, title, description, price_per_night, currency,
			province_code, province_name, ward_code, ward_name,
			address_detail, status, created_at, updated_at, deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14, $15
		)
	`

	_, err := r.db.Exec(ctx, query,
		listing.ID,
		listing.HostID,
		listing.Title,
		listing.Description,
		listing.PricePerNight,
		listing.Currency,
		listing.ProvinceCode,
		listing.ProvinceName,
		listing.WardCode,
		listing.WardName,
		listing.AddressDetail,
		listing.Status,
		listing.CreatedAt,
		listing.UpdatedAt,
		listing.DeletedAt,
	)
	return err
}

func (r *ListingRepository) FindListingByID(ctx context.Context, id string) (*model.Listing, error) {
	query := `
		SELECT
			id, host_id, title, description, price_per_night, currency,
			province_code, province_name, ward_code, ward_name,
			address_detail, status, created_at, updated_at, deleted_at
		FROM listings
		WHERE id = $1 AND deleted_at IS NULL
	`

	rows, _ := r.db.Query(ctx, query, id)
	listing, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Listing])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrListingNotFound
		}
		return nil, err
	}

	return &listing, nil
}

func (r *ListingRepository) ListListingsByStatus(ctx context.Context, status model.ListingStatus) ([]model.Listing, error) {
	query := `
		SELECT
			id, host_id, title, description, price_per_night, currency,
			province_code, province_name, ward_code, ward_name,
			address_detail, status, created_at, updated_at, deleted_at
		FROM listings
		WHERE status = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, _ := r.db.Query(ctx, query, status)
	listings, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Listing])
	if err != nil {
		return nil, err
	}

	return listings, nil
}
