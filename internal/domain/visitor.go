package domain

import (
	"time"

	"github.com/google/uuid"
)

type Visitor struct {
	ID      uuid.UUID `json:"id" db:"id"`
	UserID  string    `json:"user_id" db:"user_id"`
	IP      string    `json:"ip" db:"ip"`
	OS      string    `json:"os" db:"os"`
	City    string    `json:"city" db:"city"`
	Country string    `json:"country" db:"country"`
	Time    time.Time `json:"time" db:"time"`
}
