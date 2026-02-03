package model

import (
	"time"
)

type Province struct {
	Code      string    `db:"code"`
	FullName  string    `db:"full_name"`
	CreatedAt time.Time `db:"created_at"`
}

type District struct {
	Code         string    `db:"code"`
	FullName     string    `db:"full_name"`
	ProvinceCode string    `db:"province_code"`
	CreatedAt    time.Time `db:"created_at"`
}

type Ward struct {
	Code         string    `db:"code"`
	FullName     string    `db:"full_name"`
	DistrictCode string    `db:"district_code"`
	CreatedAt    time.Time `db:"created_at"`
}
