package repository

import (
	"context"
	rateLimitEntity "daedalus-engine-go/cmd/internal/features/ratelimit/domain/entitites"
)

type RateLimitRepository interface {
	GetRateLimit(ctx context.Context, key string) (*rateLimitEntity.RateLimit, error)
	SetRateLimit(ctx context.Context, rateLimit *rateLimitEntity.RateLimit) error
	GetRuleByEndpointAndRole(ctx context.Context, endpoint, role string) (*rateLimitEntity.RateLimitRule, error)
	RecordViolation(ctx context.Context, violation *rateLimitEntity.RateLimitViolaton) error
}
