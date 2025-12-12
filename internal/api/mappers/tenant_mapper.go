package mappers

import (
	"daedalus-engine-go/internal/api/dtos/responses"
	domain "daedalus-engine-go/internal/core/domain/entities"
)

func ToTenantResponse(tenant *domain.Tenant) responses.TenantResponse {
	return responses.TenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Slug:      tenant.Slug,
		Plan:      tenant.SubscriptionPlan,
		Status:    tenant.Status,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
	}
}

func ToCreateTenantResponse(tenant *domain.Tenant) responses.CreateTenantResponse {
	return responses.CreateTenantResponse{
		Tenant:  ToTenantResponse(tenant),
		Message: "Tenant creado exitosamente",
	}
}

func ToGetTenantResponse(tenant *domain.Tenant) responses.GetTenantResponse {
	return responses.GetTenantResponse{
		Tenant: ToTenantResponse(tenant),
	}
}

func ToUpdateTenantResponse(tenant *domain.Tenant) responses.UpdateTenantResponse {
	return responses.UpdateTenantResponse{
		Tenant:  ToTenantResponse(tenant),
		Message: "Tenant actualizado exitosamente",
	}
}
