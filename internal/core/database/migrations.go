package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

type MigrationConfig struct {
	MigratiosPath string
	Logger        *zap.Logger
}

// RunMigrations ejecuta las migraciones pendientes
func RunMigrations(db Database, cfg MigrationConfig) error {
	log := cfg.Logger

	stdDB, err := db.GetStdlibDB()

	if err != nil {
		log.Error("Error obteniendo conexión stdlib", zap.Error(err))
		return fmt.Errorf("error getting stdlib connection: %w", err)
	}
	defer stdDB.Close()

	log.Info("Verificando existencia del schema daedalus...")

	_, err = stdDB.ExecContext(context.Background(), "CREATE SCHEMA IF NOT EXISTS daedalus")
	if err != nil {
		log.Error("Error creando schema daedalus", zap.Error(err))
		return fmt.Errorf("error creating schema daedalus: %w", err)
	}
	log.Info("Schema daedalus verificado/creado correctamente")

	driver, err := postgres.WithInstance(stdDB, &postgres.Config{
		MigrationsTable: "schema_migrations",
		SchemaName:      "daedalus",
	})
	if err != nil {
		log.Error("Error creando driver de migraciones", zap.Error(err))
		return fmt.Errorf("error creating migration driver: %w", err)
	}

	// Crear instancia de migrate
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigratiosPath),
		"postgres",
		driver,
	)
	if err != nil {
		log.Error("Error creando instancia de migrate", zap.Error(err))
		return fmt.Errorf("error creating migration instance: %w", err)
	}
	defer m.Close()

	// Ejecutamos migraciones
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("No hay migraciones pendientes")
			return nil
		}

		log.Error("Error ejecutando migraciones", zap.Error(err))
		return fmt.Errorf("error running migrations: %w", err)
	}

	log.Info("Migraciones aplicadas exitosamente")

	return nil
}

func RollbackMigration(db Database, cfg MigrationConfig, steps int) error {
	log := cfg.Logger

	stdDB, err := db.GetStdlibDB()
	if err != nil {
		log.Error("Error obteniendo conexión stdlib", zap.Error(err))
		return fmt.Errorf("error getting stdlib connection: %w", err)
	}
	defer stdDB.Close()

	driver, err := postgres.WithInstance(stdDB, &postgres.Config{
		MigrationsTable: "schema_migrations",
		SchemaName:      "daedalus",
	})
	if err != nil {
		log.Error("Error creando driver de migraciones", zap.Error(err))
		return fmt.Errorf("error creating migration driver: %w", err)
	}

	// Crear instancia de migrate
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigratiosPath),
		"postgres",
		driver,
	)
	if err != nil {
		log.Error("Error creando instancia de migrate", zap.Error(err))
		return fmt.Errorf("error creating migration instance: %w", err)
	}
	defer m.Close()

	// Hacemos rollback
	if err := m.Steps(-steps); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("No hay migraciones para hacer rollback")
			return nil
		}
		log.Error("Error haciendo rollback", zap.Error(err))
		return fmt.Errorf("error rolling back migrations: %w", err)
	}

	log.Info("Rollback completado", zap.Int("steps", steps))
	return nil
}

func GetMigrationVersion(db Database, cfg MigrationConfig) (uint, bool, error) {
	log := cfg.Logger

	stdDB, err := db.GetStdlibDB()
	if err != nil {
		log.Error("Error obteniendo conexión stdlib", zap.Error(err))
		return 0, false, fmt.Errorf("error getting stdlib connection: %w", err)
	}
	defer stdDB.Close()

	driver, err := postgres.WithInstance(stdDB, &postgres.Config{
		MigrationsTable: "schema_migrations",
		SchemaName:      "daedalus",
	})
	if err != nil {
		log.Error("Error creando driver de migraciones", zap.Error(err))
		return 0, false, fmt.Errorf("error creating migration driver: %w", err)
	}

	// Crear instancia de migrate
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigratiosPath),
		"postgres",
		driver,
	)
	if err != nil {
		log.Error("Error creando instancia de migrate", zap.Error(err))
		return 0, false, fmt.Errorf("error creating migration instance: %w", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			log.Info("No hay migraciones aplicadas")
			return 0, false, nil
		}
		log.Error("Error obteniendo versión", zap.Error(err))
		return 0, false, fmt.Errorf("error getting migration version: %w", err)
	}

	log.Info("Version de migración atual", zap.Uint("version", version), zap.Bool("dirty", dirty))

	return version, dirty, nil
}
