package http

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/ramisoul84/emil-server/config"
	"github.com/ramisoul84/emil-server/internal/server/http/middleware"
	"github.com/ramisoul84/emil-server/pkg/logger"
)

type analyticsHandler interface {
	VisitStart(c *fiber.Ctx) error
	VisitEnd(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Stats(c *fiber.Ctx) error
}

type authHandler interface {
	Login(c *fiber.Ctx) error
}

type messageHandler interface {
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
}

type Server struct {
	app              *fiber.App
	analyticsHandler analyticsHandler
	authHandler      authHandler
	messageHandler   messageHandler
	cfg              *config.Config
	logger           logger.Logger
}

func NewServer(cfg *config.Config, analyticsHandler analyticsHandler, authHandler authHandler, messageHandler messageHandler) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:           cfg.Server.ReadTimeout,
		WriteTimeout:          cfg.Server.WriteTimeout,
		IdleTimeout:           cfg.Server.IdleTimeout,
		DisableStartupMessage: true,
	})

	srv := &Server{
		app:              app,
		analyticsHandler: analyticsHandler,
		authHandler:      authHandler,
		messageHandler:   messageHandler,
		logger:           logger.Get(),
		cfg:              cfg,
	}

	srv.setupMiddlewares()
	srv.setupRoutes()

	return srv
}

func (s *Server) setupMiddlewares() {
	s.app.Use(requestid.New())
	s.app.Use(middleware.CORSMiddleware(s.cfg))
	s.app.Use(middleware.ObservabilityMiddleware(s.logger))

	s.app.Use(recover.New())
}

func (s *Server) setupRoutes() {
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	api := s.app.Group("/api/v1")
	public := api.Group("/")
	public.Post("/analytics/visit-start", s.analyticsHandler.VisitStart)
	public.Post("/analytics/visit-end", s.analyticsHandler.VisitEnd)
	public.Post("/auth/login", s.authHandler.Login)
	public.Post("/message/save", s.messageHandler.Create)

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(s.cfg, s.logger))
	protected.Get("/analytics/list", s.analyticsHandler.List)
	protected.Get("/analytics/stats", s.analyticsHandler.Stats)
	protected.Get("/message/{:id}", s.messageHandler.Get)
	protected.Patch("/message/{:id}", s.messageHandler.Update)
	protected.Delete("/message/{:id}", s.messageHandler.Delete)
	protected.Get("/message", s.messageHandler.List)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	s.logger.Info().
		Str("address", addr).
		Msg("🚀 Server starting")

	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	s.logger.Info().Msg("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.ShutdownTimeout)
	defer cancel()

	return s.app.ShutdownWithContext(ctx)
}
