package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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

	h := handler.NewHandler(db)
	http.HandleFunc("/health", h.Health)
	http.HandleFunc("/auth/register", h.Register)

	log.Printf("User service starting on :%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.ServerPort), nil))
}
