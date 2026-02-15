package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ramisoul84/emil-server/config"
	"github.com/ramisoul84/emil-server/internal/domain"
)

type JWT struct {
	jwtSecret            []byte
	accessTokenExpiresIn time.Duration
}

func NewJWT(cfg *config.Config) *JWT {
	return &JWT{
		jwtSecret:            []byte(cfg.Security.JWTSecret),
		accessTokenExpiresIn: cfg.Security.AccessTokenExpiresIn,
	}
}

func (t *JWT) GenerateAccessToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.accessTokenExpiresIn.Abs())),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	accessToken, err := token.SignedString(t.jwtSecret)
	if err != nil {
		return "", domain.ErrTokenGenerate
	}

	return accessToken, nil
}

func (t *JWT) VerifyAccessToken(accessToken string) error {
	if accessToken == "" {
		return domain.ErrTokenInvalid
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return domain.ErrTokenExpired
		}
		return domain.ErrTokenInvalid
	}

	if !token.Valid {
		return domain.ErrTokenInvalid
	}

	return nil
}
