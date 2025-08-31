// Package handlers provides HTTP request handlers for the Business Exchange Marketplace API.
// This file contains authentication-related handlers for user registration and login.
package handlers

import (
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
//   {
//     "email": "user@example.com",    // Valid email address (unique)
//     "password": "securepass123"     // Minimum 8 characters
//   }
//
// Response (201 Created):
//   {
//     "message": "User created successfully",
//     "user_id": 123
//   }
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
	h.Log.Info("Registration attempt started", zap.String("ip", c.ClientIP()))
	
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Warn("Registration request validation failed", 
			zap.String("ip", c.ClientIP()), 
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	h.Log.Info("Registration attempt for user", 
		zap.String("email", req.Email), 
		zap.String("ip", c.ClientIP()))
	
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Log.Error("Registration failed - password hashing error", 
			zap.String("email", req.Email), 
			zap.String("ip", c.ClientIP()), 
			logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
		return
	}
	
	user := models.User{Email: req.Email, PasswordHash: string(hash)}
	if err := h.DB.Create(&user).Error; err != nil {
		h.Log.Warn("Registration failed - user creation error", 
			zap.String("email", req.Email), 
			zap.String("ip", c.ClientIP()), 
			logger.Err(err))
		c.JSON(http.StatusConflict, gin.H{"error": "email exists or invalid"})
		return
	}
	
	token, err := auth.GenerateToken(h.Cfg, user.ID, user.Email)
	if err != nil {
		h.Log.Error("Registration failed - token generation error", 
			zap.String("email", req.Email), 
			zap.String("ip", c.ClientIP()), 
			zap.Uint("user_id", user.ID), 
			logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}
	
	h.Log.Info("Registration successful", 
		zap.String("email", req.Email), 
		zap.String("ip", c.ClientIP()), 
		zap.Uint("user_id", user.ID))
	
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	h.Log.Info("Login attempt started", zap.String("ip", c.ClientIP()))
	
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Warn("Login request validation failed", 
			zap.String("ip", c.ClientIP()), 
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	h.Log.Info("Login attempt for user", 
		zap.String("email", req.Email), 
		zap.String("ip", c.ClientIP()))
	
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		h.Log.Warn("Login failed - user not found", 
			zap.String("email", req.Email), 
			zap.String("ip", c.ClientIP()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.Log.Warn("Login failed - invalid password", 
			zap.String("email", req.Email), 
			zap.String("ip", c.ClientIP()), 
			zap.Uint("user_id", user.ID))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	
	token, err := auth.GenerateToken(h.Cfg, user.ID, user.Email)
	if err != nil {
		h.Log.Error("Login failed - token generation error", 
			zap.String("email", req.Email), 
			zap.String("ip", c.ClientIP()), 
			zap.Uint("user_id", user.ID), 
			logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}
	
	h.Log.Info("Login successful", 
		zap.String("email", req.Email), 
		zap.String("ip", c.ClientIP()), 
		zap.Uint("user_id", user.ID))
	
	c.JSON(http.StatusOK, gin.H{"token": token})
}
