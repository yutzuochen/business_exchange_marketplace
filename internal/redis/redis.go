package redis

import (
	"business-marketplace/internal/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Initialize creates and configures a Redis client
func Initialize(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return rdb, nil
}
