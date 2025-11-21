package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/ramisoul/emil-server/config"
	"gitlab.com/ramisoul/emil-server/internal/domain"
)

type JWT struct {
	jwtSecret          []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func New(cfg config.JWTConfig) *JWT {
	return &JWT{
		jwtSecret:          []byte(cfg.JWTSecret),
		accessTokenExpiry:  cfg.AccessTokenExpiry,
		refreshTokenExpiry: cfg.RefreshTokenExpiry,
	}
}

func (t *JWT) GenerateAccessToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.accessTokenExpiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	return token.SignedString(t.jwtSecret)
}

func (t *JWT) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (t *JWT) ValidateAccessToken(tokenString string) error {
	if tokenString == "" {
		return domain.ErrInvalidAccessToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return domain.ErrExpiredAccessToken
		}
		return domain.ErrInvalidAccessToken
	}

	if !token.Valid {
		return domain.ErrInvalidAccessToken
	}

	return nil
}
