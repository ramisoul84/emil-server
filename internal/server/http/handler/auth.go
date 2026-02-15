package handler

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ramisoul84/emil-server/internal/domain"
	"github.com/ramisoul84/emil-server/pkg/logger"
)

type authService interface {
	Login(ctx context.Context, email, password string) error
}

type JWT interface {
	GenerateAccessToken(email string) (string, error)
}

type authHandler struct {
	service authService
	jwt     JWT
	logger  logger.Logger
}

func NewAuthHandler(service authService, jwt JWT) *authHandler {
	return &authHandler{
		service: service,
		jwt:     jwt,
		logger:  logger.Get(),
	}
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to parse login request")
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request body",
			Code:    fiber.StatusBadRequest,
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ErrorResponse{
			Error:   "validation_error",
			Message: "Email and password are required",
			Code:    fiber.StatusBadRequest,
		})
	}

	err := h.service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			h.logger.Warn().Msgf("Failed login attempt for email: %s", req.Email)
			return c.Status(fiber.StatusUnauthorized).JSON(domain.ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Invalid email or password",
			})
		default:
			h.logger.Error().Err(err).Msg("Login service error")
			return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{
				Error:   "server_error",
				Message: "An internal error occurred",
			})
		}
	}
	fmt.Println("GENERATING")

	accessToken, err := h.jwt.GenerateAccessToken(req.Email)
	fmt.Println(accessToken, err)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to generate token")
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate authentication token",
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.LoginResponse{
		Token: accessToken,
	})
}
