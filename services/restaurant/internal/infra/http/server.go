package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router            *chi.Mux
	RestaurantHandler *RestaurantHandler
}

func NewServer(restaurantHandler *RestaurantHandler) *Server {
	s := &Server{
		router:            chi.NewRouter(),
		RestaurantHandler: restaurantHandler,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Get("/health", s.healthCheck)
	s.router.Post("/restaurants", s.RestaurantHandler.CreateRestaurant)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (s *Server) Router() *chi.Mux {
	return s.router
}
