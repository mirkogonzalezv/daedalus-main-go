package usecases

import (
	"context"
	logger "daedalus-engine-go/cmd/common/logger"
	domain "daedalus-engine-go/cmd/core/domain/entities"
	domainErrors "daedalus-engine-go/cmd/core/domain/errors"
	"daedalus-engine-go/cmd/core/domain/repository"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
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
	log := logger.L()

	if name == "" {
		return nil, domainErrors.ErrTenantNameRequired()
	}

	if slug == "" {
		return nil, domainErrors.ErrTenantInvalidSlug()
	}

	if plan == "" {
		return nil, domainErrors.ErrTenantCodeInvalidPlan()
	}

	// Normalizar slug a logwecase y trim espaces
	slug = strings.ToLower(strings.TrimSpace(slug))

	// validar slug único
	log.Info("Verificando slug único", zap.String("slug", slug))
	existe, err := uc.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if existe != nil {
		return nil, domainErrors.ErrTenantSlugExist()
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

	tenant, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if tenant == nil {
		return nil, domainErrors.ErrTenantNotFound()
	}

	return tenant, nil
}

func (uc *TenantUseCase) ObtenerTenantPorSlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	if slug == "" {
		return nil, domainErrors.ErrTenantInvalidSlug()
	}

	// Normalizamos el slug a buscar
	slug = strings.ToLower(strings.TrimSpace(slug))

	tenant, err := uc.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if tenant == nil {
		return nil, domainErrors.ErrTenantNotFound()
	}

	return tenant, nil
}

// Actualizar Tenant
func (uc *TenantUseCase) ActualizarTenant(ctx context.Context, t *domain.Tenant) error {
	if t == nil {
		return domainErrors.ErrTenantNameRequired()
	}

	if t.ID == "" {
		return domainErrors.ErrTenantIdIsRequired()
	}

	// Validamos que el slug único (si cambió)
	if t.Slug != "" {
		existe, err := uc.repo.GetBySlug(ctx, t.Slug)
		if err != nil {
			return err
		}

		if existe != nil && existe.ID != t.ID {
			return domainErrors.ErrTenantSlugExist()
		}
	}

	return uc.repo.Update(ctx, t)
}
