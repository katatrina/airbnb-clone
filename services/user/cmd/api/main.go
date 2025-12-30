package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/services/user/config"
	"github.com/katatrina/airbnb-clone/services/user/internal/database"
	"github.com/katatrina/airbnb-clone/services/user/internal/handler"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	log.Println("connected to db")

	h := handler.NewHandler(db, cfg)
	router := gin.Default()
	router.GET("/health", h.Health)
	router.POST("/auth/register", h.Register)
	router.POST("/auth/login", h.Login)

	log.Fatal(router.Run(fmt.Sprintf(":%s", cfg.ServerPort)))
}
