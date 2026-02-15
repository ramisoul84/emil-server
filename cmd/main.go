package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ramisoul84/emil-server/config"
	"github.com/ramisoul84/emil-server/internal/repository"
	"github.com/ramisoul84/emil-server/internal/server/bot"
	"github.com/ramisoul84/emil-server/internal/server/http"
	"github.com/ramisoul84/emil-server/internal/server/http/handler"
	"github.com/ramisoul84/emil-server/internal/service"
	"github.com/ramisoul84/emil-server/internal/storage/postgres"
	"github.com/ramisoul84/emil-server/pkg/jwt"
	"github.com/ramisoul84/emil-server/pkg/logger"
)

func main() {
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	cfg, err := config.Load(env)
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// ==================== Logger ====================
	logger.InitGlobal(cfg)
	logger.Info().
		Int("PID", os.Getpid()).
		Msg("Starting Emil Server")

	// ==================== Database ====================
	db, err := postgres.New(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database.")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error().Err(err).Msg("Error closing database.")
		}
	}()

	logger.Info().Msg("Postgres connection established")

	// ==================== Telegram Bot ====================
	botServer, err := bot.NewBotServer(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create bot")
	}

	// ====================  Repository ====================
	analyticsRepository := repository.NewAnalyticsRepository(db)
	messageRepository := repository.NewMessageRepository(db)

	// ==================== Services ====================
	botService := service.NewBotService(botServer)
	analyticsService := service.NewAnalyticsService(analyticsRepository, botService)
	authService := service.NewAuthService(cfg)
	messageService := service.NewMessageService(messageRepository, botService)
	jwt := jwt.NewJWT(cfg)

	// ==================== Handler ====================
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	authHandler := handler.NewAuthHandler(authService, jwt)
	messageHandler := handler.NewMessageHandler(messageService)

	// ==================== HTTP Server ====================
	srv := http.NewServer(cfg, analyticsHandler, authHandler, messageHandler)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info().Msg("Received shutdown signal")
		if err := srv.Shutdown(); err != nil {
			logger.Error().Err(err).Msg("Error during server shutdown")
		}
	}()

	if err := srv.Start(); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}

}
