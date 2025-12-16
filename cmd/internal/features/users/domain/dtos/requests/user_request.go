package requests

type CreateUserRequestDTO struct {
	Name     string `json:"name"`
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type GetUserByIdRequestDTO struct {
	ID string `uri:"id" binding:"required"`
}

type GetUserByEmailAndTenantIdDTO struct {
	Email    string `uri:"email" binding:"required"`
	TenantID string `uri:"tenant_id" binding:"required"`
}

// DTO para actualizar el usuario
type UpdateUserRequestDTO struct {
	Name     string `json:"name,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
	Status   string `json:"status,omitempty"`
}
