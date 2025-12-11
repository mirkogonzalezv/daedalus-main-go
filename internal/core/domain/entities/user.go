package domain

import "time"

type User struct {
	ID           string
	TenantID     string // uuid fk
	Name         string
	Email        string
	PasswordHash string
	Role         string // owner | admin | user
	Status       string // active | disabled
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
