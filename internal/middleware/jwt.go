package middleware

import (
	"fmt"
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
		requestID := c.GetString("request_id")
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		logger.Info("JWT middleware: Starting authentication check",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("user_agent", userAgent))

		var tokenString string

		// Debug: Log all cookies received
		cookieHeader := c.GetHeader("Cookie")
		logger.Info("JWT middleware: All cookies received",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("cookie_header", cookieHeader))

		// First, try to get token from cookie (preferred method)
		if cookie, err := c.Cookie("authToken"); err == nil && cookie != "" {
			tokenString = cookie
			logger.Info("JWT middleware: Token found in cookie",
				zap.String("request_id", requestID),
				zap.String("ip", clientIP),
				zap.String("token_length", fmt.Sprintf("%d", len(tokenString))))
		} else {
			logger.Info("JWT middleware: No authToken cookie found - falling back to Authorization header",
				zap.String("request_id", requestID),
				zap.String("ip", clientIP),
				zap.String("cookie_error", err.Error()))

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
		logger.Info("JWT middleware: Starting token validation",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("token_length", fmt.Sprintf("%d", len(tokenString))))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Error("JWT middleware: Invalid signing method",
					zap.String("request_id", requestID),
					zap.String("ip", clientIP),
					zap.String("method", fmt.Sprintf("%T", token.Method)))
				return nil, jwt.ErrSignatureInvalid
			}
			logger.Info("JWT middleware: Token signing method validated",
				zap.String("request_id", requestID),
				zap.String("ip", clientIP))
			return []byte(config.Secret), nil
		})

		if err != nil {
			logger.Warn("JWT middleware: Token validation failed",
				zap.String("error", err.Error()),
				zap.String("request_id", requestID),
				zap.String("ip", clientIP),
				zap.String("user_agent", userAgent),
				zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Check if token is valid
		if !token.Valid {
			logger.Warn("JWT middleware: Token marked as invalid",
				zap.String("request_id", requestID),
				zap.String("ip", clientIP),
				zap.String("user_agent", userAgent),
				zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		logger.Info("JWT middleware: Token validation successful",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP))

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			logger.Info("JWT middleware: Extracting token claims",
				zap.String("request_id", requestID),
				zap.String("ip", clientIP),
				zap.Int("claims_count", len(claims)))

			// Validate issuer
			if issuer, exists := claims["iss"]; !exists || issuer != config.Issuer {
				logger.Warn("JWT middleware: Invalid or missing token issuer",
					zap.String("request_id", requestID),
					zap.String("ip", clientIP),
					zap.Any("found_issuer", issuer),
					zap.String("expected_issuer", config.Issuer),
					zap.Bool("issuer_exists", exists))
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid token issuer",
				})
				c.Abort()
				return
			}

			logger.Info("JWT middleware: Token issuer validated",
				zap.String("request_id", requestID),
				zap.String("ip", clientIP),
				zap.String("issuer", config.Issuer))

			// Store the token string in context for proxy handlers
			c.Set("jwt_token", tokenString)

			// Set user info in context
			if userID, exists := claims["uid"]; exists {
				// Convert userID to uint (JWT numbers are typically float64)
				if userIDFloat, ok := userID.(float64); ok {
					userIDUint := uint(userIDFloat)
					c.Set("user_id", userIDUint)
					logger.Info("JWT middleware: User ID extracted from uid claim",
						zap.String("request_id", requestID),
						zap.String("ip", clientIP),
						zap.Uint("user_id", userIDUint))
				} else {
					logger.Error("JWT middleware: Invalid user ID type in JWT uid claim",
						zap.Any("user_id", userID),
						zap.String("request_id", requestID),
						zap.String("ip", clientIP),
						zap.String("user_id_type", fmt.Sprintf("%T", userID)))
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "Invalid token format",
					})
					c.Abort()
					return
				}
			} else if userID, exists := claims["sub"]; exists {
				// Fallback to sub claim for backwards compatibility
				logger.Info("JWT middleware: Falling back to sub claim for user ID",
					zap.String("request_id", requestID),
					zap.String("ip", clientIP))
				if userIDFloat, ok := userID.(float64); ok {
					userIDUint := uint(userIDFloat)
					c.Set("user_id", userIDUint)
					logger.Info("JWT middleware: User ID extracted from sub claim",
						zap.String("request_id", requestID),
						zap.String("ip", clientIP),
						zap.Uint("user_id", userIDUint))
				} else {
					logger.Error("JWT middleware: Invalid user ID type in JWT sub claim",
						zap.Any("user_id", userID),
						zap.String("request_id", requestID),
						zap.String("ip", clientIP),
						zap.String("user_id_type", fmt.Sprintf("%T", userID)))
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "Invalid token format",
					})
					c.Abort()
					return
				}
			} else {
				logger.Warn("JWT middleware: No user ID found in token claims",
					zap.String("request_id", requestID),
					zap.String("ip", clientIP),
					zap.Any("available_claims", claims))
			}

			if email, exists := claims["email"]; exists {
				c.Set("user_email", email)
				logger.Info("JWT middleware: User email extracted from claims",
					zap.String("request_id", requestID),
					zap.String("ip", clientIP),
					zap.String("user_email", fmt.Sprintf("%v", email)))
			}
			if role, exists := claims["role"]; exists {
				c.Set("user_role", role)
				logger.Info("JWT middleware: User role extracted from claims",
					zap.String("request_id", requestID),
					zap.String("ip", clientIP),
					zap.String("user_role", fmt.Sprintf("%v", role)))
			}
		} else {
			logger.Error("JWT middleware: Failed to extract JWT claims",
				zap.String("request_id", requestID),
				zap.String("ip", clientIP),
				zap.String("claims_type", fmt.Sprintf("%T", token.Claims)))
		}

		logger.Info("JWT middleware: Authentication successful - proceeding to next handler",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Any("user_id", c.GetString("user_id")))

		c.Next()
	}
}

// OptionalJWT middleware that doesn't require JWT but sets user info if present
func OptionalJWT(config JWTConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		logger.Info("OptionalJWT middleware: Starting optional authentication check",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("user_agent", userAgent))

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Info("OptionalJWT middleware: No Authorization header found - proceeding without authentication",
				zap.String("request_id", requestID),
				zap.String("ip", clientIP))
			c.Next()
			return
		}

		// Try to parse JWT token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Info("OptionalJWT middleware: Invalid Authorization header format - proceeding without authentication",
				zap.String("request_id", requestID),
				zap.String("ip", clientIP),
				zap.String("auth_header_format", authHeader))
			c.Next()
			return
		}

		tokenString := parts[1]
		logger.Info("OptionalJWT middleware: Found Bearer token - attempting validation",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("token_length", fmt.Sprintf("%d", len(tokenString))))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Warn("OptionalJWT middleware: Invalid signing method in optional token",
					zap.String("request_id", requestID),
					zap.String("ip", clientIP),
					zap.String("method", fmt.Sprintf("%T", token.Method)))
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.Secret), nil
		})

		if err != nil || !token.Valid {
			logger.Info("OptionalJWT middleware: Token validation failed - proceeding without authentication",
				zap.String("request_id", requestID),
				zap.String("ip", clientIP),
				zap.String("error", fmt.Sprintf("%v", err)),
				zap.Bool("token_valid", token != nil && token.Valid))
			c.Next()
			return
		}

		logger.Info("OptionalJWT middleware: Token validation successful",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP))

		// Set user info if token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if issuer, exists := claims["iss"]; exists && issuer == config.Issuer {
				logger.Info("OptionalJWT middleware: Token issuer validated - extracting user info",
					zap.String("request_id", requestID),
					zap.String("ip", clientIP),
					zap.String("issuer", config.Issuer))

				if userID, exists := claims["sub"]; exists {
					c.Set("user_id", userID)
					logger.Info("OptionalJWT middleware: User ID extracted from sub claim",
						zap.String("request_id", requestID),
						zap.String("ip", clientIP),
						zap.Any("user_id", userID))
				}
				if email, exists := claims["email"]; exists {
					c.Set("user_email", email)
					logger.Info("OptionalJWT middleware: User email extracted from claims",
						zap.String("request_id", requestID),
						zap.String("ip", clientIP),
						zap.String("user_email", fmt.Sprintf("%v", email)))
				}
				if role, exists := claims["role"]; exists {
					c.Set("user_role", role)
					logger.Info("OptionalJWT middleware: User role extracted from claims",
						zap.String("request_id", requestID),
						zap.String("ip", clientIP),
						zap.String("user_role", fmt.Sprintf("%v", role)))
				}
			} else {
				logger.Warn("OptionalJWT middleware: Token issuer validation failed - proceeding without authentication",
					zap.String("request_id", requestID),
					zap.String("ip", clientIP),
					zap.Any("found_issuer", issuer),
					zap.String("expected_issuer", config.Issuer),
					zap.Bool("issuer_exists", exists))
			}
		} else {
			logger.Warn("OptionalJWT middleware: Failed to extract JWT claims - proceeding without authentication",
				zap.String("request_id", requestID),
				zap.String("ip", clientIP),
				zap.String("claims_type", fmt.Sprintf("%T", token.Claims)))
		}

		userIDValue, userIDExists := c.Get("user_id")
		logger.Info("OptionalJWT middleware: Optional authentication complete - proceeding to next handler",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("path", c.Request.URL.Path),
			zap.Any("user_id", userIDValue),
			zap.Bool("user_authenticated", userIDExists))

		c.Next()
	}
}
