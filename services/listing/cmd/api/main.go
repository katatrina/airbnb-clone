package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katatrina/airbnb-clone/pkg/middleware"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/pkg/token"
	"github.com/katatrina/airbnb-clone/services/listing/config"
	"github.com/katatrina/airbnb-clone/services/listing/internal/handler"
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

	router.NoRoute(func(c *gin.Context) {
		response.NotFound(c, response.CodeRouteNotFound, "The requested endpoint does not exist")
	})

	router.GET("/health", listingHandler.Health)

	v1 := router.Group("/api/v1")
	{
		public := v1.Group("")
		{
			public.GET("/listings", listingHandler.ListActiveListings)
			public.GET("/listings/:id", listingHandler.GetActiveListing)
			public.GET("/provinces", listingHandler.ListProvinces)
			public.GET("/provinces/:code/districts", listingHandler.ListDistrictsByProvince)
			public.GET("/districts/:code/wards", listingHandler.ListWardsByDistrict)
		}

		hostListings := v1.Group("/me/listings")
		hostListings.Use(middleware.AuthMiddleware(tokenMaker))
		{
			hostListings.POST("", listingHandler.CreateListing)
			hostListings.GET("", listingHandler.ListHostListings)
			hostListings.GET("/:id", listingHandler.GetHostListing)
			hostListings.PATCH("/:id/basic-info", listingHandler.UpdateListingBasicInfo)
			hostListings.PATCH("/:id/address", listingHandler.UpdateListingAddress)
			hostListings.DELETE("/:id", listingHandler.DeleteListing)
			hostListings.POST("/:id/publish", listingHandler.PublishListing)
			hostListings.POST("/:id/deactivate", listingHandler.DeactivateListing)
			hostListings.POST("/:id/reactivate", listingHandler.ReactivateListing)
		}
	}

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting server on :%s", addr)
	if err = router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
