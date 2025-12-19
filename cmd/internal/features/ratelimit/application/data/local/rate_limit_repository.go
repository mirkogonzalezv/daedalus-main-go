package local

import (
	"context"
	rateLimitEntity "daedalus-engine-go/cmd/internal/features/ratelimit/domain/entitites"
	rateLimitRepository "daedalus-engine-go/cmd/internal/features/ratelimit/domain/repository"
	"sync"
	"time"
)

type MemoryRateLimitRepository struct {
	rateLimits sync.Map
	rules      map[string]*rateLimitEntity.RateLimitRule
	violations []*rateLimitEntity.RateLimitViolaton
	rulesMutex sync.RWMutex
	violMutex  sync.Mutex
}

func NewMemoryRateLimitRepository() rateLimitRepository.RateLimitRepository {
	repo := &MemoryRateLimitRepository{
		rules:      make(map[string]*rateLimitEntity.RateLimitRule),
		violations: make([]*rateLimitEntity.RateLimitViolaton, 0),
	}

	repo.initialiceDefaultRules()
	return repo
}

// GetRateLimit implements [repository.RateLimitRepository].
func (r *MemoryRateLimitRepository) GetRateLimit(ctx context.Context, key string) (*rateLimitEntity.RateLimit, error) {
	values, exists := r.rateLimits.Load(key)
	if !exists {
		return nil, nil
	}

	rateLimit := values.(*rateLimitEntity.RateLimit)
	return rateLimit, nil
}

// GetRuleByEndpointAndRole implements [repository.RateLimitRepository].
func (r *MemoryRateLimitRepository) GetRuleByEndpointAndRole(ctx context.Context, endpoint string, role string) (*rateLimitEntity.RateLimitRule, error) {
	r.rulesMutex.RLock()
	defer r.rulesMutex.RUnlock()

	key := endpoint + ":" + role
	rule, exists := r.rules[key]
	if !exists {
		defaultKey := endpoint + ":default"
		rule, exists = r.rules[defaultKey]
		if !exists {
			return nil, nil
		}
	}

	return rule, nil
}

// RecordViolation implements [repository.RateLimitRepository].
func (r *MemoryRateLimitRepository) RecordViolation(ctx context.Context, violation *rateLimitEntity.RateLimitViolaton) error {
	r.violMutex.Lock()
	defer r.violMutex.Unlock()

	r.violations = append(r.violations, violation)
	return nil
}

// SetRateLimit implements [repository.RateLimitRepository].
func (r *MemoryRateLimitRepository) SetRateLimit(ctx context.Context, rateLimit *rateLimitEntity.RateLimit) error {
	r.rateLimits.Store(rateLimit.Key, rateLimit)
	return nil
}

func (r *MemoryRateLimitRepository) initialiceDefaultRules() {
	rules := []*rateLimitEntity.RateLimitRule{
		{
			ID:          "api_default_user",
			Name:        "API Default User",
			Endpoint:    "/api",
			Role:        "user",
			MaxRequests: 100,
			WindowSize:  time.Minute,
			Enabled:     true,
		},
		{
			ID:          "api_default_admin",
			Name:        "API Default Admin",
			Endpoint:    "/api",
			Role:        "admin",
			MaxRequests: 500,
			WindowSize:  time.Minute,
			Enabled:     true,
		},
		{
			ID:          "api_default_system_admin",
			Name:        "API Default System Admin",
			Endpoint:    "/api",
			Role:        "system_admin",
			MaxRequests: 1000,
			WindowSize:  time.Minute,
			Enabled:     true,
		},
	}

	for _, rule := range rules {
		key := rule.Endpoint + ":" + rule.Role
		r.rules[key] = rule
	}
}
