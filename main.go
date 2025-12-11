package main

import (
	"daedalus-engine-go/internal/common/database"
	"daedalus-engine-go/internal/config"
	usecases "daedalus-engine-go/internal/core/application/use_cases"
	"daedalus-engine-go/internal/core/infraestructure/controllers"
	"daedalus-engine-go/internal/core/infraestructure/repository/local"
	"daedalus-engine-go/internal/core/infraestructure/routes"
	"daedalus-engine-go/internal/core/logger"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	// Cargar .env solo en entornos no production
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	config.LoadEnv(env)

	logger.Init(env)
	log := logger.L()
	log.Info("Daedalus Engine iniciando...", zap.String("env", env))

	cfg, err := config.CargarVariables()

	if err != nil {
		log.Fatal("Error al cargar configuración", zap.Error(err))
	}

	log.Info("Configuración cargada correctamente")

	_ = cfg

	migrationDB, err := database.NuevaBaseDeDatos(cfg)

	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos", zap.Error(err))
	}

	migrationCfg := database.MigrationConfig{
		MigratiosPath: cfg.MigrationPath,
		Logger:        log,
	}

	log.Info("Verificando y ejecutando migraciones ...", zap.String("path", cfg.MigrationPath))

	if err := database.RunMigrations(migrationDB, migrationCfg); err != nil {
		log.Fatal("Error ejecutando migraciones", zap.Error(err))
	}

	migrationDB.Close()

	db, err := database.NuevaBaseDeDatos(cfg)
	if err != nil {
		log.Fatal("No se puede conectar a la base de datos", zap.Error(err))
	}

	defer db.Close()

	log.Info("Conexión a base de datos establecida")

	// Migraciones completadas exitosamente

	// Repositorios
	tenantRepo := local.NewTenantRepository(db)

	// UseCases
	tenantUseCases := usecases.NewTenantUseCase(tenantRepo)

	// Controller
	tenantController := controllers.NewTenantController(tenantUseCases)

	router := gin.Default()

	apiRouter := routes.NewAPIRouter(tenantController)

	apiRouter.RegisterRouter(router)

	log.Info("Daedalus Engine iniciado correctamente... ✅")
	router.Run(":" + strconv.Itoa(cfg.Port))

}
