package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
)

func (r *ListingRepository) GetProvinceByCode(ctx context.Context, code string) (*model.Province, error) {
	query := `
		SELECT code, full_name, created_at
		FROM provinces
		WHERE code = $1
	`

	rows, _ := r.db.Query(ctx, query, code)
	province, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Province])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrProvinceCodeNotFound
		}
		return nil, err
	}

	return &province, nil
}

func (r *ListingRepository) GetWardByCode(ctx context.Context, code string) (*model.Ward, error) {
	query := `
		SELECT code, full_name, province_code, created_at
		FROM wards
		WHERE code = $1
	`

	rows, _ := r.db.Query(ctx, query, code)
	ward, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Ward])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrWardCodeNotFound
		}
		return nil, err
	}

	return &ward, nil
}

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

func (r *ListingRepository) ListWardsByProvinceCode(ctx context.Context, provinceCode string) ([]model.Ward, error) {
	query := `
		SELECT code, full_name, province_code, created_at
		FROM wards
		WHERE province_code = $1
		ORDER BY full_name
	`

	rows, _ := r.db.Query(ctx, query, provinceCode)
	wards, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Ward])
	if err != nil {
		return nil, err
	}

	return wards, nil
}
