package container

import (
	"daedalus-engine-go/internal/common/logger"
	usecases "daedalus-engine-go/internal/core/application/use_cases"
	"daedalus-engine-go/internal/core/infraestructure/controllers"
	"daedalus-engine-go/internal/core/infraestructure/repository/local"
	"database/sql"
)

type Container struct {
	TenantController *controllers.TenantController
}

func NewContainer(db *sql.DB) *Container {
	log := logger.L()
	// Repositories
	tenantRepo := local.NewTenantRepository(db)

	// Use Cases
	tenantUseCase := usecases.NewTenantUseCase(tenantRepo)

	// Controllers
	tenantController := controllers.NewTenantController(tenantUseCase, log)

	return &Container{
		TenantController: tenantController,
	}
}
