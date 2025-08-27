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

func main() {
	fmt.Println("========= lets start =================")
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	zapLogger := logger.New(cfg.AppEnv)
	defer zapLogger.Sync()

	// 嘗試連接數據庫，但不因失敗而退出
	var db *gorm.DB
	dbRetryCount := 0
	maxDbRetries := 5

	for dbRetryCount < maxDbRetries {
		db, err = database.Connect(cfg, zapLogger)
		if err == nil {
			zapLogger.Info("數據庫連接成功")
			break
		}

		dbRetryCount++
		zapLogger.Sugar().Warnw("數據庫連接失敗，重試中", "error", err, "attempt", dbRetryCount)

		if dbRetryCount < maxDbRetries {
			time.Sleep(time.Duration(dbRetryCount) * time.Second)
		}
	}

	if db == nil {
		zapLogger.Error("無法連接到數據庫，但繼續啟動服務器")
	} else {
		zapLogger.Info("數據庫連接成功，開始運行遷移...")

		// 運行數據庫遷移
		if err := database.RunMigrations(db); err != nil {
			zapLogger.Error("database migrations failed :( )", logger.Err(err))
		} else {
			zapLogger.Info("數據庫遷移完成")
		}

		// 只有在數據庫連接成功時才嘗試種子數據
		zapLogger.Info("開始填充種子數據...")
		if err := database.SeedData(db, cfg); err != nil {
			zapLogger.Error("database seeding failed", logger.Err(err))
		} else {
			zapLogger.Info("種子數據填充完成")
		}
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
