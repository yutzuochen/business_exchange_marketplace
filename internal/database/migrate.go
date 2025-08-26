package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"trade_company/internal/config"

	"github.com/golang-migrate/migrate/v4"
	migrateMySQL "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// RunMigrations runs database migrations using golang-migrate
func RunMigrations(db *gorm.DB) error {
	// Create a separate database connection for migrations to avoid conflicts
	// Load config to get DSN
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config for migrations: %w", err)
	}
	dsn := cfg.MySQLDSN()
	if !strings.Contains(dsn, "multiStatements=") {
		if strings.Contains(dsn, "?") {
			dsn += "&multiStatements=true"
		} else {
			dsn += "?multiStatements=true"
		}
	}
	// migrationDB, err := sql.Open("mysql", dsn)
	migrationDB, err := sql.Open("mysql", migrationDSN(cfg))
	// Create a separate database connection for migrations
	// migrationDB, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		return fmt.Errorf("failed to open migration database: %w", err)
	}
	defer migrationDB.Close()

	// Test the migration connection
	if err := migrationDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping migration database: %w", err)
	}

	// Create MySQL driver instance with separate connection
	driver, err := migrateMySQL.WithInstance(migrationDB, &migrateMySQL.Config{})
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

	// Run migrations
	log.Println("Running database migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			log.Printf("Warning: failed to close migrate instance on error - src: %v, db: %v", srcErr, dbErr)
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Close migrate instance after successful migration
	if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
		log.Printf("Warning: failed to close migrate instance - src: %v, db: %v", srcErr, dbErr)
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
	// Load config to get DSN
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config for migrations: %w", err)
	}

	// Create a separate database connection for migrations
	// migrationDB, err := sql.Open("mysql", cfg.MySQLDSN())
	migrationDB, err := sql.Open("mysql", migrationDSN(cfg))
	if err != nil {
		return fmt.Errorf("failed to open migration database: %w", err)
	}
	defer migrationDB.Close()

	// Create MySQL driver instance
	driver, err := migrateMySQL.WithInstance(migrationDB, &migrateMySQL.Config{})
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
	defer func() {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			log.Printf("Warning: failed to close migrate instance - src: %v, db: %v", srcErr, dbErr)
		}
	}()

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
	// Load config to get DSN
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config for migrations: %w", err)
	}

	// Create a separate database connection for migrations
	migrationDB, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		return fmt.Errorf("failed to open migration database: %w", err)
	}
	defer migrationDB.Close()

	// Create MySQL driver instance
	driver, err := migrateMySQL.WithInstance(migrationDB, &migrateMySQL.Config{})
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
	defer func() {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			log.Printf("Warning: failed to close migrate instance - src: %v, db: %v", srcErr, dbErr)
		}
	}()

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
	// Load config to get DSN
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config for migrations: %w", err)
	}

	// Create a separate database connection for migrations
	migrationDB, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		return fmt.Errorf("failed to open migration database: %w", err)
	}
	defer migrationDB.Close()

	// Create MySQL driver instance
	driver, err := migrateMySQL.WithInstance(migrationDB, &migrateMySQL.Config{})
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
	defer func() {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			log.Printf("Warning: failed to close migrate instance - src: %v, db: %v", srcErr, dbErr)
		}
	}()

	// Force version
	if err := m.Force(version); err != nil {
		return fmt.Errorf("failed to force version: %w", err)
	}

	log.Printf("Successfully forced migration version to %d", version)
	return nil
}

func migrationDSN(cfg *config.Config) string {
	dsn := cfg.MySQLDSN()
	if !strings.Contains(dsn, "multiStatements=") {
		if strings.Contains(dsn, "?") {
			dsn += "&multiStatements=true"
		} else {
			dsn += "?multiStatements=true"
		}
	}
	return dsn
}
