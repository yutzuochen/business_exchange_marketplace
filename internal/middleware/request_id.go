package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists in header
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// Set request ID in context and header
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		
		c.Next()
	}
}
