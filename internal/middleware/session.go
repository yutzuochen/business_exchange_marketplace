package middleware

import (
	"net/http"
	"strings"

	"trade_company/internal/auth"
	"trade_company/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SessionAuth struct {
	sessionManager *auth.SessionManager
	config         *config.Config
}

func NewSessionAuth(redisClient *redis.Client, db *gorm.DB, config *config.Config) *SessionAuth {
	sessionManager := auth.NewSessionManager(redisClient, db, config)
	return &SessionAuth{
		sessionManager: sessionManager,
		config:         config,
	}
}

// SessionAuthRequired middleware that requires a valid session
func (sa *SessionAuth) SessionAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := sa.getSessionID(c)
		if sessionID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		session, err := sa.sessionManager.GetSession(sessionID)
		if err != nil {
			// Clear invalid session cookie
			sa.clearSessionCookie(c)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired session",
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", session.UserID)
		c.Set("session_id", session.SessionID)
		c.Set("ip_address", session.IPAddress)
		c.Set("user_agent", session.UserAgent)

		c.Next()
	}
}

// OptionalSessionAuth middleware that sets user info if session exists
func (sa *SessionAuth) OptionalSessionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := sa.getSessionID(c)
		if sessionID == "" {
			c.Next()
			return
		}

		session, err := sa.sessionManager.GetSession(sessionID)
		if err != nil {
			// Clear invalid session cookie
			sa.clearSessionCookie(c)
			c.Next()
			return
		}

		// Set user info in context
		c.Set("user_id", session.UserID)
		c.Set("session_id", session.SessionID)
		c.Set("ip_address", session.IPAddress)
		c.Set("user_agent", session.UserAgent)

		c.Next()
	}
}

// AdminRequired middleware that requires admin role
func (sa *SessionAuth) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// Check if user is admin (this would need to be implemented based on your user model)
		// For now, we'll assume admin check is done elsewhere
		c.Next()
	}
}

// getSessionID extracts session ID from cookie
func (sa *SessionAuth) getSessionID(c *gin.Context) string {
	// Check for session cookie
	cookie, err := c.Cookie("sid")
	if err != nil {
		return ""
	}

	// Validate session ID format (basic check)
	if len(cookie) != 64 { // 32 bytes = 64 hex characters
		return ""
	}

	return cookie
}

// setSessionCookie sets the session cookie
func (sa *SessionAuth) setSessionCookie(c *gin.Context, sessionID string) {
	// Determine cookie domain
	domain := sa.config.SessionCookieDomain
	if domain == "" {
		// Extract domain from request host
		host := c.Request.Host
		if strings.Contains(host, ":") {
			host = strings.Split(host, ":")[0]
		}
		// Set to apex domain if it's a subdomain
		if strings.Count(host, ".") > 1 {
			parts := strings.Split(host, ".")
			if len(parts) >= 2 {
				domain = "." + strings.Join(parts[len(parts)-2:], ".")
			}
		}
	}

	// Set secure flag based on environment
	secure := sa.config.SessionCookieSecure
	if c.Request.TLS == nil && sa.config.AppEnv == "development" {
		secure = false
	}

	// Set cookie
	c.SetCookie(
		"sid",                           // name
		sessionID,                       // value
		sa.config.SessionTTLMinutes*60,  // maxAge (seconds)
		"/",                             // path
		domain,                          // domain
		secure,                          // secure
		sa.config.SessionCookieHttpOnly, // httpOnly
	)
}

// clearSessionCookie clears the session cookie
func (sa *SessionAuth) clearSessionCookie(c *gin.Context) {
	domain := sa.config.SessionCookieDomain
	if domain == "" {
		host := c.Request.Host
		if strings.Contains(host, ":") {
			host = strings.Split(host, ":")[0]
		}
		if strings.Count(host, ".") > 1 {
			parts := strings.Split(host, ".")
			if len(parts) >= 2 {
				domain = "." + strings.Join(parts[len(parts)-2:], ".")
			}
		}
	}

	secure := sa.config.SessionCookieSecure
	if c.Request.TLS == nil && sa.config.AppEnv == "development" {
		secure = false
	}

	c.SetCookie(
		"sid",                           // name
		"",                              // value
		-1,                              // maxAge (expire immediately)
		"/",                             // path
		domain,                          // domain
		secure,                          // secure
		sa.config.SessionCookieHttpOnly, // httpOnly
	)
}

// GetUserID gets the user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	if id, ok := userID.(uint); ok {
		return id, true
	}

	return 0, false
}

// GetSessionID gets the session ID from context
func GetSessionID(c *gin.Context) (string, bool) {
	sessionID, exists := c.Get("session_id")
	if !exists {
		return "", false
	}

	if id, ok := sessionID.(string); ok {
		return id, true
	}

	return "", false
}
