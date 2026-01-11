package repository

import (
	"context"

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
