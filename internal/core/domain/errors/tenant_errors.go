package errors

const (
	TenantNameRequired          = "TENANT_001"
	TenantNameNotFound          = "TENANT_002"
	TenantSlugNotFound          = "TENANT_004"
	TenantCodeSlugExist         = "TENANT_004"
	TenantCodeInvalidSlug       = "TENANT_005"
	TenantCodeInvalidPlan       = "TENANT_006"
	TenantCodeInvalidOptionPlan = "TENANT_007"
	TenantIdIsRequired          = "TENANT_008"
)

func ErrTenantNameRequired() *DomainError {
	return NewValidationError(TenantNameNotFound, "Nombre del tenant es requerido")
}

func ErrTenantIdIsRequired() *DomainError {
	return NewValidationError(TenantIdIsRequired, "El id del tenant es requerido")
}

func ErrTenantNotFound() *DomainError {
	return NewValidationError(TenantNameNotFound, "Tenant no encontrado")
}

func ErrTenantSlugNotFound() *DomainError {
	return NewValidationError(TenantSlugNotFound, "Slug del tenant es requerido")
}

func ErrTenantSlugExist() *DomainError {
	return NewValidationError(TenantCodeSlugExist, "Slug del tenant ya existe")
}

func ErrTenantInvalidSlug() *DomainError {
	return NewValidationError(TenantCodeInvalidSlug, "Slug del tenant es requerido")
}

func ErrTenantCodeInvalidPlan() *DomainError {
	return NewValidationError(TenantCodeInvalidPlan, "Plan del tenant es requerido")
}

func ErrTenantInvalidOptionPlan() *DomainError {
	return NewValidationError(TenantCodeInvalidOptionPlan, "Plan del tenant es inválido")
}
