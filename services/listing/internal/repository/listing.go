package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
)

func (r *ListingRepository) CreateListing(ctx context.Context, listing model.Listing) error {
	query := `
		INSERT INTO listings (
			id, host_id, title, description, price_per_night, currency,
			province_code, province_name, district_code, district_name,
			ward_code, ward_name, address_detail,
			status, created_at, updated_at, deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13,
			$14, $15, $16, $17
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
		listing.DistrictCode,
		listing.DistrictName,
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
			province_code, province_name, district_code, district_name,
			ward_code, ward_name, address_detail,
			status, created_at, updated_at, deleted_at
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

func (r *ListingRepository) ListListingsByStatus(
	ctx context.Context,
	status model.ListingStatus,
	limit,
	offset int,
) ([]model.Listing, error) {
	query := `
		SELECT
			id, host_id, title, description, price_per_night, currency,
			province_code, province_name, district_code, district_name,
			ward_code, ward_name, address_detail,
			status, created_at, updated_at, deleted_at
		FROM listings
		WHERE status = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, _ := r.db.Query(ctx, query, status, limit, offset)
	listings, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Listing])
	if err != nil {
		return nil, err
	}

	return listings, nil
}

func (r *ListingRepository) CountListingSByStatus(
	ctx context.Context,
	status model.ListingStatus,
) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM listings
		WHERE status = $1 AND deleted_at IS NULL
	`

	var count int64
	err := r.db.QueryRow(ctx, query, status).Scan(&count)
	if err != nil {
		return 0, nil
	}

	return count, nil
}

func (r *ListingRepository) UpdateListingStatus(ctx context.Context, id string, status model.ListingStatus) error {
	query := `
		UPDATE listings
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrListingNotFound
	}

	return nil
}

func (r *ListingRepository) UpdateListingBasicInfo(ctx context.Context, id string, params model.UpdateListingBasicInfoParams) (*model.Listing, error) {
	var setClauses []string
	var args []interface{}
	paramIndex := 1

	if params.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", paramIndex))
		args = append(args, *params.Title)
		paramIndex++
	}

	if params.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", paramIndex))
		args = append(args, *params.Description)
		paramIndex++
	}

	if params.PricePerNight != nil {
		setClauses = append(setClauses, fmt.Sprintf("price_per_night = $%d", paramIndex))
		args = append(args, *params.PricePerNight)
		paramIndex++
	}

	if len(setClauses) == 0 {
		return nil, nil // No need to update or return error
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", paramIndex))
	args = append(args, time.Now())
	paramIndex++

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE listings
		SET %s
		WHERE id = $%d AND deleted_at IS NULL
		RETURNING id, host_id, title, description, price_per_night, currency,
			province_code, province_name, district_code, district_name,
			ward_code, ward_name, address_detail,
			status, created_at, updated_at, deleted_at
	`, strings.Join(setClauses, ", "), paramIndex)

	rows, _ := r.db.Query(ctx, query, args...)
	listing, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Listing])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrListingNotFound
		}
		return nil, err
	}

	return &listing, nil
}

func (r *ListingRepository) UpdateListingAddress(ctx context.Context, params model.UpdateListingAddressParams) (*model.Listing, error) {
	var setClauses []string
	var args []interface{}
	paramIndex := 1

	if params.ProvinceCode != nil {
		setClauses = append(setClauses, fmt.Sprintf("province_code = $%d", paramIndex))
		args = append(args, *params.ProvinceCode)
		paramIndex++
	}

	if params.ProvinceName != nil {
		setClauses = append(setClauses, fmt.Sprintf("province_name = $%d", paramIndex))
		args = append(args, *params.ProvinceName)
		paramIndex++
	}

	if params.DistrictCode != nil {
		setClauses = append(setClauses, fmt.Sprintf("district_code = $%d", paramIndex))
		args = append(args, *params.DistrictCode)
		paramIndex++
	}

	if params.DistrictName != nil {
		setClauses = append(setClauses, fmt.Sprintf("district_name = $%d", paramIndex))
		args = append(args, *params.DistrictName)
		paramIndex++
	}

	if params.WardCode != nil {
		setClauses = append(setClauses, fmt.Sprintf("ward_code = $%d", paramIndex))
		args = append(args, *params.WardCode)
		paramIndex++
	}

	if params.WardName != nil {
		setClauses = append(setClauses, fmt.Sprintf("ward_name = $%d", paramIndex))
		args = append(args, *params.WardName)
		paramIndex++
	}

	if params.AddressDetail != nil {
		setClauses = append(setClauses, fmt.Sprintf("address_detail = $%d", paramIndex))
		args = append(args, *params.AddressDetail)
		paramIndex++
	}

	if len(setClauses) == 0 {
		return nil, nil
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", paramIndex))
	args = append(args, time.Now())
	paramIndex++

	args = append(args, params.ListingID)

	query := fmt.Sprintf(`
		UPDATE listings
		SET %s
		WHERE id = $%d AND deleted_at IS NULL
		RETURNING id, host_id, title, description, price_per_night, currency,
			province_code, province_name, district_code, district_name,
			ward_code, ward_name, address_detail,
			status, created_at, updated_at, deleted_at
	`, strings.Join(setClauses, ", "), paramIndex)

	rows, _ := r.db.Query(ctx, query, args...)
	listing, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Listing])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrListingNotFound
		}
		return nil, err
	}

	return &listing, nil
}

func (r *ListingRepository) ListHostListings(ctx context.Context, hostID string) ([]model.Listing, error) {
	query := `
		SELECT id, host_id, title, description, price_per_night, currency,
			province_code, province_name, district_code, district_name,
			ward_code, ward_name, address_detail,
			status, created_at, updated_at, deleted_at
		FROM listings
		WHERE host_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, _ := r.db.Query(ctx, query, hostID)
	listings, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Listing])
	if err != nil {
		return nil, err
	}

	return listings, nil
}

func (r *ListingRepository) DeleteListingByID(ctx context.Context, listingID string) error {
	query := `
		UPDATE listings
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, listingID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return model.ErrListingNotFound
	}

	return nil
}
