package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katatrina/food-delivery/services/restaurant/config"
	httpserver "github.com/katatrina/food-delivery/services/restaurant/internal/infra/http"
	"github.com/katatrina/food-delivery/services/restaurant/internal/infra/postgres"
	"github.com/katatrina/food-delivery/services/restaurant/internal/service"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	RunE:  runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to create db connection pool: %w", err)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	store := postgres.NewStore(pool)
	restaurantSvc := service.NewRestaurantService(store)
	restaurantHandler := httpserver.NewRestaurantHandler(restaurantSvc)
	server := httpserver.NewServer(restaurantHandler)

	log.Println("Starting server on :8080")
	return http.ListenAndServe(":8080", server.Router())
}
