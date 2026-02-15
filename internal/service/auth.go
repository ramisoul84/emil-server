package service

import (
	"context"
	"fmt"

	"github.com/ramisoul84/emil-server/config"
	"github.com/ramisoul84/emil-server/internal/domain"
	"github.com/ramisoul84/emil-server/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	email          string
	hashedPassword string
	logger         logger.Logger
}

func NewAuthService(cfg *config.Config) *authService {
	return &authService{
		email:          cfg.Security.Email,
		hashedPassword: cfg.Security.HashedPassword,
		logger:         logger.Get(),
	}
}

func (s *authService) Login(ctx context.Context, email, password string) error {
	logger := s.logger.WithFields(
		map[string]any{
			"layer":  "auth_service",
			"method": "login",
		},
	)
	logger.Info().Msg("Login")

	if email != s.email {
		logger.Warn().Msg("Login attempt with wrong email")
		return domain.ErrInvalidCredentials
	}

	err := bcrypt.CompareHashAndPassword([]byte(s.hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			logger.Warn().Msg("Invalid password attempt")
			return domain.ErrInvalidCredentials
		}
		logger.Error().Err(err).Msg("Failed to verify password")
		return fmt.Errorf("failed to verify password: %w", err)
	}

	logger.Info().Msg("Admin logged in successfully")
	return nil
}
