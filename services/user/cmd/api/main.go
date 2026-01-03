package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/services/user/config"
	"github.com/katatrina/airbnb-clone/services/user/internal/handler"
	"github.com/katatrina/airbnb-clone/services/user/internal/middleware"
	"github.com/katatrina/airbnb-clone/services/user/internal/repository"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	db, err := repository.NewUserRepository(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	log.Println("connected to db")

	h := handler.NewUserHandler(db, cfg)
	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.GET("/health", h.Health)
	v1.POST("/auth/register", h.Register)
	v1.POST("/auth/login", h.Login)

	v1.GET("/users/me", middleware.AuthMiddleware(cfg.JWTSecret), h.GetMe)

	log.Fatal(router.Run(fmt.Sprintf(":%s", cfg.ServerPort)))
}
