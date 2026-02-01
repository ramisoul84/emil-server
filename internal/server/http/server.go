package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/ramisoul/emil-server/config"
	"gitlab.com/ramisoul/emil-server/internal/domain"
	"gitlab.com/ramisoul/emil-server/internal/server/http/handlers"
	"gitlab.com/ramisoul/emil-server/internal/server/http/middlewares"
	"gitlab.com/ramisoul/emil-server/pkg/logger"
)

type BotService interface {
	BroadcastMessage(string) (int, []error)
}

type AuthService interface {
	Login(email, password string) (string, string, error)
}

type MessageService interface {
	CreateMessage(ctx context.Context, req *domain.CreateMessageRequest) error
	GetMessageByID(ctx context.Context, id uuid.UUID) (*domain.Message, error)
	MarkMessageAsRead(ctx context.Context, id uuid.UUID) error
	DeleteMessage(ctx context.Context, id uuid.UUID) error
	GetMessagesList(ctx context.Context, limit, offset int) ([]*domain.Message, int, error)
}

type AnalyticsService interface {
	SaveVisitor(ctx context.Context, visitor *domain.Visitor) error
	GetVisitors(ctx context.Context, limit, offset int) ([]*domain.Visitor, int, int, error)
	GetCountByUserId(ctx context.Context, userId string) (int, error)
}

type Tokener interface {
	ValidateAccessToken(tokenString string) error
}

type Server struct {
	echo             *echo.Echo
	server           *http.Server
	log              logger.Logger
	authService      AuthService
	messageService   MessageService
	analyticsService AnalyticsService
	bot              BotService
	jwt              Tokener
}

func New(
	cfg config.ServerConfig,
	log logger.Logger,
	authService AuthService,
	messageService MessageService,
	analyticsService AnalyticsService,
	bot BotService,
	jwt Tokener) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"https://emilsuliman.com",
			"https://www.emilsuliman.com",
			"http://localhost:4200",
		},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())

	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server := &Server{e, srv, log, authService, messageService, analyticsService, bot, jwt}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	api := s.echo.Group("/api/v1")

	messageHandler := handlers.NewMessageHandler(s.messageService, s.bot, s.log)
	authHandler := handlers.NewAuthHandler(s.authService, s.log)
	analyticsHandler := handlers.NewAnalyticsHandler(s.analyticsService, s.bot, s.log)

	message := api.Group("/message")

	message.POST("", messageHandler.Create)

	message.Use(middlewares.AuthMiddleware(s.jwt, s.log))

	message.GET("", messageHandler.List)
	message.GET("/:id", messageHandler.Get)
	message.PUT("/:id", messageHandler.Update)
	message.DELETE("/:id", messageHandler.Delete)

	auth := api.Group("/auth")
	auth.POST("/login", authHandler.Login)

	analytics := api.Group("/analytics")
	analytics.POST("/track", analyticsHandler.TrackVisitor)
	analytics.GET("/visitors", analyticsHandler.GetVisitors)

}

func (s *Server) Start() error {
	s.log.WithFields(map[string]any{
		"port": s.server.Addr,
	}).Info("Starting HTTP server")

	s.server.Handler = s.echo

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.log.WithError(err).Error("Server failed to start")
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down HTTP server gracefully")

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.WithError(err).Error("Server shutdown failed")
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	s.log.Info("HTTP server shutdown completed")
	return nil
}
