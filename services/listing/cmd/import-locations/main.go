package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katatrina/airbnb-clone/services/listing/config"
)

// Ward represents a ward in the JSON file
type Ward struct {
	Name          string `json:"name"`
	Code          int    `json:"code"`
	Codename      string `json:"codename"`
	DivisionType  string `json:"division_type"`
	ShortCodename string `json:"short_codename"`
}

// District represents a district in the JSON file
type District struct {
	Name          string `json:"name"`
	Code          int    `json:"code"`
	Codename      string `json:"codename"`
	DivisionType  string `json:"division_type"`
	ShortCodename string `json:"short_codename"`
	Wards         []Ward `json:"wards"`
}

// Province represents a province in the JSON file
type Province struct {
	Name         string     `json:"name"`
	Code         int        `json:"code"`
	Codename     string     `json:"codename"`
	DivisionType string     `json:"division_type"`
	PhoneCode    int        `json:"phone_code"`
	Districts    []District `json:"districts"`
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create database connection
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Test connection
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Read and parse JSON file from same directory
	jsonFilePath := "cmd/import-locations/vn-provinces.json"
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		log.Fatalf("Failed to read JSON file %s: %v", jsonFilePath, err)
	}

	var provinces []Province
	if err := json.Unmarshal(data, &provinces); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	log.Printf("Parsed %d provinces from JSON file", len(provinces))

	// Import data
	if err := importLocations(ctx, dbPool, provinces); err != nil {
		log.Fatalf("Failed to import locations: %v", err)
	}

	log.Println("Import completed successfully!")
}

func importLocations(ctx context.Context, pool *pgxpool.Pool, provinces []Province) error {
	// Start transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	now := time.Now()
	provinceCount := 0
	districtCount := 0
	wardCount := 0

	// Insert provinces, districts and wards
	for _, province := range provinces {
		provinceCode := strconv.Itoa(province.Code)

		// Insert province
		_, err := tx.Exec(ctx,
			`INSERT INTO provinces (code, full_name, created_at)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (code) DO UPDATE
			 SET full_name = EXCLUDED.full_name`,
			provinceCode, province.Name, now)
		if err != nil {
			return fmt.Errorf("failed to insert province %s: %w", provinceCode, err)
		}
		provinceCount++

		// Insert districts for this province
		for _, district := range province.Districts {
			districtCode := strconv.Itoa(district.Code)

			_, err := tx.Exec(ctx,
				`INSERT INTO districts (code, full_name, province_code, created_at)
				 VALUES ($1, $2, $3, $4)
				 ON CONFLICT (code) DO UPDATE
				 SET full_name = EXCLUDED.full_name, province_code = EXCLUDED.province_code`,
				districtCode, district.Name, provinceCode, now)
			if err != nil {
				return fmt.Errorf("failed to insert district %s: %w", districtCode, err)
			}
			districtCount++

			// Insert wards for this district
			for _, ward := range district.Wards {
				wardCode := strconv.Itoa(ward.Code)

				_, err := tx.Exec(ctx,
					`INSERT INTO wards (code, full_name, district_code, created_at)
					 VALUES ($1, $2, $3, $4)
					 ON CONFLICT (code) DO UPDATE
					 SET full_name = EXCLUDED.full_name, district_code = EXCLUDED.district_code`,
					wardCode, ward.Name, districtCode, now)
				if err != nil {
					return fmt.Errorf("failed to insert ward %s: %w", wardCode, err)
				}
				wardCount++
			}
		}

		log.Printf("Imported province: %s with %d districts", province.Name, len(province.Districts))
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Total imported: %d provinces, %d districts, %d wards", provinceCount, districtCount, wardCount)
	return nil
}
