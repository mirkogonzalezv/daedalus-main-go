package domain

import "time"

type AuditLog struct {
	ID        string
	TenantId  string
	UserId    *string
	Action    string
	Meta      map[string]any //jsonb
	IP        *string
	CreatedAt time.Time
}
