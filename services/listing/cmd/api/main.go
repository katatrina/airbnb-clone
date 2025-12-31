package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/services/listing/config"
	"github.com/katatrina/airbnb-clone/services/listing/internal/database"
	"github.com/katatrina/airbnb-clone/services/listing/internal/handler"
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
	v1 := router.Group("/api/v1")
	v1.GET("/health", h.Health)
	v1.GET("/provinces", h.ListProvinces)

	log.Fatal(router.Run(fmt.Sprintf(":%s", cfg.ServerPort)))
}
