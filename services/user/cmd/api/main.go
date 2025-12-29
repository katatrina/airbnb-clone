package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/katatrina/airbnb-clone/services/user/config"
	"github.com/katatrina/airbnb-clone/services/user/internal/handler"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	h := handler.NewHandler()
	http.HandleFunc("/health", h.Health)

	log.Printf("User service starting on :%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.ServerPort), nil))
}
