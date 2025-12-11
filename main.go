package main

import (
	"daedalus-engine-go/internal/common/database"
	"daedalus-engine-go/internal/config"
	"daedalus-engine-go/internal/core/logger"
	"os"

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

	db, err := database.NuevaBaseDeDatos(cfg)

	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos", zap.Error(err))
	}

	defer db.Close()

	log.Info("Conexión a base de datos establecida")

	migrationCfg := database.MigrationConfig{
		MigratiosPath: cfg.MigrationPath,
		Logger:        log,
	}

	log.Info("Verificando y ejecutando migraciones...", zap.String("path", cfg.MigrationPath))

	if err := database.RunMigrations(db, migrationCfg); err != nil {
		log.Fatal("Error ejecutando migraciones", zap.Error(err))
	}

	// Verificar versión actual de migraciones
	version, dirty, err := database.GetMigrationVersion(db, migrationCfg)
	if err != nil {
		log.Warn("No se pudo obtener versión de migraciones", zap.Error(err))
	} else {
		log.Info("Estado de migraciones", zap.Uint("version", version), zap.Bool("dirty", dirty))
	}

	log.Info("Daedalus Engine iniciado correctamente... ✅")

}
