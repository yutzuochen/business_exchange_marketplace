package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"

	"trade_company/internal/config"
	"trade_company/internal/database"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Parse command line flags
	var (
		action  = flag.String("action", "up", "Migration action: up, down, status, force")
		version = flag.Int("version", 0, "Version to force (for force action)")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg, nil)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Execute migration action
	switch *action {
	case "up":
		if err := database.RunMigrations(db); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations completed successfully")

	case "down":
		if err := database.RollbackMigrations(db); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		log.Println("Migration rollback completed successfully")

	case "status":
		if err := database.GetMigrationStatus(db); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

	case "force":
		if err := database.ForceVersion(db, *version); err != nil {
			log.Fatalf("Failed to force version: %v", err)
		}
		log.Printf("Forced version to %d", *version)

	default:
		log.Fatalf("Unknown action: %s. Use: up, down, status, or force", *action)
	}
}
