package model

import (
	"time"
)

type Province struct {
	Code      int32     `db:"code"`
	FullName  string    `db:"full_name"`
	CreatedAt time.Time `db:"created_at"`
}

type District struct {
	Code         int32     `db:"code"`
	FullName     string    `db:"full_name"`
	ProvinceCode int32     `db:"province_code"`
	CreatedAt    time.Time `db:"created_at"`
}

type Ward struct {
	Code         int32     `db:"code"`
	FullName     string    `db:"full_name"`
	DistrictCode int32     `db:"district_code"`
	CreatedAt    time.Time `db:"created_at"`
}
