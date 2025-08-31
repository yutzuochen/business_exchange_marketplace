// Package main serves as the entry point for the Business Exchange Marketplace backend service.
// This service provides REST APIs and GraphQL endpoints for managing business listings,
// user authentication, messaging, favorites, transactions, and leads.
//
// Key features:
// - JWT-based authentication with Redis session management
// - MySQL database with automatic migrations and seeding
// - Redis caching for session management and performance
// - RESTful API endpoints for business operations
// - GraphQL API for flexible data queries
// - Graceful shutdown handling
// - Retry logic for database connections
// - Structured logging with Zap
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"trade_company/internal/config"
	"trade_company/internal/database"
	"trade_company/internal/logger"
	"trade_company/internal/models"
	"trade_company/internal/redisclient"
	"trade_company/internal/router"

	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// main is the application entry point that initializes all services and starts the HTTP server.
// It performs the following initialization sequence:
// 1. Load environment configuration from .env file
// 2. Initialize structured logging
// 3. Connect to MySQL database with retry logic
// 4. Run database migrations and seed initial data
// 5. Connect to Redis for caching (optional)
// 6. Initialize HTTP router with middleware
// 7. Start HTTP server with graceful shutdown support
func main() {
	fmt.Println("========= Business Exchange Marketplace Starting =================")
	
	// Load environment variables from .env file (development/testing only)
	_ = godotenv.Load()

	// Load application configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize structured logger (Zap) based on environment
	zapLogger := logger.New(cfg.AppEnv)
	defer zapLogger.Sync() // Flush any buffered log entries on exit

	// Database Connection with Retry Logic
	// Attempt to connect to MySQL database with exponential backoff
	// The service can start without database connection for health checks
	var db *gorm.DB
	dbRetryCount := 0
	maxDbRetries := 5

	zapLogger.Info("Attempting to connect to database...")
	for dbRetryCount < maxDbRetries {
		db, err = database.Connect(cfg, zapLogger)
		if err == nil {
			zapLogger.Info("Database connection established successfully")
			break
		}

		dbRetryCount++
		zapLogger.Sugar().Warnw("Database connection failed, retrying...", 
			"error", err, 
			"attempt", dbRetryCount,
			"max_retries", maxDbRetries)

		// Exponential backoff: wait 1s, 2s, 3s, 4s, 5s between retries
		if dbRetryCount < maxDbRetries {
			time.Sleep(time.Duration(dbRetryCount) * time.Second)
		}
	}

	// Database initialization (migrations and seeding)
	// Service can function without database for basic health checks
	if db == nil {
		zapLogger.Error("Unable to connect to database after retries, continuing without database")
	} else {
		zapLogger.Info("Running database migrations...")

		// Apply database schema migrations to ensure tables are up-to-date
		if err := database.RunMigrations(db); err != nil {
			zapLogger.Error("Database migrations failed", logger.Err(err))
		} else {
			zapLogger.Info("Database migrations completed successfully")
		}

		// Seed initial data (users, sample listings, etc.) for development/testing
		zapLogger.Info("Seeding initial database data...")
		if err := database.SeedData(db, cfg); err != nil {
			zapLogger.Error("Database seeding failed", logger.Err(err))
		} else {
			zapLogger.Info("Database seeding completed successfully")
		}
	}

	// Redis Connection (Optional)
	// Redis is used for session management and caching
	// Service can function without Redis but with reduced performance
	var redisClient *redis.Client
	if cfg.RedisAddr != "" {
		zapLogger.Info("Connecting to Redis for session management...")
		r, rerr := redisclient.Connect(cfg)
		if rerr != nil {
			zapLogger.Warn("Redis connection failed; continuing without Redis", logger.Err(rerr))
		} else {
			defer r.Close() // Ensure Redis connection is closed on shutdown
			redisClient = r
			zapLogger.Info("Redis connection established successfully")
		}
	} else {
		zapLogger.Info("Redis not configured, skipping Redis connection")
	}

	// Initialize HTTP Router and Middleware
	// Creates Gin router with all routes, middleware, and dependencies injected
	engine := router.NewRouter(cfg, zapLogger, db, redisClient)

	// HTTP Server Configuration
	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,        // Listen on configured port (default: 8080)
		Handler:           engine,                   // Gin router handles all requests
		ReadHeaderTimeout: 20 * time.Second,        // Prevent slowloris attacks
	}

	// Start HTTP Server in Background Goroutine
	// This allows the main goroutine to handle shutdown signals
	go func() {
		zapLogger.Sugar().Infow("HTTP server starting", 
			"addr", srv.Addr,
			"environment", cfg.AppEnv,
			"database_connected", db != nil,
			"redis_connected", redisClient != nil)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("HTTP server failed to start", logger.Err(err))
		}
	}()

	// Graceful Shutdown Handling
	// Wait for interrupt signal (CTRL+C) or termination signal from Docker/Kubernetes
	zapLogger.Info("Server is ready. Press CTRL+C to shutdown gracefully...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until signal received
	
	zapLogger.Info("Shutdown signal received, initiating graceful shutdown...")
	
	// Give server 10 seconds to finish handling existing requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("Forced server shutdown due to timeout", logger.Err(err))
	}
	
	zapLogger.Info("Business Exchange Marketplace server has shut down successfully")

	_ = models.ErrPlaceholder // Prevent unused import error when models only used in migrations
}
