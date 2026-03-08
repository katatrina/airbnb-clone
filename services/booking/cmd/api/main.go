package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katatrina/airbnb-clone/pkg/token"
	"github.com/katatrina/airbnb-clone/services/booking/config"
	"github.com/katatrina/airbnb-clone/services/booking/internal/client"
	"github.com/katatrina/airbnb-clone/services/booking/internal/handler"
	"github.com/katatrina/airbnb-clone/pkg/middleware"
	"github.com/katatrina/airbnb-clone/services/booking/internal/repository"
	"github.com/katatrina/airbnb-clone/services/booking/internal/service"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	listingClient := client.NewListingClient(cfg.ListingServiceURL)
	bookingRepo := repository.NewBookingRepository(db)
	bookingService := service.NewBookingService(bookingRepo, listingClient)
	bookingHandler := handler.NewBookingHandler(bookingService)

	router := gin.Default()

	router.GET("/health", bookingHandler.Health)

	v1 := router.Group("/api/v1")
	{
		protected := v1.Group("").Use(middleware.AuthMiddleware(tokenMaker))
		{
			protected.POST("/bookings", bookingHandler.CreateBooking)

			protected.GET("/me/bookings", bookingHandler.ListGuestBookings)
			protected.GET("/me/bookings/:id", bookingHandler.GetBooking)
			protected.POST("/me/bookings/:id/cancel", bookingHandler.CancelBooking)

			protected.GET("/me/hosting/bookings", bookingHandler.ListHostBookings)
			protected.POST("/me/hosting/bookings/:id/confirm", bookingHandler.ConfirmBooking)
			protected.POST("/me/hosting/bookings/:id/reject", bookingHandler.RejectBooking)
		}
	}

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Booking service starting on %s", addr)
	if err = router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
