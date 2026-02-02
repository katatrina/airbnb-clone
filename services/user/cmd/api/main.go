package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katatrina/airbnb-clone/pkg/token"
	"github.com/katatrina/airbnb-clone/services/user/config"
	"github.com/katatrina/airbnb-clone/services/user/internal/handler"
	"github.com/katatrina/airbnb-clone/services/user/internal/middleware"
	"github.com/katatrina/airbnb-clone/services/user/internal/repository"
	"github.com/katatrina/airbnb-clone/services/user/internal/service"
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
		log.Fatalf("Failed create token maker: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, tokenMaker)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", userHandler.Health)

		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		users := v1.Group("/me")
		users.Use(middleware.AuthMiddleware(tokenMaker))
		{
			users.GET("", userHandler.GetMe)
		}
	}

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting server on :%s", addr)
	if err = router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
