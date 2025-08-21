package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS middleware configuration
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow localhost and 127.0.0.1 with any port for development
		allowed := false
		if origin != "" {
			// Allow localhost with any port
			if strings.HasPrefix(origin, "http://localhost:") {
				allowed = true
			}
			// Allow 127.0.0.1 with any port
			if strings.HasPrefix(origin, "http://127.0.0.1:") {
				allowed = true
			}
			// Allow specific network IPs for development
			if strings.HasPrefix(origin, "http://192.168.") {
				allowed = true
			}
			// Allow specific network IPs for development
			if strings.HasPrefix(origin, "http://172.") {
				allowed = true
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Request-ID, Origin")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
