package entitites

import "time"

type RateLimit struct {
	Key          string
	RequestCount int
	WindowStart  time.Time
	WindowSize   time.Duration
	MaxRequests  int
	ResetAt      time.Time
}

type RateLimitRule struct {
	ID          string
	Name        string
	Endpoint    string
	Role        string
	MaxRequests int
	WindowSize  time.Duration
	Enabled     bool
}

type RateLimitViolaton struct {
	ID        string
	UserID    string
	TenantID  *string
	Endpoint  string
	IP        string
	UserAgent string
	Timestamp time.Time
	Attempts  int
}
