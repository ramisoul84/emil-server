package domain

type VisitStartData struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Referrer  string `json:"referrer"`
	UserAgent string `json:"user_agent"`
}

type VisitData struct {
	SessionID string         `json:"session_id"`
	UserID    string         `json:"user_id"`
	Referrer  string         `json:"referrer"`
	UserAgent string         `json:"user_agent"`
	StartTime string         `json:"start_time"`
	Duration  float64        `json:"duration"`
	Actions   map[string]int `json:"actions"`
}

type Data struct {
	ID             int     `json:"id"`
	SessionID      string  `json:"session_id" db:"session_id"`
	UserID         string  `json:"user_id" db:"user_id"`
	IP             string  `json:"ip"`
	Country        string  `json:"country"`
	City           string  `json:"city"`
	OS             string  `json:"os"`
	StartTime      string  `json:"start_time" db:"start_time"`
	Duration       float64 `json:"duration"`
	ActiveDuration float64 `json:"active_duration" db:"active_duration"`
	ActionsCount   int     `json:"actions_count" db:"actions_count"`
}

type Stats struct {
	TotalVisits       int     `json:"total_visits" db:"total_visits"`
	UniqueUsers       int     `json:"unique_users" db:"unique_users"`
	AvgDuration       float64 `json:"avg_duration" db:"avg_duration"`
	AvgActiveDuration float64 `json:"avg_active_duration" db:"avg_active_duration"`
	AvgActions        float64 `json:"avg_actions" db:"avg_actions"`
}
