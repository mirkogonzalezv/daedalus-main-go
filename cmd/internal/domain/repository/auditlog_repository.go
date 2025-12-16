package repository

import (
	"context"
	domain "daedalus-engine-go/cmd/internal/domain/entities"
)

type AuditLogRepository interface {
	Insert(ctx context.Context, log *domain.AuditLog) error
}
