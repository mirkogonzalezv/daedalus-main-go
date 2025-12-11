package responses

import "time"

type TenantResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Plan      string    `json:"plan"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTenantResponse struct {
	Tenant  TenantResponse `json:"tenant"`
	Message string         `json:"message"`
}

type GetTenantResponse struct {
	Tenant TenantResponse `json:"tenant"`
}
