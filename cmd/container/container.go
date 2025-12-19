package container

import (
	"daedalus-engine-go/cmd/common/logger"
	"daedalus-engine-go/cmd/config"
	auditRepository "daedalus-engine-go/cmd/internal/features/auditlog/application/data/local"
	tenantRepository "daedalus-engine-go/cmd/internal/features/tenants/application/data/local"
	tenantUsecases "daedalus-engine-go/cmd/internal/features/tenants/application/use_cases"
	tenantController "daedalus-engine-go/cmd/internal/features/tenants/infraestructure/handlers"
	userRepository "daedalus-engine-go/cmd/internal/features/users/application/data/local"
	userUseCases "daedalus-engine-go/cmd/internal/features/users/application/use_cases"
	userController "daedalus-engine-go/cmd/internal/features/users/infraestructure/handlers"
	"time"

	authRepository "daedalus-engine-go/cmd/internal/features/auth/application/data/local"
	authUseCases "daedalus-engine-go/cmd/internal/features/auth/application/uses_cases"
	authController "daedalus-engine-go/cmd/internal/features/auth/infraestructure/handlers"
	authService "daedalus-engine-go/cmd/internal/features/auth/infraestructure/services"
	"daedalus-engine-go/cmd/internal/pkg/services"
	"database/sql"

	rateLimitRepository "daedalus-engine-go/cmd/internal/features/ratelimit/application/data/local"
	rateLimitService "daedalus-engine-go/cmd/internal/features/ratelimit/infraestructure/services"
)

type Container struct {
	AuthController   *authController.AuthController
	TenantController *tenantController.TenantController
	UserController   *userController.UserController
	AuthService      *authService.AuthService
	RateLimitService *rateLimitService.RateLimitService
}

func NewContainer(db *sql.DB, cfg *config.Config) *Container {
	log := logger.L()

	// Audit
	auditLogRepo := auditRepository.NewAuditlogRepository(db)
	auditSvc := services.NewAuditService(auditLogRepo, log)

	// Tenant
	tenantRepo := tenantRepository.NewTenantRepository(db)
	tenantUseCase := tenantUsecases.NewTenantUseCase(tenantRepo, log, auditSvc)
	tenantController := tenantController.NewTenantController(tenantUseCase, log)

	// User
	userRepo := userRepository.NewUserRepository(db)
	userUseCase := userUseCases.NewUserUseCase(userRepo, log)
	userController := userController.NewUserController(userUseCase, log)

	// Auth
	authRepo := authRepository.NewAuthRepository(db)
	// TODO: Inyectar secret JWT por variables de entorno
	authService := authService.NewAuthService(cfg.JwtSecret, +time.Duration(cfg.JwtExpiredMin)*time.Minute,
		time.Duration(cfg.JwtRefreshExpireDays)*time.Hour*24,
		log)
	authUseCase := authUseCases.NewAuthUseCase(authRepo, userRepo, authService, log)
	authController := authController.NewAuthController(authUseCase, log)

	// Ratelimit
	rateLimitRepo := rateLimitRepository.NewMemoryRateLimitRepository()
	rateLimitSvc := rateLimitService.NewRateLimitService(rateLimitRepo, log)

	return &Container{
		TenantController: tenantController,
		UserController:   userController,
		AuthController:   authController,
		AuthService:      authService,
		RateLimitService: rateLimitSvc,
	}
}
