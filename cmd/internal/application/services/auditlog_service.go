package services

import (
	"context"
	domain "daedalus-engine-go/cmd/internal/domain/entities"
	"daedalus-engine-go/cmd/internal/domain/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuditService struct {
	repo repository.AuditLogRepository
	log  *zap.Logger
}

func NewAuditService(repo repository.AuditLogRepository, log *zap.Logger) *AuditService {
	return &AuditService{
		repo: repo,
		log:  log,
	}
}

func (s *AuditService) LogAction(ctx context.Context, tenantID, userID, action, ip string, meta map[string]interface{}) {
	var tenantIDPtr *string
	var userIDPtr *string
	var ipPtr *string

	// TenantID puede ser opcional (*string)
	if tenantID != "" {
		tenantIDPtr = &tenantID
	}

	// UserID puede ser opcional (*string)
	if userID != "" {
		userIDPtr = &userID
	}

	// IP puede ser opcional (*string)
	if ip != "" {
		ipPtr = &ip
	}

	auditLog := &domain.AuditLog{
		ID:       uuid.NewString(),
		TenantId: tenantIDPtr,
		UserId:   userIDPtr,
		Action:   action,
		Meta:     meta,
		IP:       ipPtr,
	}

	go func() {
		if err := s.repo.Insert(context.Background(), auditLog); err != nil {
			s.log.Error("Failed to insert audit log", zap.Error(err))
		}
	}()
}
