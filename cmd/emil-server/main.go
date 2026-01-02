package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/ramisoul/emil-server/config"
	"gitlab.com/ramisoul/emil-server/internal/repository"
	"gitlab.com/ramisoul/emil-server/internal/server/bot"
	"gitlab.com/ramisoul/emil-server/internal/server/http"
	"gitlab.com/ramisoul/emil-server/internal/service"
	"gitlab.com/ramisoul/emil-server/internal/storage/postgres"
	"gitlab.com/ramisoul/emil-server/pkg/hash"
	"gitlab.com/ramisoul/emil-server/pkg/jwt"
	"gitlab.com/ramisoul/emil-server/pkg/logger"
)

func main() {
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	log := logger.New(env, "emil-server", "1.0.0")

	log.WithFields(map[string]any{
		"environment": env,
		"version":     "1.0.0",
	}).Info("Starting Emil Server")

	cfg, err := config.Load(env)
	if err != nil {
		log.WithError(err).Error("Failed to load configuration")
		os.Exit(1)
	}

	log.Info("Configuration loaded successfully")

	// Bot
	bot, err := bot.NewBot(*cfg)
	if err != nil {
		log.WithError(err).Error("Failed to create a bot")
	}

	bot.Start()

	// PostgreSQL Connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := postgres.New(cfg.Database)
	if err != nil {
		log.WithError(err).Error("Failed to connect to PostgreSQL")
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Error("Failed to close database connection")
		} else {
			log.Info("Database connection closed")
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		log.WithError(err).Error("Failed to ping database")
		os.Exit(1)
	}

	log.Info("PostgreSQL connection established and verified")

	// Utilities
	jwt := jwt.New(cfg.JWT)
	hasher := hash.NewHash(10)

	// Repositories
	messageRepository := repository.NewMessageRepository(db, log)
	analyticsRepository := repository.NewAnalyticsRepository(db, log)

	// Services
	authService := service.NewAuthService(jwt, hasher, log)
	messageService := service.NewMessageService(messageRepository, log)
	analyticsService := service.NewAnalyticsService(analyticsRepository, log)

	// HTTP Server
	server := http.New(cfg.Server, log, authService, messageService, analyticsService, bot, jwt)

	go func() {
		if err := server.Start(); err != nil {
			log.WithError(err).Error("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Received shutdown signal")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	log.Info("Initiating graceful shutdown...")

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("Server shutdown failed")
		os.Exit(1)
	}

	log.Info("Server shutdown completed successfully")
}
