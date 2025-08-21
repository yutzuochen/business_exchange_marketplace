package database

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// RunMigrations runs database migrations using golang-migrate
func RunMigrations(db *gorm.DB) error {
	// Get the underlying *sql.DB from GORM
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create MySQL driver instance
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create mysql driver: %w", err)
	}

	// Get migrations path
	migrationsPath := "file://migrations"
	if os.Getenv("MIGRATIONS_PATH") != "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "mysql", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	log.Println("Running database migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("Database is up to date, no migrations needed")
	} else {
		log.Println("Database migrations completed successfully")
	}

	return nil
}

// RollbackMigrations rolls back the last migration
func RollbackMigrations(db *gorm.DB) error {
	// Get the underlying *sql.DB from GORM
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create MySQL driver instance
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create mysql driver: %w", err)
	}

	// Get migrations path
	migrationsPath := "file://migrations"
	if os.Getenv("MIGRATIONS_PATH") != "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "mysql", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Rollback last migration
	log.Println("Rolling back last migration...")
	if err := m.Steps(-1); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Migration rollback completed successfully")
	return nil
}

// GetMigrationStatus gets the current migration status
func GetMigrationStatus(db *gorm.DB) error {
	// Get the underlying *sql.DB from GORM
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create MySQL driver instance
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create mysql driver: %w", err)
	}

	// Get migrations path
	migrationsPath := "file://migrations"
	if os.Getenv("MIGRATIONS_PATH") != "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "mysql", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Get current version
	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			log.Println("Migration status: No migrations have been run")
			return nil
		}
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	log.Printf("Migration status: Version %d, Dirty: %t", version, dirty)
	return nil
}

// ForceVersion forces the migration version to a specific version
func ForceVersion(db *gorm.DB, version int) error {
	// Get the underlying *sql.DB from GORM
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create MySQL driver instance
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create mysql driver: %w", err)
	}

	// Get migrations path
	migrationsPath := "file://migrations"
	if os.Getenv("MIGRATIONS_PATH") != "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "mysql", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Force version
	if err := m.Force(version); err != nil {
		return fmt.Errorf("failed to force version: %w", err)
	}

	log.Printf("Successfully forced migration version to %d", version)
	return nil
}
