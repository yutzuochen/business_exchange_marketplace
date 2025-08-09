package main

import (
	"context"
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
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	zapLogger := logger.New(cfg.AppEnv)
	defer zapLogger.Sync()

	db, err := database.Connect(cfg, zapLogger)
	if err != nil {
		zapLogger.Fatal("db connect", logger.Err(err))
	}

	if err := database.AutoMigrate(db); err != nil {
		zapLogger.Fatal("db automigrate", logger.Err(err))
	}

	redis, err := redisclient.Connect(cfg)
	if err != nil {
		zapLogger.Fatal("redis connect", logger.Err(err))
	}
	defer redis.Close()

	engine := router.NewRouter(cfg, zapLogger, db, redis)

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           engine,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		zapLogger.Sugar().Infow("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("server error", logger.Err(err))
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("server shutdown", logger.Err(err))
	}
	zapLogger.Info("server exited")

	_ = models.ErrPlaceholder // avoid unused import if models only used in migration
} 