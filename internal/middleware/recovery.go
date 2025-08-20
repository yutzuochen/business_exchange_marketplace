package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery middleware for handling panics
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID, exists := c.Get("request_id")
		if !exists {
			requestID = "unknown"
		}
		
		logger.Error("Panic recovered",
			zap.Any("panic", recovered),
			zap.String("request_id", requestID.(string)),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("stack", string(debug.Stack())),
		)
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"request_id": requestID,
		})
	})
}
