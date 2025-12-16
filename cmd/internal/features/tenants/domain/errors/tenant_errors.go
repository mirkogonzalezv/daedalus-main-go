package errors

import (
	domainErrors "daedalus-engine-go/cmd/internal/pkg/errors"
)

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

func ErrTenantNameRequired() *domainErrors.DomainError {
	return domainErrors.NewValidationError(TenantNameNotFound, "Nombre del tenant es requerido")
}

func ErrTenantIdIsRequired() *domainErrors.DomainError {
	return domainErrors.NewValidationError(TenantIdIsRequired, "El id del tenant es requerido")
}

func ErrTenantNotFound() *domainErrors.DomainError {
	return domainErrors.NewValidationError(TenantNameNotFound, "Tenant no encontrado")
}

func ErrTenantSlugNotFound() *domainErrors.DomainError {
	return domainErrors.NewValidationError(TenantSlugNotFound, "Slug del tenant es requerido")
}

func ErrTenantSlugExist() *domainErrors.DomainError {
	return domainErrors.NewValidationError(TenantCodeSlugExist, "Slug del tenant ya existe")
}

func ErrTenantInvalidSlug() *domainErrors.DomainError {
	return domainErrors.NewValidationError(TenantCodeInvalidSlug, "Slug del tenant es requerido")
}

func ErrTenantCodeInvalidPlan() *domainErrors.DomainError {
	return domainErrors.NewValidationError(TenantCodeInvalidPlan, "Plan del tenant es requerido")
}

func ErrTenantInvalidOptionPlan() *domainErrors.DomainError {
	return domainErrors.NewValidationError(TenantCodeInvalidOptionPlan, "Plan del tenant es inválido")
}
