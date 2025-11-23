package middlewares

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.com/ramisoul/emil-server/pkg/logger"
)

type Tokener interface {
	ValidateAccessToken(tokenString string) error
}

func AuthMiddleware(jwt Tokener, log logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/api/v1/auth/login" || c.Path() == "/api/v1/auth/refresh" {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header required",
				})
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
				})
			}

			token := parts[1]
			err := jwt.ValidateAccessToken(token)
			if err != nil {
				log.WithError(err).Warn("Invalid token attempt")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
			}

			return next(c)
		}
	}
}
