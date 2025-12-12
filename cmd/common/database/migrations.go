package database

import (
	"context"
	"database/sql"
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
func RunMigrations(db *sql.DB, cfg MigrationConfig) error {
	log := cfg.Logger

	log.Info("Verificando existencia del schema daedalus...")

	_, err := db.ExecContext(context.Background(), "CREATE SCHEMA IF NOT EXISTS daedalus")
	if err != nil {
		log.Error("Error creando schema daedalus", zap.Error(err))
		return fmt.Errorf("error creating schema daedalus: %w", err)
	}
	log.Info("Schema daedalus verificado/creado correctamente")

	driver, err := postgres.WithInstance(db, &postgres.Config{
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

func RollbackMigration(db *sql.DB, cfg MigrationConfig, steps int) error {
	log := cfg.Logger

	driver, err := postgres.WithInstance(db, &postgres.Config{
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
