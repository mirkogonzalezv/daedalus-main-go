package domain

import "time"

type Conversation struct {
	ID        string
	TenantId  string
	UserId    string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
