package domain

import (
	"time"

	"github.com/google/uuid"
)

type Visitor struct {
	ID        uuid.UUID `json:"id" db:"id"`
	IP        string    `json:"ip" db:"ip"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	City      string    `json:"city" db:"city"`
	Country   string    `json:"country" db:"country"`
	Time      time.Time `json:"time" db:"time"`
}
