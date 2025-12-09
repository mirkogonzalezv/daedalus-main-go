package main

import (
	"daedalus-engine-go/internal/config"
	"daedalus-engine-go/internal/core/database"
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

}
