package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnv(env string) {
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
