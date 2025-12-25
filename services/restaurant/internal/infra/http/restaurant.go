package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/food-delivery/services/restaurant/internal/model"
	"github.com/katatrina/food-delivery/services/restaurant/internal/service"
)

type RestaurantService interface {
	CreateRestaurant(ctx context.Context, cmd service.CreateRestaurantCmd) (*model.Restaurant, error)
}

type RestaurantHandler struct {
	restaurantSvc RestaurantService
}

func NewRestaurantHandler(restaurantSvc RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{restaurantSvc: restaurantSvc}
}

func (h *RestaurantHandler) CreateRestaurant(ctx *gin.Context) {
	var cmd service.CreateRestaurantCmd
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurant, err := h.restaurantSvc.CreateRestaurant(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, restaurant)
}
