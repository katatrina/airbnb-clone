package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/services/listing/config"
	"github.com/katatrina/airbnb-clone/services/listing/internal/db"
	"github.com/katatrina/airbnb-clone/services/listing/internal/handler"
)

func main() {
	fmt.Println(response.TestVal)

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create db pool: %v", err)
	}

	if err = dbPool.Ping(ctx); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	log.Println("connected to db")

	listingRepo := db.NewListingRepository(dbPool)
	h := handler.NewHandler(listingRepo, cfg)

	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.GET("/health", h.Health)
	v1.GET("/provinces", h.ListProvinces)
	v1.GET("/wards", h.ListWards)

	log.Fatal(router.Run(fmt.Sprintf(":%s", cfg.ServerPort)))
}
