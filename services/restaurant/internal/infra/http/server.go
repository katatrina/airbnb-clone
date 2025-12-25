package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router            *gin.Engine
	RestaurantHandler *RestaurantHandler
}

func NewServer(restaurantHandler *RestaurantHandler) *Server {
	s := &Server{
		router:            gin.Default(),
		RestaurantHandler: restaurantHandler,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", s.healthCheck)
	s.router.POST("/restaurants", s.RestaurantHandler.CreateRestaurant)
}

func (s *Server) healthCheck(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}

func (s *Server) Router() *gin.Engine {
	return s.router
}
