package database

import (
	"daedalus-engine-go/cmd/common/logger"
	"daedalus-engine-go/cmd/config"
	infra "daedalus-engine-go/cmd/infra/database"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

func NuevaBaseDeDatos(cfg *config.Config) (*sql.DB, error) {

	log := logger.L()

	switch cfg.DBType {
	case "postgres":
		log.Info("Conectando a PostgreSQL...")
		pool, err := infra.NuevaConexionPostgres(
			infra.PostgresConfig{
				Host:     cfg.DBHost,
				Port:     cfg.DBPort,
				User:     cfg.DBUser,
				Password: cfg.DBPassword,
				DBName:   cfg.DBName,
				SSLMode:  cfg.DBSSLMode,
			})
		if err != nil {
			log.Error("Error al conectar con PostgresSQL", zap.Error(err))
			return nil, err
		}

		return pool, nil
	default:
		err := fmt.Errorf("motor no soportado: %s", cfg.DBType)
		log.Error("motor de base de datos no soportado", zap.Error(err))
		return nil, fmt.Errorf("motor de base de datos no soportado: %s", cfg.DBType)
	}

}
