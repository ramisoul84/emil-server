package domain

import "time"

type Message struct {
	ID      int       `json:"id" db:"id"`
	UserID  string    `json:"user_id" db:"user_id"`
	Name    string    `json:"name" db:"name"`
	Email   string    `json:"email" db:"email"`
	Text    string    `json:"text" db:"text"`
	Time    time.Time `json:"time" db:"time"`
	Unread  bool      `json:"unread" db:"unread"`
	IP      string    `json:"ip"`
	City    string    `json:"city" db:"city"`
	Country string    `json:"country" db:"country"`
}
