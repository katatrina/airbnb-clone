package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katatrina/airbnb-clone/services/user/config"
	"github.com/katatrina/airbnb-clone/services/user/internal/handler"
	"github.com/katatrina/airbnb-clone/services/user/internal/middleware"
	"github.com/katatrina/airbnb-clone/services/user/internal/repository"
	"github.com/katatrina/airbnb-clone/services/user/internal/service"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create db pool: %v", err)
	}
	defer db.Close()

	if err = db.Ping(ctx); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, cfg)
	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.GET("/health", userHandler.Health)
	v1.POST("/auth/register", userHandler.Register)
	v1.POST("/auth/login", userHandler.Login)

	v1.GET("/users/me", middleware.AuthMiddleware(cfg.JWTSecret), userHandler.GetMe)

	log.Fatal(router.Run(fmt.Sprintf(":%s", cfg.ServerPort)))
}
