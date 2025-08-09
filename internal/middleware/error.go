package middleware

import "github.com/gin-gonic/gin"

func JSONError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}
