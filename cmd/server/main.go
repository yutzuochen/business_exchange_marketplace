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
)

func main() {
	fmt.Println("========= lets start =================")
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

	// Seed database with sample data
	if err := database.SeedData(db); err != nil {
		zapLogger.Warn("database seeding failed", logger.Err(err))
	}

	var redisClient *redis.Client
	if cfg.RedisAddr != "" {
		r, rerr := redisclient.Connect(cfg)
		if rerr != nil {
			zapLogger.Warn("redis connect failed; continuing without redis", logger.Err(rerr))
		} else {
			defer r.Close()
			redisClient = r
		}
	}

	engine := router.NewRouter(cfg, zapLogger, db, redisClient)

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           engine,
		ReadHeaderTimeout: 20 * time.Second,
	}
	fmt.Printf("srv: %+v\n", srv)

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
