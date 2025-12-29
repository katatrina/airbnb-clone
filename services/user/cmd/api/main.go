package main

import (
	"log"
	"net/http"

	"github.com/katatrina/airbnb-clone/services/user/internal/handler"
)

func main() {
	h := handler.NewHandler()
	http.HandleFunc("/health", h.Health)

	log.Println("User service starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
