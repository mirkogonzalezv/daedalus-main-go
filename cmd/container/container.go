package container

import (
	"daedalus-engine-go/cmd/common/logger"
	usecases "daedalus-engine-go/cmd/core/application/use_cases"
	"daedalus-engine-go/cmd/core/infraestructure/controllers"
	"daedalus-engine-go/cmd/core/infraestructure/repository/local"
	"database/sql"
)

type Container struct {
	TenantController *controllers.TenantController
	UserController   *controllers.UserController
}

func NewContainer(db *sql.DB) *Container {
	log := logger.L()
	// Repositories
	tenantRepo := local.NewTenantRepository(db)

	// Use Cases
	tenantUseCase := usecases.NewTenantUseCase(tenantRepo)

	// Controllers
	tenantController := controllers.NewTenantController(tenantUseCase, log)

	userRepo := local.NewUserRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo, log)
	userController := controllers.NewUserController(userUseCase, log)

	return &Container{
		TenantController: tenantController,
		UserController:   userController,
	}
}
