package tenant

import (
	"context"
	domain "daedalus-engine-go/internal/core/domain/entities"
	"daedalus-engine-go/internal/core/domain/repository"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type TenantUseCase struct {
	repo repository.TenantRepository
}

func NewTenantUseCase(repo repository.TenantRepository) *TenantUseCase {
	return &TenantUseCase{repo: repo}
}

// Aqui es donde aplicaremos la logica de negocio
func (uc *TenantUseCase) CreateTenant(ctx context.Context, name string, slug string, plan string) (*domain.Tenant, error) {
	// Validaciones básicas

	if name == "" {
		return nil, errors.New("tenant name es requerido")
	}

	if slug == "" {
		return nil, errors.New("tenant slug es requerido")
	}

	// Normalizar slug a logwecase y trim espaces
	slug = strings.ToLower(strings.TrimSpace(slug))

	// validar slug único
	existe, err := uc.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if existe != nil {
		return nil, errors.New("tenant slug ya existe")
	}

	// Creamos la instancia domain de tenant
	nuevoTenant := &domain.Tenant{
		ID:               uuid.NewString(),
		Name:             name,
		Slug:             slug,
		SubscriptionPlan: plan,
		Status:           "active", // estado por defecto
	}

	// Guardamos
	err = uc.repo.Create(ctx, nuevoTenant)

	if err != nil {
		return nil, err
	}

	return nuevoTenant, nil
}

func (uc *TenantUseCase) ObtenerTenantPorId(ctx context.Context, id string) (*domain.Tenant, error) {
	if id == "" {
		return nil, errors.New("id es requerido")
	}

	tenant, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if tenant == nil {
		return nil, errors.New("tenant no existe")
	}

	return tenant, nil
}

func (uc *TenantUseCase) ObtenerTenantPorSlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	if slug == "" {
		return nil, errors.New("slug es requerido")
	}

	// Normalizamos el slug a buscar
	slug = strings.ToLower(strings.TrimSpace(slug))

	tenant, err := uc.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if tenant == nil {
		return nil, errors.New("tenant no encontrado")
	}

	return tenant, nil
}

// Actualizar Tenant
func (uc *TenantUseCase) ActualizarTenant(ctx context.Context, t *domain.Tenant) error {
	if t == nil {
		return errors.New("tenant es requerido")
	}

	if t.ID == "" {
		return errors.New("tenant ID es requerido")
	}

	// Validamos que el slug único (si cambió)
	if t.Slug != "" {
		existe, err := uc.repo.GetBySlug(ctx, t.Slug)
		if err != nil {
			return err
		}

		if existe != nil && existe.ID != t.ID {
			return errors.New("tenant slug ya esta en uso")
		}
	}

	return uc.repo.Update(ctx, t)
}
