package domain

import "time"

type Session struct {
	ID         string
	TenantID   string
	UserID     string
	RefreshJWT string // hash/identifier del refresh token
	ExpiresAt  time.Time
	IssueAt    time.Time
	IP         *string
	UserAgent  *string
	Revoked    bool
}
