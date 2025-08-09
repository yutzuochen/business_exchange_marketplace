package middleware

import (
	"net/http"
	"strings"

	"trade_company/internal/auth"
	"trade_company/internal/config"

	"github.com/gin-gonic/gin"
)

const ContextUserID = "userID"

func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimPrefix(authz, "Bearer ")
		claims, err := auth.ParseToken(cfg, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set(ContextUserID, claims.UserID)
		c.Next()
	}
}
