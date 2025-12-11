package container

import (
	usecases "daedalus-engine-go/internal/core/application/use_cases"
	"daedalus-engine-go/internal/core/infraestructure/controllers"
	"daedalus-engine-go/internal/core/infraestructure/repository/local"
	"database/sql"
)

type Container struct {
	TenantController *controllers.TenantController
}

func NewContainer(db *sql.DB) *Container {
	// Repositories
	tenantRepo := local.NewTenantRepository(db)

	// Use Cases
	tenantUseCase := usecases.NewTenantUseCase(tenantRepo)

	// Controllers
	tenantController := controllers.NewTenantController(tenantUseCase)

	return &Container{
		TenantController: tenantController,
	}
}
