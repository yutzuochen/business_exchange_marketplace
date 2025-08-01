package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server Configuration
	ServerHost string
	ServerPort string
	GinMode    string

	// Database Configuration
	DBHost               string
	DBPort               string
	DBName               string
	DBUser               string
	DBPassword           string
	DBMaxConnections     int
	DBMaxIdleConnections int

	// Redis Configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// JWT Configuration
	JWTSecret      string
	JWTExpiryHours time.Duration

	// Email Configuration
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string

	// File Upload Configuration
	UploadPath       string
	MaxFileSize      int64
	AllowedFileTypes []string

	// Cache Configuration
	CacheDefaultTTL time.Duration
	CacheSearchTTL  time.Duration
	CacheSessionTTL time.Duration

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// Pagination
	DefaultPageSize int
	MaxPageSize     int

	// Application Settings
	AppName      string
	AppURL       string
	ContactEmail string

	// Security
	CSRFSecret    string
	SessionSecret string
}

func Load() *Config {
	return &Config{
		// Server Configuration
		ServerHost: getEnv("SERVER_HOST", "localhost"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),

		// Database Configuration
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "3306"),
		DBName:               getEnv("DB_NAME", "business_marketplace"),
		DBUser:               getEnv("DB_USER", "root"),
		DBPassword:           getEnv("DB_PASSWORD", ""),
		DBMaxConnections:     getEnvAsInt("DB_MAX_CONNECTIONS", 20),
		DBMaxIdleConnections: getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 10),

		// Redis Configuration
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		// JWT Configuration
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-this"),
		JWTExpiryHours: time.Duration(getEnvAsInt("JWT_EXPIRY_HOURS", 24)) * time.Hour,

		// Email Configuration
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", "noreply@businessmarketplace.com"),

		// File Upload Configuration
		UploadPath:       getEnv("UPLOAD_PATH", "./uploads"),
		MaxFileSize:      int64(getEnvAsInt("MAX_FILE_SIZE", 10485760)), // 10MB
		AllowedFileTypes: []string{"jpg", "jpeg", "png", "gif", "pdf"},

		// Cache Configuration
		CacheDefaultTTL: time.Duration(getEnvAsInt("CACHE_DEFAULT_TTL", 300)) * time.Second,
		CacheSearchTTL:  time.Duration(getEnvAsInt("CACHE_SEARCH_TTL", 300)) * time.Second,
		CacheSessionTTL: time.Duration(getEnvAsInt("CACHE_SESSION_TTL", 1800)) * time.Second,

		// Rate Limiting
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   time.Duration(getEnvAsInt("RATE_LIMIT_WINDOW", 60)) * time.Second,

		// Pagination
		DefaultPageSize: getEnvAsInt("DEFAULT_PAGE_SIZE", 20),
		MaxPageSize:     getEnvAsInt("MAX_PAGE_SIZE", 100),

		// Application Settings
		AppName:      getEnv("APP_NAME", "Business Marketplace"),
		AppURL:       getEnv("APP_URL", "http://localhost:8080"),
		ContactEmail: getEnv("CONTACT_EMAIL", "support@businessmarketplace.com"),

		// Security
		CSRFSecret:    getEnv("CSRF_SECRET", "your-csrf-secret"),
		SessionSecret: getEnv("SESSION_SECRET", "your-session-secret"),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}
