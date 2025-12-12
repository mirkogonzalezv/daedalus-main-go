package repository

import (
	"context"
	domain "daedalus-engine-go/cmd/core/domain/entities"
)

type AuditLogRepository interface {
	Insert(ctx context.Context, log *domain.AuditLog) error
}
