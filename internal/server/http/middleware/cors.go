package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ramisoul84/emil-server/config"
)

func CORSMiddleware(cfg *config.Config) fiber.Handler {
	if !cfg.Server.EnableCORS {
		return func(c *fiber.Ctx) error {
			fmt.Println("⚠️ CORS is disabled")
			return c.Next()
		}
	}

	allowOrigins := strings.Join(cfg.Server.CORSAllowedOrigins, ",")
	allowMethods := strings.Join(cfg.Server.CORSAllowedMethods, ",")
	allowHeaders := strings.Join(cfg.Server.CORSAllowedHeaders, ",")
	exposeHeaders := strings.Join(cfg.Server.ExposeHeaders, ",")

	config := cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: cfg.Server.AllowCredentials,
		ExposeHeaders:    exposeHeaders,
		MaxAge:           cfg.Server.MaxAge,
	}

	return cors.New(config)
}
