package middleware

import (
	"fmt"
	"net/http"
	"time"

	"trade_company/internal/config"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redisClient *redis.Client
	config      *config.Config
}

func NewRateLimiter(redisClient *redis.Client, config *config.Config) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		config:      config,
	}
}

// RateLimitLogin limits login attempts per IP address
func (rl *RateLimiter) RateLimitLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:login:%s", ip)

		if !rl.checkRateLimit(key, rl.config.RateLimitLoginPerMinute, time.Minute) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many login attempts. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitSignup limits signup attempts per IP address
func (rl *RateLimiter) RateLimitSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:signup:%s", ip)

		if !rl.checkRateLimit(key, rl.config.RateLimitSignupPerHour, time.Hour) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many signup attempts. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitForgotPassword limits password reset requests per email
func (rl *RateLimiter) RateLimitForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			c.Abort()
			return
		}

		key := fmt.Sprintf("rate_limit:forgot_password:%s", req.Email)

		if !rl.checkRateLimit(key, rl.config.RateLimitForgotPasswordPerHour, time.Hour) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many password reset requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitContactSeller limits contact seller form submissions per IP
func (rl *RateLimiter) RateLimitContactSeller() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:contact_seller:%s", ip)

		if !rl.checkRateLimit(key, rl.config.RateLimitContactSellerPerHour, time.Hour) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many contact requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkRateLimit checks if the rate limit has been exceeded
func (rl *RateLimiter) checkRateLimit(key string, limit int, window time.Duration) bool {
	ctx := context.Background()

	// Get current count
	count, err := rl.redisClient.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		// Redis error, allow request
		return true
	}

	if count >= limit {
		return false
	}

	// Increment counter
	pipe := rl.redisClient.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err = pipe.Exec(ctx)

	if err != nil {
		// Redis error, allow request
		return true
	}

	return true
}
