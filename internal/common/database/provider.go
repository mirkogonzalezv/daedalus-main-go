package database

import (
	"daedalus-engine-go/internal/config"
	"daedalus-engine-go/internal/core/logger"
	infra "daedalus-engine-go/internal/infra/database"
	"fmt"

	"go.uber.org/zap"
)

func NuevaBaseDeDatos(cfg *config.Config) (Database, error) {

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

		return &PostgresDB{pool: pool}, nil
	case "other_db":
		err := fmt.Errorf("other_db aun no implementado")
		log.Warn("Intentando usar una abse de datos no implementada", zap.Error(err))
		return nil, fmt.Errorf("other_db aun no implementado, puedes configurar Mysql , Mongo u otro")
	default:
		err := fmt.Errorf("motor no soportado: %s", cfg.DBType)
		log.Error("motor de base de datos no soportado", zap.Error(err))
		return nil, fmt.Errorf("motor de base de datos no soportado: %s", cfg.DBType)
	}

}
