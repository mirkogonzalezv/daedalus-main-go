package main

import (
	"daedalus-engine-go/internal/config"
	"daedalus-engine-go/internal/core/logger"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	// Cargar .env solo en entornos no production
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	loadEnv(env)

	logger.Init(env)
	log := logger.L()
	log.Info("Daedalus Engine iniciando...", zap.String("env", env))

	cfg, err := config.CargarVariables()

	if err != nil {
		log.Fatal("Error al cargar configuración", zap.Error(err))
	}

	log.Info("Configuración cargada correctamente")

	_ = cfg

}

func loadEnv(env string) {
	switch env {
	case "development", "dev":
		fmt.Println("Cargando .env.dev...")
		_ = godotenv.Load(".env.dev")
	case "qa":
		fmt.Println("Cargando .env.qa...")
		_ = godotenv.Load(".env.qa")
	case "production", "prod":
		fmt.Println("Usando variables de entorno del sistema")
	default:
		fmt.Printf("Entorno desconocido '%s' usando .env.dev por defecto\n", env)
		_ = godotenv.Load(".env.dev")
	}
}
