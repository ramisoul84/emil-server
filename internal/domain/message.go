package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID      uuid.UUID `json:"id" db:"id"`
	Name    string    `json:"name" db:"name"`
	Email   string    `json:"email" db:"email"`
	Text    string    `json:"text" db:"text"`
	Time    time.Time `json:"time" db:"time"`
	Unread  bool      `json:"unread" db:"unread"`
	Country string    `json:"country" db:"country"`
}

type CreateMessageRequest struct {
	Name  string `json:"name" validate:"required,email"`
	Email string `json:"email" validate:"required"`
	Text  string `json:"text" validate:"required"`
}
