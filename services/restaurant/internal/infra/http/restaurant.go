package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/katatrina/food-delivery/services/restaurant/internal/service"
)

type RestaurantHandler struct {
	restaurantSvc *service.RestaurantService
}

func NewRestaurantHandler(restaurantSvc *service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{restaurantSvc: restaurantSvc}
}

func (h *RestaurantHandler) CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var cmd service.CreateRestaurantCmd
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	restaurant, err := h.restaurantSvc.CreateRestaurant(context.Background(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(restaurant)
}
