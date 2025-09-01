package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
	Issuer string
}

// JWT middleware for authentication
func JWT(config JWTConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// First, try to get token from cookie (preferred method)
		if cookie, err := c.Cookie("authToken"); err == nil && cookie != "" {
			tokenString = cookie
		} else {
			// Fallback to Authorization header for backwards compatibility
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Authentication required: no token found in cookie or Authorization header",
				})
				c.Abort()
				return
			}

			// Check Bearer token format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid authorization header format",
				})
				c.Abort()
				return
			}

			tokenString = parts[1]
		}

		// Parse and validate JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.Secret), nil
		})

		if err != nil {
			logger.Warn("JWT validation failed",
				zap.String("error", err.Error()),
				zap.String("request_id", c.GetString("request_id")),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Check if token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Validate issuer
			if issuer, exists := claims["iss"]; !exists || issuer != config.Issuer {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid token issuer",
				})
				c.Abort()
				return
			}

			// Set user info in context
			if userID, exists := claims["uid"]; exists {
				// Convert userID to uint (JWT numbers are typically float64)
				if userIDFloat, ok := userID.(float64); ok {
					c.Set("user_id", uint(userIDFloat))
				} else {
					logger.Error("Invalid user ID type in JWT claims",
						zap.Any("user_id", userID),
						zap.String("request_id", c.GetString("request_id")),
					)
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "Invalid token format",
					})
					c.Abort()
					return
				}
			} else if userID, exists := claims["sub"]; exists {
				// Fallback to sub claim for backwards compatibility
				if userIDFloat, ok := userID.(float64); ok {
					c.Set("user_id", uint(userIDFloat))
				} else {
					logger.Error("Invalid user ID type in JWT claims",
						zap.Any("user_id", userID),
						zap.String("request_id", c.GetString("request_id")),
					)
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "Invalid token format",
					})
					c.Abort()
					return
				}
			}
			if email, exists := claims["email"]; exists {
				c.Set("user_email", email)
			}
			if role, exists := claims["role"]; exists {
				c.Set("user_role", role)
			}
		}

		c.Next()
	}
}

// OptionalJWT middleware that doesn't require JWT but sets user info if present
func OptionalJWT(config JWTConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Try to parse JWT token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.Secret), nil
		})

		if err != nil || !token.Valid {
			c.Next()
			return
		}

		// Set user info if token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if issuer, exists := claims["iss"]; exists && issuer == config.Issuer {
				if userID, exists := claims["sub"]; exists {
					c.Set("user_id", userID)
				}
				if email, exists := claims["email"]; exists {
					c.Set("user_email", email)
				}
				if role, exists := claims["role"]; exists {
					c.Set("user_role", role)
				}
			}
		}

		c.Next()
	}
}
