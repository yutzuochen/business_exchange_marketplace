// Package handlers provides HTTP request handlers for the Business Exchange Marketplace API.
// This file contains authentication-related handlers for user registration and login.
package handlers

import (
	"fmt"
	"net/http"

	"trade_company/internal/auth"
	"trade_company/internal/config"
	"trade_company/internal/logger"
	"trade_company/internal/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler handles authentication-related HTTP requests.
// Provides endpoints for user registration, login, and token management.
//
// Dependencies:
//   - DB: GORM database connection for user persistence
//   - Cfg: Application configuration for JWT settings
//   - Log: Structured logger for security event logging
type AuthHandler struct {
	DB  *gorm.DB       // Database connection for user operations
	Cfg *config.Config // Configuration for JWT token generation
	Log *zap.Logger    // Logger for authentication events
}

// registerRequest defines the JSON payload structure for user registration.
//
// Validation rules:
//   - Email: Must be a valid email format (RFC 5322)
//   - Password: Minimum 8 characters for security
type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`    // User's email address (unique identifier)
	Password string `json:"password" binding:"required,min=8"` // Plain text password (hashed before storage)
}

// loginRequest defines the JSON payload structure for user authentication.
//
// Fields:
//   - Email: User's registered email address
//   - Password: Plain text password for verification
type loginRequest struct {
	Email    string `json:"email" binding:"required,email"` // User's email address
	Password string `json:"password" binding:"required"`    // Plain text password for verification
}

// Register handles new user registration requests.
//
// This endpoint creates a new user account with email and password authentication.
// The password is securely hashed using bcrypt before database storage.
//
// HTTP Method: POST
// Endpoint: /api/v1/auth/register
// Content-Type: application/json
//
// Request Body:
//
//	{
//	  "email": "user@example.com",    // Valid email address (unique)
//	  "password": "securepass123"     // Minimum 8 characters
//	}
//
// Response (201 Created):
//
//	{
//	  "message": "User created successfully",
//	  "user_id": 123
//	}
//
// Error Responses:
//   - 400 Bad Request: Invalid email format or password too short
//   - 409 Conflict: Email already exists
//   - 500 Internal Server Error: Database or hashing failure
//
// Security features:
//   - bcrypt password hashing with default cost (10)
//   - Email uniqueness validation
//   - Input validation and sanitization
//   - Comprehensive security event logging
func (h *AuthHandler) Register(c *gin.Context) {
	requestID := c.GetString("request_id")
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	h.Log.Info("AuthHandler: Registration attempt started",
		zap.String("request_id", requestID),
		zap.String("ip", clientIP),
		zap.String("user_agent", userAgent),
		zap.String("endpoint", "/api/v1/auth/register"))

	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Warn("AuthHandler: Registration request validation failed",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Error(err),
			zap.String("validation_error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.Log.Info("AuthHandler: Registration request validated successfully",
		zap.String("request_id", requestID),
		zap.String("email", req.Email),
		zap.String("ip", clientIP),
		zap.String("user_agent", userAgent),
		zap.Int("password_length", len(req.Password)))

	h.Log.Info("AuthHandler: Starting password hashing",
		zap.String("request_id", requestID),
		zap.String("email", req.Email),
		zap.String("ip", clientIP))

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Log.Error("AuthHandler: Registration failed - password hashing error",
			zap.String("request_id", requestID),
			zap.String("email", req.Email),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
		return
	}

	h.Log.Info("AuthHandler: Password hashing successful - creating user",
		zap.String("request_id", requestID),
		zap.String("email", req.Email),
		zap.String("ip", clientIP))

	user := models.User{Email: req.Email, PasswordHash: string(hash)}
	if err := h.DB.Create(&user).Error; err != nil {
		h.Log.Warn("AuthHandler: Registration failed - user creation error",
			zap.String("request_id", requestID),
			zap.String("email", req.Email),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			logger.Err(err),
			zap.String("database_error", err.Error()))
		c.JSON(http.StatusConflict, gin.H{"error": "email exists or invalid"})
		return
	}

	h.Log.Info("AuthHandler: User created successfully - generating JWT token",
		zap.String("request_id", requestID),
		zap.String("email", req.Email),
		zap.String("ip", clientIP),
		zap.Uint("user_id", user.ID))

	token, err := auth.GenerateToken(h.Cfg, user.ID, user.Email)
	if err != nil {
		h.Log.Error("AuthHandler: Registration failed - token generation error",
			zap.String("request_id", requestID),
			zap.String("email", req.Email),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Uint("user_id", user.ID),
			logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	h.Log.Info("AuthHandler: Registration successful - returning token",
		zap.String("request_id", requestID),
		zap.String("email", req.Email),
		zap.String("ip", clientIP),
		zap.String("user_agent", userAgent),
		zap.Uint("user_id", user.ID),
		zap.Int("token_length", len(token)))

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	requestID := c.GetString("request_id")
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	h.Log.Info("AuthHandler: Login attempt started",
		zap.String("request_id", requestID),
		zap.String("ip", clientIP),
		zap.String("user_agent", userAgent),
		zap.String("endpoint", "/api/v1/auth/login"))

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Warn("AuthHandler: Login request validation failed",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Error(err),
			zap.String("validation_error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.Log.Info("AuthHandler: Login request validated successfully, Searching for user in database",
		zap.String("request_id", requestID),
		zap.String("email", req.Email),
		zap.String("ip", clientIP),
		zap.String("user_agent", userAgent),
		zap.Int("password_length", len(req.Password)))

	// h.Log.Info("AuthHandler: Searching for user in database",
	// 	zap.String("request_id", requestID),
	// 	zap.String("email", req.Email),
	// 	zap.String("ip", clientIP))

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		h.Log.Warn("AuthHandler: Login failed - user not found",
			zap.String("request_id", requestID),
			zap.String("email", req.Email),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			logger.Err(err),
			zap.String("database_error", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	h.Log.Info("AuthHandler: User found - verifying password",
		zap.String("request_id", requestID),
		zap.String("email", req.Email),
		zap.String("ip", clientIP),
		zap.Uint("user_id", user.ID),
		zap.Bool("user_is_active", user.IsActive))

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.Log.Warn("AuthHandler: Login failed - invalid password",
			zap.String("request_id", requestID),
			zap.String("email", req.Email),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Uint("user_id", user.ID),
			logger.Err(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	h.Log.Info("AuthHandler: Password verification successful - generating JWT token",
		zap.String("request_id", requestID),
		zap.String("email", req.Email),
		zap.String("ip", clientIP),
		zap.Uint("user_id", user.ID))

	token, err := auth.GenerateToken(h.Cfg, user.ID, user.Email)
	if err != nil {
		h.Log.Error("AuthHandler: Login failed - token generation error",
			zap.String("request_id", requestID),
			zap.String("email", req.Email),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Uint("user_id", user.ID),
			logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	h.Log.Info("AuthHandler: JWT token generated successfully - setting cookie",
		zap.String("request_id", requestID),
		zap.String("email", req.Email),
		zap.String("ip", clientIP),
		zap.Uint("user_id", user.ID),
		zap.Int("token_length", len(token)),
		zap.Int("expire_minutes", h.Cfg.JWTExpireMinutes))

	// Set JWT token as HTTP-only cookie for security
	if h.Cfg.AppEnv == "development" {
		// For development, use standard SetCookie with localhost domain
		// SameSite=Lax works better for localhost development than SameSite=None
		c.SetCookie(
			"authToken",                    // Cookie name
			token,                          // JWT token value
			int(h.Cfg.JWTExpireMinutes*60), // Max age in seconds
			"/",                            // Path (all routes)
			"localhost",                    // Domain (localhost for cross-port support)
			false,                          // Secure flag (false for HTTP development)
			true,                           // HttpOnly flag (prevents JavaScript access)
		)
		h.Log.Info("AuthHandler: Development cookie set with localhost domain",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("domain", "localhost"),
			zap.String("app_env", h.Cfg.AppEnv),
			zap.Bool("secure", false),
			zap.Bool("http_only", true))
	} else {
		// Production cookie with Secure flag
		c.SetCookie(
			"authToken",                    // Cookie name
			token,                          // JWT token value
			int(h.Cfg.JWTExpireMinutes*60), // Max age in seconds
			"/",                            // Path (all routes)
			"",                             // Domain (empty for production)
			true,                           // Secure flag (requires HTTPS)
			true,                           // HttpOnly flag (prevents JavaScript access)
		)
	}

	h.Log.Info("AuthHandler: Login successful - cookie set, returning response",
		zap.String("request_id", requestID),
		zap.String("email", req.Email),
		zap.String("ip", clientIP),
		zap.String("user_agent", userAgent),
		zap.Uint("user_id", user.ID),
		zap.String("app_env", h.Cfg.AppEnv),
		zap.Int("cookie_max_age", int(h.Cfg.JWTExpireMinutes*60)))

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user_id": user.ID,
	})
}

// Logout handles user logout requests by clearing the authentication cookie.
//
// This endpoint securely logs out the user by setting the authToken cookie
// to expire immediately, effectively clearing it from the browser.
//
// HTTP Method: POST
// Endpoint: /api/v1/auth/logout
//
// Response (200 OK):
//
//	{
//	  "message": "Logout successful"
//	}
//
// Security features:
//   - Immediately expires the authentication cookie
//   - Clears session on the client side
//   - Prevents session hijacking after logout
func (h *AuthHandler) Logout(c *gin.Context) {
	requestID := c.GetString("request_id")
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Try to get user info before clearing session
	userID, userIDExists := c.Get("user_id")
	userEmail, emailExists := c.Get("user_email")

	h.Log.Info("AuthHandler: Logout request started",
		zap.String("request_id", requestID),
		zap.String("ip", clientIP),
		zap.String("user_agent", userAgent),
		zap.String("endpoint", "/api/v1/auth/logout"),
		zap.Any("user_id", userID),
		zap.Bool("user_authenticated", userIDExists),
		zap.Any("user_email", userEmail),
		zap.Bool("email_exists", emailExists))

	// Clear the authentication cookie by setting it to expire immediately
	if h.Cfg.AppEnv == "development" {
		// Development logout cookie with localhost domain
		c.SetCookie(
			"authToken", // Cookie name
			"",          // Empty value
			-1,          // Max age -1 (expires immediately)
			"/",         // Path (all routes)
			"localhost", // Domain (localhost for development)
			false,       // Secure flag (false for HTTP development)
			true,        // HttpOnly flag
		)
		h.Log.Info("AuthHandler: Development logout cookie cleared with localhost domain",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("domain", "localhost"),
			zap.String("app_env", h.Cfg.AppEnv))
	} else {
		// Production logout cookie
		c.SetCookie(
			"authToken", // Cookie name
			"",          // Empty value
			-1,          // Max age -1 (expires immediately)
			"/",         // Path (all routes)
			"",          // Domain (empty for production)
			true,        // Secure flag (requires HTTPS)
			true,        // HttpOnly flag
		)
	}

	h.Log.Info("AuthHandler: Logout successful - cookie cleared, returning response",
		zap.String("request_id", requestID),
		zap.String("ip", clientIP),
		zap.String("user_agent", userAgent),
		zap.Any("logged_out_user_id", userID),
		zap.Any("logged_out_user_email", userEmail),
		zap.String("app_env", h.Cfg.AppEnv))

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

// Me handles requests to get the current user's information.
//
// This endpoint returns the authenticated user's profile information
// based on the JWT token in the request context.
//
// HTTP Method: GET
// Endpoint: /api/v1/auth/me
//
// Response (200 OK):
//
//	{
//	  "data": {
//	    "id": 123,
//	    "email": "user@example.com",
//	    "username": "username",
//	    "first_name": "First",
//	    "last_name": "Last"
//	  }
//	}
//
// Security features:
//   - Requires valid JWT token
//   - Returns only the authenticated user's data
func (h *AuthHandler) Me(c *gin.Context) {
	requestID := c.GetString("request_id")
	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	h.Log.Info("AuthHandler: Me request started",
		zap.String("request_id", requestID),
		zap.String("ip", clientIP),
		zap.String("user_agent", userAgent),
		zap.String("endpoint", "/api/v1/auth/me"))

	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		h.Log.Warn("AuthHandler: Me request failed - no user ID in context",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.String("auth_error", "no_user_id_in_context"))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	h.Log.Info("AuthHandler: User ID found in context - validating type",
		zap.String("request_id", requestID),
		zap.String("ip", clientIP),
		zap.Any("user_id_raw", userID),
		zap.String("user_id_type", fmt.Sprintf("%T", userID)))

	userIDValue, ok := userID.(uint)
	if !ok {
		h.Log.Error("AuthHandler: Me request failed - invalid user ID type in context",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Any("user_id_value", userID),
			zap.String("expected_type", "uint"),
			zap.String("actual_type", fmt.Sprintf("%T", userID)))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.Log.Info("AuthHandler: User ID validated - fetching user from database",
		zap.String("request_id", requestID),
		zap.String("ip", clientIP),
		zap.Uint("user_id", userIDValue))

	// Get user information from database
	var user models.User
	if err := h.DB.First(&user, userIDValue).Error; err != nil {
		h.Log.Error("AuthHandler: Me request failed - user not found in database",
			zap.String("request_id", requestID),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Uint("user_id", userIDValue),
			logger.Err(err),
			zap.String("database_error", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	h.Log.Info("AuthHandler: User found in database - returning user information",
		zap.String("request_id", requestID),
		zap.String("ip", clientIP),
		zap.String("user_agent", userAgent),
		zap.Uint("user_id", userIDValue),
		zap.String("user_email", user.Email),
		zap.String("user_role", user.Role),
		zap.Bool("user_is_active", user.IsActive),
		zap.Bool("has_username", user.Username != ""),
		zap.Bool("has_first_name", user.FirstName != ""),
		zap.Bool("has_last_name", user.LastName != ""))

	// Return user information (excluding sensitive data like password hash)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"username":   user.Username,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"role":       user.Role,
			"is_active":  user.IsActive,
		},
	})
}
