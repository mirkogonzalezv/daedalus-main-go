package local

import (
	"context"
	domain "daedalus-engine-go/cmd/core/domain/entities"
	"daedalus-engine-go/cmd/core/domain/repository"
	"database/sql"
	"encoding/json"
)

type PostgresAuditlogRepository struct {
	db *sql.DB
}

func NewAuditlogRepository(db *sql.DB) repository.AuditLogRepository {
	return &PostgresAuditlogRepository{db: db}
}

// Insert implements repository.AuditLogRepository.
func (p *PostgresAuditlogRepository) Insert(ctx context.Context, log *domain.AuditLog) error {
	var metaJSON []byte
	var err error

	if log.Meta != nil {
		metaJSON, err = json.Marshal(log.Meta)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO daedalus.audit_logs(
		id, tenant_id, user_id, action, meta, ip, created_at
		) VALUES ($1,$2,$3,$4,$5,$6, NOW())
	`
	_, err = p.db.ExecContext(ctx, query, log.ID, log.TenantId, log.UserId, log.Action, metaJSON, log.IP)
	return err
}
