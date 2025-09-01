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

	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBMaxIdleConns int
	DBMaxOpenConns int
	Params         map[string]string

	RedisAddr              string
	RedisPassword          string
	RedisDB                int
	RedisDefaultTTLSeconds int

	JWTSecret        string
	JWTIssuer        string
	JWTExpireMinutes int

	CORSAllowedOrigins string
	CORSAllowedMethods string
	CORSAllowedHeaders string

	// Members service configuration
	SendGridAPIKey    string
	SendGridFromEmail string
	SendGridFromName  string

	// Session management
	SessionSecret         string
	SessionTTLMinutes     int
	SessionCookieDomain   string
	SessionCookieSecure   bool
	SessionCookieHttpOnly bool
	SessionCookieSameSite string

	// Rate limiting
	RateLimitLoginPerMinute        int
	RateLimitSignupPerHour         int
	RateLimitForgotPasswordPerHour int
	RateLimitContactSellerPerHour  int

	// Security
	PasswordMinLength      int
	MaxLoginAttempts       int
	LockoutDurationMinutes int

	// 2FA
	TwoFactorIssuer string

	// File upload limits
	MaxFileSizeMB      int
	MaxTotalSizeMB     int
	MaxFilesPerRequest int
	MaxAvatarSizeMB    int
	GlobalBodyLimitMB  int

	// API 和靜態文件基礎 URL - 根據環境自動設置
	APIBaseURL    string
	StaticBaseURL string
}

func Load() (*Config, error) {
	cfg := &Config{}
	cfg.AppName = getEnv("APP_NAME", "trade_company")
	cfg.AppEnv = getEnv("APP_ENV", "development")

	// Cloud Run 會自動設置 PORT 環境變量，優先使用它
	if port := os.Getenv("PORT"); port != "" {
		cfg.AppPort = port
	} else {
		cfg.AppPort = getEnv("APP_PORT", "8080")
	}

	// Local test
	cfg.DBHost = getEnv("DB_HOST", "127.0.0.1") // this should be noted
	cfg.DBPort = getEnv("DB_PORT", "3306")
	cfg.DBUser = getEnv("DB_USER", "app")
	cfg.DBPassword = getEnv("DB_PASSWORD", "app_password")
	cfg.DBName = getEnv("DB_NAME", "business_exchange")
	cfg.DBMaxIdleConns = getEnvInt("DB_MAX_IDLE_CONNS", 10)
	cfg.DBMaxOpenConns = getEnvInt("DB_MAX_OPEN_CONNS", 50)
	// cfg.Params = map[string]string{
	//     "parseTime":      "true",
	//     "charset":        "utf8mb4",
	//     "loc":            "Local",
	//     "timeout":        "30s",
	//     "readTimeout":    "30s",
	//     "writeTimeout":   "30s",
	//     "multiStatements":"true", // 關鍵
	// }

	// empty by default so Redis is optional in environments without it
	cfg.RedisAddr = getEnv("REDIS_ADDR", "")
	cfg.RedisPassword = getEnv("REDIS_PASSWORD", "")
	cfg.RedisDB = getEnvInt("REDIS_DB", 0)
	cfg.RedisDefaultTTLSeconds = getEnvInt("REDIS_DEFAULT_TTL_SECONDS", 60)

	cfg.JWTSecret = getEnv("JWT_SECRET", "your-local-jwt-secret")
	cfg.JWTIssuer = getEnv("JWT_ISSUER", "trade_company")
	cfg.JWTExpireMinutes = getEnvInt("JWT_EXPIRE_MINUTES", 60)

	cfg.CORSAllowedOrigins = getEnv("CORS_ALLOWED_ORIGINS", "*")
	cfg.CORSAllowedMethods = getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	cfg.CORSAllowedHeaders = getEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization")

	// Members service configuration
	cfg.SendGridAPIKey = getEnv("SENDGRID_API_KEY", "")
	cfg.SendGridFromEmail = getEnv("SENDGRID_FROM_EMAIL", "noreply@business-exchange.com")
	cfg.SendGridFromName = getEnv("SENDGRID_FROM_NAME", "Business Exchange")

	// Session management
	cfg.SessionSecret = getEnv("SESSION_SECRET", "changeme-session-secret")
	cfg.SessionTTLMinutes = getEnvInt("SESSION_TTL_MINUTES", 1440) // 24 hours
	cfg.SessionCookieDomain = getEnv("SESSION_COOKIE_DOMAIN", "")
	cfg.SessionCookieSecure = getEnvBool("SESSION_COOKIE_SECURE", true)
	cfg.SessionCookieHttpOnly = getEnvBool("SESSION_COOKIE_HTTP_ONLY", true)
	cfg.SessionCookieSameSite = getEnv("SESSION_COOKIE_SAME_SITE", "Lax")

	// Rate limiting
	cfg.RateLimitLoginPerMinute = getEnvInt("RATE_LIMIT_LOGIN_PER_MINUTE", 5)
	cfg.RateLimitSignupPerHour = getEnvInt("RATE_LIMIT_SIGNUP_PER_HOUR", 3)
	cfg.RateLimitForgotPasswordPerHour = getEnvInt("RATE_LIMIT_FORGOT_PASSWORD_PER_HOUR", 3)
	cfg.RateLimitContactSellerPerHour = getEnvInt("RATE_LIMIT_CONTACT_SELLER_PER_HOUR", 10)

	// Security
	cfg.PasswordMinLength = getEnvInt("PASSWORD_MIN_LENGTH", 8)
	cfg.MaxLoginAttempts = getEnvInt("MAX_LOGIN_ATTEMPTS", 5)
	cfg.LockoutDurationMinutes = getEnvInt("LOCKOUT_DURATION_MINUTES", 30)

	// 2FA
	cfg.TwoFactorIssuer = getEnv("TWO_FACTOR_ISSUER", "Business Exchange")

	// File upload limits
	cfg.MaxFileSizeMB = getEnvInt("MAX_FILE_SIZE_MB", 5)
	cfg.MaxTotalSizeMB = getEnvInt("MAX_TOTAL_SIZE_MB", 25)
	cfg.MaxFilesPerRequest = getEnvInt("MAX_FILES_PER_REQUEST", 10)
	cfg.MaxAvatarSizeMB = getEnvInt("MAX_AVATAR_SIZE_MB", 1)
	cfg.GlobalBodyLimitMB = getEnvInt("GLOBAL_BODY_LIMIT_MB", 30)

	// API 和靜態文件基礎 URL - 根據環境自動設置
	if cfg.AppEnv == "production" {
		// 生產環境：使用 Cloud Run 的 URL
		cfg.APIBaseURL = getEnv("API_BASE_URL", "https://business-exchange-backend-430730011391.us-central1.run.app")
		cfg.StaticBaseURL = getEnv("STATIC_BASE_URL", "https://business-exchange-backend-430730011391.us-central1.run.app")
	} else {
		// 本地環境：使用 localhost
		cfg.APIBaseURL = getEnv("API_BASE_URL", "http://127.0.0.1:8080")
		cfg.StaticBaseURL = getEnv("STATIC_BASE_URL", "http://127.0.0.1:8080")
	}

	return cfg, nil
}

func (c *Config) MySQLDSN() string {
	// Check if DB_HOST is a Unix socket path (Cloud SQL)
	if len(c.DBHost) > 0 && c.DBHost[0] == '/' {
		// Unix socket connection for Cloud SQL
		// Add additional connection parameters for Cloud SQL
		return fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true&charset=utf8mb4&loc=Local&timeout=30s&readTimeout=30s&writeTimeout=30s", c.DBUser, c.DBPassword, c.DBHost, c.DBName)
	}
	// TCP connection for local development
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local&timeout=30s&readTimeout=30s&writeTimeout=30s", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
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

func getEnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}
