package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
)

func (r *ListingRepository) ListProvinces(ctx context.Context) ([]model.Province, error) {
	query := `
		SELECT code, full_name, created_at
		FROM provinces 
		ORDER BY full_name
	`

	rows, _ := r.db.Query(ctx, query)
	provinces, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Province])
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

func (r *ListingRepository) ListWards(ctx context.Context, provinceCode string) ([]model.Ward, error) {
	query := `
		SELECT code, full_name, province_code, created_at
		FROM wards
		WHERE province_code = $1
		ORDER BY full_name
	`

	rows, err := r.db.Query(ctx, query, provinceCode)
	if err != nil {
		return nil, err
	}

	wards, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Ward])
	if err != nil {
		return nil, err
	}

	return wards, nil
}
