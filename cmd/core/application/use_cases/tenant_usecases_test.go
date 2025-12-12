package usecases

import (
	"context"
	"daedalus-engine-go/cmd/common/logger"
	domain "daedalus-engine-go/cmd/core/domain/entities"
	"testing"
)

type mockTenantRepository struct {
	tenants map[string]*domain.Tenant
}

func newMockTenantRepository() *mockTenantRepository {
	return &mockTenantRepository{
		tenants: make(map[string]*domain.Tenant),
	}
}

func (m *mockTenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	m.tenants[tenant.ID] = tenant
	return nil
}

func (m *mockTenantRepository) GetById(ctx context.Context, id string) (*domain.Tenant, error) {
	tenant, exists := m.tenants[id]
	if !exists {
		return nil, nil
	}
	return tenant, nil
}

func (m *mockTenantRepository) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	for _, tenant := range m.tenants {
		if tenant.Slug == slug {
			return tenant, nil
		}
	}
	return nil, nil
}

func (m *mockTenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	m.tenants[tenant.ID] = tenant
	return nil
}

func TestCreateTenant_Success(t *testing.T) {
	logger.Init("test")
	// Arrange (Preparar)
	mockRepo := newMockTenantRepository()
	useCase := NewTenantUseCase(mockRepo)

	ctx := context.Background()

	// Ejecutar
	tenant, err := useCase.CreateTenant(ctx, "Test Company", "test-company", "pro")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if tenant == nil {
		t.Fatal("Expected tenant to be created, got nil")
	}

	if tenant.Name != "Test Company" {
		t.Errorf("Excepted name 'Test Company', got '%s'", tenant.Name)
	}

	if tenant.Slug != "test-company" {
		t.Errorf("Expected slug 'test-company', got '%s'", tenant.Slug)
	}

	if tenant.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", tenant.Status)
	}
}
