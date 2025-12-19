package services

import (
	"context"
	rateLimitEntity "daedalus-engine-go/cmd/internal/features/ratelimit/domain/entitites"
	rateLimitRepository "daedalus-engine-go/cmd/internal/features/ratelimit/domain/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RateLimitService struct {
	repo rateLimitRepository.RateLimitRepository
	log  *zap.Logger
}

type RateLimitResult struct {
	Allowed    bool
	Remaining  int
	ResetAt    time.Time
	RetryAfter time.Duration
}

func NewRateLimitService(repo rateLimitRepository.RateLimitRepository, log *zap.Logger) *RateLimitService {
	return &RateLimitService{
		repo: repo,
		log:  log,
	}
}

func (s *RateLimitService) CheckRateLimit(ctx context.Context, userID, tenantID, endpoint, role, ip, userAgent string) (*RateLimitResult, error) {
	// Obtnemos las regla por endpoint y regla
	rule, err := s.repo.GetRuleByEndpointAndRole(ctx, endpoint, role)
	if err != nil || rule == nil {
		return &RateLimitResult{Allowed: true}, nil
	}

	if !rule.Enabled {
		return &RateLimitResult{Allowed: true}, nil
	}

	// Generamos rate limit key
	key := s.generateKey(userID, tenantID, endpoint)

	rateLimit, err := s.repo.GetRateLimit(ctx, key)
	if err != nil {
		s.log.Error("Falló al obtener el rate limit", zap.Error(err))
		return &RateLimitResult{Allowed: true}, nil
	}

	now := time.Now()

	if rateLimit == nil || now.After(rateLimit.ResetAt) {
		// Creamos una nueva ventana
		rateLimit = &rateLimitEntity.RateLimit{
			Key:          key,
			RequestCount: 1,
			WindowStart:  now,
			WindowSize:   rule.WindowSize,
			MaxRequests:  rule.MaxRequests,
			ResetAt:      now.Add(rule.WindowSize),
		}

		if err := s.repo.SetRateLimit(ctx, rateLimit); err != nil {
			s.log.Error("Falló al configurar el rate limit", zap.Error(err))
		}

		return &RateLimitResult{
			Allowed:   true,
			Remaining: rule.MaxRequests - 1,
			ResetAt:   rateLimit.ResetAt,
		}, nil
	}

	if rateLimit.RequestCount >= rule.MaxRequests {
		// Registro de violation
		violation := &rateLimitEntity.RateLimitViolaton{
			ID:        uuid.NewString(),
			UserID:    userID,
			TenantID:  &tenantID,
			Endpoint:  endpoint,
			IP:        ip,
			UserAgent: userAgent,
			Timestamp: now,
			Attempts:  rateLimit.RequestCount + 1,
		}

		s.repo.RecordViolation(ctx, violation)

		s.log.Warn("Rate limit excedido",
			zap.String("user_id", userID),
			zap.String("endpoint", endpoint),
			zap.String("ip", ip))

		return &RateLimitResult{
			Allowed:    false,
			Remaining:  0,
			ResetAt:    rateLimit.ResetAt,
			RetryAfter: time.Until(rateLimit.ResetAt),
		}, nil
	}

	// Incremantamos el contador
	rateLimit.RequestCount++
	s.repo.SetRateLimit(ctx, rateLimit)

	return &RateLimitResult{
		Allowed:   true,
		Remaining: rule.MaxRequests - rateLimit.RequestCount,
		ResetAt:   rateLimit.ResetAt,
	}, nil
}

func (s *RateLimitService) CheckRateLimitByIP(ctx context.Context, ip, endpoint, userAgent string) (*RateLimitResult, error) {
	maxRequest := 5
	windowSize := time.Minute

	key := fmt.Sprintf("ip:%s:endpoint:%s", ip, endpoint)

	rateLimit, err := s.repo.GetRateLimit(ctx, key)
	if err != nil {
		s.log.Error("Falló al obtener el rate limit", zap.Error(err))
		return &RateLimitResult{Allowed: true}, nil
	}

	now := time.Now()

	if rateLimit == nil || now.After(rateLimit.ResetAt) {
		rateLimit = &rateLimitEntity.RateLimit{
			Key:          key,
			RequestCount: 1,
			WindowStart:  now,
			WindowSize:   windowSize,
			MaxRequests:  maxRequest,
			ResetAt:      now.Add(windowSize),
		}

		if err := s.repo.SetRateLimit(ctx, rateLimit); err != nil {
			s.log.Error("Falló al configurar el rate limit", zap.Error(err))
		}

		return &RateLimitResult{
			Allowed:   true,
			Remaining: maxRequest - 1,
			ResetAt:   rateLimit.ResetAt,
		}, nil
	}

	if rateLimit.RequestCount >= maxRequest {
		violation := &rateLimitEntity.RateLimitViolaton{
			ID:        uuid.NewString(),
			UserID:    "anonymous",
			TenantID:  nil,
			Endpoint:  endpoint,
			IP:        ip,
			UserAgent: userAgent,
			Timestamp: now,
			Attempts:  rateLimit.RequestCount + 1,
		}

		s.repo.RecordViolation(ctx, violation)

		s.log.Warn("Rate limit por IP excedido",
			zap.String("endpoint", endpoint),
			zap.String("ip", ip))

		return &RateLimitResult{
			Allowed:    false,
			Remaining:  0,
			ResetAt:    rateLimit.ResetAt,
			RetryAfter: time.Until(rateLimit.ResetAt),
		}, nil
	}

	// Incrementamos el contador
	rateLimit.RequestCount++
	s.repo.SetRateLimit(ctx, rateLimit)

	return &RateLimitResult{
		Allowed:   true,
		Remaining: maxRequest - rateLimit.RequestCount,
		ResetAt:   rateLimit.ResetAt,
	}, nil
}

// Funcion privada del servicio generateKey
func (s *RateLimitService) generateKey(userID, tenantID, endpoint string) string {
	if tenantID != "" && tenantID != "null" {
		return fmt.Sprintf("tenant:%s:user:%s:endpoint:%s", tenantID, userID, endpoint)
	}
	return fmt.Sprintf("global:user:%s:endpoint:%s", userID, endpoint)
}
