package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katatrina/airbnb-clone/pkg/token"
	"github.com/katatrina/airbnb-clone/services/listing/config"
	"github.com/katatrina/airbnb-clone/services/listing/internal/handler"
	"github.com/katatrina/airbnb-clone/services/listing/internal/middleware"
	"github.com/katatrina/airbnb-clone/services/listing/internal/repository"
	"github.com/katatrina/airbnb-clone/services/listing/internal/service"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer db.Close()

	if err = db.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	tokenMaker, err := token.NewJWTMaker([]byte(cfg.JWTSecret), cfg.JWTExpiry)
	if err != nil {
		log.Fatalf("Failed to create token maker: %v", err)
	}

	listingRepo := repository.NewListingRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	listingService := service.NewListingService(listingRepo, locationRepo, tokenMaker)
	listingHandler := handler.NewListingHandler(listingService)

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", listingHandler.Health)

		public := v1.Group("")
		{
			public.GET("/listings", listingHandler.ListActiveListings)
			public.GET("/listings/:id", listingHandler.GetActiveListing)
			public.GET("/provinces", listingHandler.ListProvinces)
			public.GET("/districts", listingHandler.ListDistrictsByProvince)
			public.GET("/wards", listingHandler.ListWardsByDistrict)
		}

		protected := v1.Group("").Use(middleware.AuthMiddleware(tokenMaker))
		{
			protected.POST("/listings", listingHandler.CreateListing)
			protected.POST("/listings/:id/publish", listingHandler.PublishListing)
			protected.POST("/listings/:id/deactivate", listingHandler.DeactivateListing)
			protected.POST("/listings/:id/reactivate", listingHandler.ReactivateListing)
			protected.PATCH("/listings/:id/basic-info", listingHandler.UpdateListingBasicInfo)
			protected.PATCH("/listings/:id/address", listingHandler.UpdateListingAddress)
			protected.DELETE("/listings/:id", listingHandler.DeleteListing)

			protected.GET("/me/listings", listingHandler.ListHostListings)
			protected.GET("/me/listings/:id", listingHandler.GetUserListing)
		}
	}

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting server on :%s", addr)
	if err = router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
