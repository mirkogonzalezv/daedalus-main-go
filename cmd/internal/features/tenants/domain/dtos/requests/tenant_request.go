package requests

type CreateTenantRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	Plan string `json:"plan"`
}

type GetTenantByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type GetTenantBySlugRequest struct {
	Slug string `uri:"slug" binding:"required"`
}

// DTO para actualizar el tenant
type UpdateTenantRequest struct {
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
	Plan string `json:"plan,omitempty"`
}
