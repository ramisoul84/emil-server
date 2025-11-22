package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/ramisoul/emil-server/internal/domain"
	"gitlab.com/ramisoul/emil-server/pkg/logger"
)

type AuthService interface {
	Login(email, password string) (string, string, error)
}

type authHandler struct {
	service AuthService
	log     logger.Logger
}

func NewAuthHandler(service AuthService, log logger.Logger) *authHandler {
	return &authHandler{service, log}
}

func (h *authHandler) Login(c echo.Context) error {
	log := h.log.WithFields(map[string]any{
		"layer":     "handlers",
		"operation": "login",
	})

	if c.Request().Body == nil {
		log.Info("Empty request body")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Request body is required",
		})
	}

	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		log.WithError(err).Info("Invalid request body")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format: " + err.Error(),
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email and password are required",
		})
	}

	accessToken, refreshToken, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid credentials",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
