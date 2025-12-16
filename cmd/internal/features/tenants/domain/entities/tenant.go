package domain

import "time"

type Tenant struct {
	ID               string //uuid
	Name             string
	Slug             string // opcional, para subdominios
	SubscriptionPlan string // free | pro | enterprise
	Status           string // active | suspended
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
