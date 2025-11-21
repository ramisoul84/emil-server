package service

import (
	"gitlab.com/ramisoul/emil-server/internal/domain"
	"gitlab.com/ramisoul/emil-server/pkg/logger"
)

type Tokener interface {
	GenerateAccessToken(email string) (string, error)
	GenerateRefreshToken() (string, error)
}

type Hahser interface {
	VerifyHash(hashedPassword, plainPassword string) error
}

type authService struct {
	jwt  Tokener
	hash Hahser
	log  logger.Logger
}

func NewAuthService(jwt Tokener, hash Hahser, log logger.Logger) *authService {
	return &authService{jwt, hash, log}
}

func (s *authService) Login(email, password string) (string, string, error) {
	log := s.log.WithFields(map[string]any{
		"layer":     "service",
		"operation": "login",
	})

	if email != "rami@mail.ru" {
		return "", "", domain.ErrInvalidCredentials
	}

	err := s.hash.VerifyHash("$2a$10$4P8JPBTHRD/otMj25/d5YO8mpDdxfaCvw.TybZwtAsOIgbinpwWxK", password)
	if err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, err := s.jwt.GenerateAccessToken(email)
	if err != nil {
		log.WithError(err).Error("faild to generate access token")
		return "", "", domain.ErrInternalServer
	}

	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		log.WithError(err).Error("faild to generate refresh token")
		return "", "", domain.ErrInternalServer
	}
	return accessToken, refreshToken, nil
}
