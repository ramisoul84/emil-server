package hash

import (
	"fmt"

	"gitlab.com/ramisoul/emil-server/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type hash struct {
	cost int
}

func NewHash(cost int) *hash {
	return &hash{cost}
}

func (h *hash) HashPassword(password string) (string, error) {
	if password == "" {
		return "", domain.ErrBadRequest
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

func (h *hash) VerifyHash(hashedPassword, plainPassword string) error {
	if hashedPassword == "" || plainPassword == "" {
		return domain.ErrBadRequest
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return domain.ErrInvalidCredentials
		}
		return fmt.Errorf("failed to verify password %w", err)
	}
	return nil
}
