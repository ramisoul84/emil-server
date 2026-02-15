package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/ramisoul84/emil-server/config"
)

func CORSMiddleware(cfg *config.Config) fiber.Handler {
	// Debug: Print what's being configured
	fmt.Println("🔧 CORS Config:")
	fmt.Println("  EnableCORS:", cfg.Server.EnableCORS)
	fmt.Println("  AllowOrigins:", strings.Join(cfg.Server.CORSAllowedOrigins, ","))
	fmt.Println("  AllowMethods:", strings.Join(cfg.Server.CORSAllowedMethods, ","))
	fmt.Println("  AllowHeaders:", strings.Join(cfg.Server.CORSAllowedHeaders, ","))
	fmt.Println("  AllowCredentials:", cfg.Server.AllowCredentials)

	if !cfg.Server.EnableCORS {
		return func(c *fiber.Ctx) error {
			fmt.Println("⚠️ CORS is disabled")
			return c.Next()
		}
	}

	// Create CORS config
	corsConfig := cors.Config{
		AllowOrigins:     strings.Join(cfg.Server.CORSAllowedOrigins, ","),
		AllowMethods:     strings.Join(cfg.Server.CORSAllowedMethods, ","),
		AllowHeaders:     strings.Join(cfg.Server.CORSAllowedHeaders, ","),
		ExposeHeaders:    strings.Join(cfg.Server.ExposeHeaders, ","),
		AllowCredentials: cfg.Server.AllowCredentials,
		MaxAge:           cfg.Server.MaxAge,
	}

	// Return CORS middleware
	return cors.New(corsConfig)
}
