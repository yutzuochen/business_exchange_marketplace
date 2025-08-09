package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBMaxIdleConns  int
	DBMaxOpenConns  int

	RedisAddr              string
	RedisPassword          string
	RedisDB                int
	RedisDefaultTTLSeconds int

	JWTSecret       string
	JWTIssuer       string
	JWTExpireMinutes int

	CORSAllowedOrigins string
	CORSAllowedMethods string
	CORSAllowedHeaders string
}

func Load() (*Config, error) {
	cfg := &Config{}
	cfg.AppName = getEnv("APP_NAME", "trade_company")
	cfg.AppEnv = getEnv("APP_ENV", "development")
	cfg.AppPort = getEnv("APP_PORT", "8080")

	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort = getEnv("DB_PORT", "3306")
	cfg.DBUser = getEnv("DB_USER", "app")
	cfg.DBPassword = getEnv("DB_PASSWORD", "app_password")
	cfg.DBName = getEnv("DB_NAME", "trade_company")
	cfg.DBMaxIdleConns = getEnvInt("DB_MAX_IDLE_CONNS", 10)
	cfg.DBMaxOpenConns = getEnvInt("DB_MAX_OPEN_CONNS", 50)

	cfg.RedisAddr = getEnv("REDIS_ADDR", "localhost:6379")
	cfg.RedisPassword = getEnv("REDIS_PASSWORD", "")
	cfg.RedisDB = getEnvInt("REDIS_DB", 0)
	cfg.RedisDefaultTTLSeconds = getEnvInt("REDIS_DEFAULT_TTL_SECONDS", 60)

	cfg.JWTSecret = getEnv("JWT_SECRET", "changeme-super-secret")
	cfg.JWTIssuer = getEnv("JWT_ISSUER", "trade_company")
	cfg.JWTExpireMinutes = getEnvInt("JWT_EXPIRE_MINUTES", 60)

	cfg.CORSAllowedOrigins = getEnv("CORS_ALLOWED_ORIGINS", "*")
	cfg.CORSAllowedMethods = getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	cfg.CORSAllowedHeaders = getEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization")
	return cfg, nil
}

func (c *Config) MySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
} 