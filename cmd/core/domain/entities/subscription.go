package domain

import "time"

type Subscription struct {
	ID                   string
	TenantID             string
	StripeCustomerId     *string
	StripeSubscriptionID *string
	PeriodStart          *time.Time
	PeriodEnd            *time.Time
	Status               string // active | canceled | unpadi
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
