package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ramisoul84/emil-server/config"
)

func CORSMiddleware(cfg *config.Config) fiber.Handler {
	if !cfg.Server.EnableCORS {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.Server.CORSAllowedOrigins, ","),
		AllowMethods:     strings.Join(cfg.Server.CORSAllowedMethods, ","),
		AllowHeaders:     strings.Join(cfg.Server.CORSAllowedHeaders, ","),
		ExposeHeaders:    strings.Join(cfg.Server.ExposeHeaders, ","),
		AllowCredentials: cfg.Server.AllowCredentials,
		MaxAge:           cfg.Server.MaxAge,
	})
}
