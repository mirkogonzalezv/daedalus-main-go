package container

import (
	"daedalus-engine-go/cmd/common/logger"
	auditRepository "daedalus-engine-go/cmd/internal/features/auditlog/application/data/local"
	tenantRepository "daedalus-engine-go/cmd/internal/features/tenants/application/data/local"
	tenantUsecases "daedalus-engine-go/cmd/internal/features/tenants/application/use_cases"
	tenantController "daedalus-engine-go/cmd/internal/features/tenants/infraestructure/handlers"
	userRepository "daedalus-engine-go/cmd/internal/features/users/application/data/local"
	userUseCases "daedalus-engine-go/cmd/internal/features/users/application/use_cases"
	userController "daedalus-engine-go/cmd/internal/features/users/infraestructure/handlers"
	"daedalus-engine-go/cmd/internal/pkg/services"
	"database/sql"
)

type Container struct {
	TenantController *tenantController.TenantController
	UserController   *userController.UserController
}

func NewContainer(db *sql.DB) *Container {
	log := logger.L()

	//audit service
	auditLogRepo := auditRepository.NewAuditlogRepository(db)
	//audit service
	auditSvc := services.NewAuditService(auditLogRepo, log)
	// Repositories
	tenantRepo := tenantRepository.NewTenantRepository(db)
	// Use Cases
	tenantUseCase := tenantUsecases.NewTenantUseCase(tenantRepo, log, auditSvc)
	// Controllers
	tenantController := tenantController.NewTenantController(tenantUseCase, log)

	userRepo := userRepository.NewUserRepository(db)
	userUseCase := userUseCases.NewUserUseCase(userRepo, log)
	userController := userController.NewUserController(userUseCase, log)

	return &Container{
		TenantController: tenantController,
		UserController:   userController,
	}
}
