package main

//
// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"
//
// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/katatrina/airbnb-clone/services/listing/config"
// 	"github.com/katatrina/airbnb-clone/services/listing/internal/repository"
// )
//
// type Ward struct {
// 	Code         string `json:"Code"`
// 	FullName     string `json:"FullName"`
// 	ProvinceCode string `json:"ProvinceCode"`
// }
//
// type Province struct {
// 	Code     string `json:"Code"`
// 	FullName string `json:"FullName"`
// 	Wards    []Ward `json:"Wards"`
// }
//
// func main() {
// 	// Load configuration
// 	cfg, err := config.LoadConfig(".env")
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}
//
// 	// Create repository connection
// 	ctx := context.Background()
// 	dbPool, err := repository.NewListingRepository(ctx, cfg.DatabaseURL)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to repository: %v", err)
// 	}
// 	defer dbPool.Close()
//
// 	log.Println("Connected to repository successfully")
//
// 	// Read and parse JSON file from same directory
// 	jsonFilePath := "cmd/import-locations/data.json"
// 	data, err := os.ReadFile(jsonFilePath)
// 	if err != nil {
// 		log.Fatalf("Failed to read JSON file %s: %v", jsonFilePath, err)
// 	}
//
// 	var provinces []Province
// 	if err := json.Unmarshal(data, &provinces); err != nil {
// 		log.Fatalf("Failed to parse JSON: %v", err)
// 	}
//
// 	// Import data
// 	if err := importLocations(ctx, dbPool, provinces); err != nil {
// 		log.Fatalf("Failed to import locations: %v", err)
// 	}
//
// 	log.Println("Import completed successfully!")
// }
//
// func importLocations(ctx context.Context, pool *pgxpool.Pool, provinces []Province) error {
// 	// Start transaction
// 	tx, err := pool.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to begin transaction: %w", err)
// 	}
// 	defer tx.Rollback(ctx)
//
// 	now := time.Now()
// 	provinceCount := 0
// 	wardCount := 0
//
// 	// Insert provinces and wards
// 	for _, province := range provinces {
// 		// Insert province
// 		_, err := tx.Exec(ctx,
// 			`INSERT INTO provinces (code, full_name, created_at)
// 			 VALUES ($1, $2, $3)
// 			 ON CONFLICT (code) DO UPDATE
// 			 SET full_name = EXCLUDED.full_name`,
// 			province.Code, province.FullName, now)
// 		if err != nil {
// 			return fmt.Errorf("failed to insert province %s: %w", province.Code, err)
// 		}
// 		provinceCount++
//
// 		// Insert wards for this province
// 		for _, ward := range province.Wards {
// 			_, err := tx.Exec(ctx,
// 				`INSERT INTO wards (code, full_name, province_code, created_at)
// 				 VALUES ($1, $2, $3, $4)
// 				 ON CONFLICT (code) DO UPDATE
// 				 SET full_name = EXCLUDED.full_name, province_code = EXCLUDED.province_code`,
// 				ward.Code, ward.FullName, province.Code, now)
// 			if err != nil {
// 				return fmt.Errorf("failed to insert ward %s: %w", ward.Code, err)
// 			}
// 			wardCount++
// 		}
//
// 		log.Printf("Imported province: %s with %d wards", province.FullName, len(province.Wards))
// 	}
//
// 	// Commit transaction
// 	if err := tx.Commit(ctx); err != nil {
// 		return fmt.Errorf("failed to commit transaction: %w", err)
// 	}
//
// 	log.Printf("Total imported: %d provinces, %d wards", provinceCount, wardCount)
// 	return nil
// }
