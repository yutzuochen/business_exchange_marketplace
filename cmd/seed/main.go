package main

import (
	"log"
	"trade_company/internal/config"
	"trade_company/internal/database"
)

func main() {
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

	// Run seed data
	log.Println("Starting database seeding...")
	if err := database.SeedData(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}
