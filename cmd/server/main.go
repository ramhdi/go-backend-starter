package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-backend-starter/internal/api/handlers"
	"go-backend-starter/internal/api/routes"
	"go-backend-starter/internal/config"
	"go-backend-starter/internal/db/postgres"
	"go-backend-starter/internal/repository"
	"go-backend-starter/internal/service"
	"go-backend-starter/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Configure logger
	utils.ConfigureLogger(cfg.Server.Environment)

	// Log successful configuration load
	log.Info().Msg("Configuration loaded successfully")
	log.Info().Str("environment", cfg.Server.Environment).
		Int("port", cfg.Server.Port).
		Str("db_host", cfg.Database.Host).
		Int("db_port", cfg.Database.Port).
		Msg("Application configuration")

	// Set up database connection
	log.Info().Msg("Connecting to database...")
	db, err := postgres.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()
	log.Info().Msg("Database connection established")

	// Initialize layers
	repo := repository.NewPostgresRepository(db.Pool)
	srvc := service.NewService(repo, cfg.JWT.Secret, cfg.JWT.Expiration)
	handler := handlers.NewHandler(srvc)

	// Set up Gin router
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())

	// Set up routes
	routes.Setup(router, handler, srvc)

	// Create server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Int("port", cfg.Server.Port).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
