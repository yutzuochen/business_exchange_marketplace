package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"trade_company/internal/auth"
	"trade_company/internal/config"
	"trade_company/internal/middleware"
	"trade_company/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MembersAuthHandler struct {
	DB             *gorm.DB
	RedisClient    *redis.Client
	Config         *config.Config
	SessionManager *auth.SessionManager
	EmailService   *auth.EmailService
}

func NewMembersAuthHandler(db *gorm.DB, redisClient *redis.Client, config *config.Config) *MembersAuthHandler {
	sessionManager := auth.NewSessionManager(redisClient, db, config)
	emailService := auth.NewEmailService(config)

	return &MembersAuthHandler{
		DB:             db,
		RedisClient:    redisClient,
		Config:         config,
		SessionManager: sessionManager,
		EmailService:   emailService,
	}
}

type signupRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`

	// Seller-specific fields
	CompanyName  string `json:"company_name"`
	TaxID        string `json:"tax_id"`
	ContactPhone string `json:"contact_phone"`

	// Anti-bot fields
	Honeypot string `json:"website"`   // Hidden field to catch bots
	FormTime int64  `json:"form_time"` // Time when form was rendered
}

type membersLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type verifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type resetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// Signup handles user registration
func (h *MembersAuthHandler) Signup(c *gin.Context) {
	var req signupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Anti-bot checks
	if req.Honeypot != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check if form was submitted too quickly (less than 800ms)
	if req.FormTime > 0 {
		elapsed := time.Now().UnixMilli() - req.FormTime
		if elapsed < 800 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
	}

	// Check if email already exists
	var existingUser models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Generate email verification token
	verificationToken := h.EmailService.GenerateVerificationToken()

	// Create user
	user := models.User{
		Email:                  req.Email,
		PasswordHash:           string(hashedPassword),
		FirstName:              req.FirstName,
		LastName:               req.LastName,
		Phone:                  req.Phone,
		Role:                   h.getDefaultRole(req.Role),
		IsActive:               false, // Must verify email first
		EmailVerificationToken: verificationToken,
		CompanyName:            req.CompanyName,
		TaxID:                  req.TaxID,
		ContactPhone:           req.ContactPhone,
		EmailNotifications:     true,
		MarketingEmails:        false,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Send verification email
	if err := h.EmailService.SendVerificationEmail(&user, verificationToken); err != nil {
		// Log error but don't fail the request
		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully. Please check your email for verification.",
			"warning": "Verification email could not be sent. Please contact support.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully. Please check your email for verification.",
	})
}

// Login handles user authentication
func (h *MembersAuthHandler) Login(c *gin.Context) {
	var req membersLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not verified. Please check your email."})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.recordFailedLogin(c, req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if account is locked
	if h.isAccountLocked(req.Email) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Account temporarily locked due to too many failed attempts"})
		return
	}

	// Create session
	session, err := h.SessionManager.CreateSession(user.ID, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Set session cookie
	h.setSessionCookie(c, session.SessionID)

	// Update last login time
	h.DB.Model(&user).Update("last_login_at", time.Now())

	// Log successful login
	h.recordSuccessfulLogin(c, user.ID)

	// Check if 2FA is required
	if user.TwoFactorEnabled {
		c.JSON(http.StatusOK, gin.H{
			"message":      "2FA required",
			"requires_2fa": true,
			"session_id":   session.SessionID,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"role":       user.Role,
		},
	})
}

// VerifyEmail handles email verification
func (h *MembersAuthHandler) VerifyEmail(c *gin.Context) {
	var req verifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by verification token
	var user models.User
	if err := h.DB.Where("email_verification_token = ?", req.Token).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification token"})
		return
	}

	// Check if token is expired (24 hours)
	if time.Since(user.CreatedAt) > 24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token expired"})
		return
	}

	// Activate user
	now := time.Now()
	updates := map[string]interface{}{
		"is_active":                true,
		"email_verified_at":        &now,
		"email_verification_token": "",
	}

	if err := h.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully. You can now log in.",
	})
}

// ForgotPassword handles password reset requests
func (h *MembersAuthHandler) ForgotPassword(c *gin.Context) {
	var req forgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Don't reveal if email exists or not
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been sent.",
		})
		return
	}

	// Generate reset token
	resetToken := h.EmailService.GeneratePasswordResetToken()

	// Create or update password reset token
	expiresAt := time.Now().Add(30 * time.Minute)
	resetTokenRecord := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: expiresAt,
	}

	// Delete existing tokens for this user
	h.DB.Where("user_id = ?", user.ID).Delete(&models.PasswordResetToken{})

	// Create new token
	if err := h.DB.Create(&resetTokenRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	// Send reset email
	if err := h.EmailService.SendPasswordResetEmail(&user, resetToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a password reset link has been sent.",
	})
}

// ResetPassword handles password reset
func (h *MembersAuthHandler) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find reset token
	var resetToken models.PasswordResetToken
	if err := h.DB.Where("token = ? AND used = ? AND expires_at > ?",
		req.Token, false, time.Now()).First(&resetToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Update user password
	if err := h.DB.Model(&models.User{}).Where("id = ?", resetToken.UserID).
		Update("password_hash", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Mark token as used
	h.DB.Model(&resetToken).Update("used", true)

	// Revoke all existing sessions for this user
	h.revokeAllUserSessions(resetToken.UserID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully. Please log in with your new password.",
	})
}

// Logout handles user logout
func (h *MembersAuthHandler) Logout(c *gin.Context) {
	sessionID, exists := middleware.GetSessionID(c)
	if !exists {
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
		return
	}

	// Revoke session
	if err := h.SessionManager.RevokeSession(sessionID); err != nil {
		// Log error but don't fail the request
	}

	// Clear session cookie
	h.clearSessionCookie(c)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Helper methods
func (h *MembersAuthHandler) getDefaultRole(requestedRole string) string {
	if requestedRole == "seller" || requestedRole == "admin" {
		return "user" // Default to user role, admin can promote later
	}
	return "user"
}

func (h *MembersAuthHandler) setSessionCookie(c *gin.Context, sessionID string) {
	// This would be implemented in the session middleware
	// For now, we'll use a simple approach
	c.SetCookie("sid", sessionID, h.Config.SessionTTLMinutes*60, "/", "",
		h.Config.SessionCookieSecure, h.Config.SessionCookieHttpOnly)
}

func (h *MembersAuthHandler) clearSessionCookie(c *gin.Context) {
	c.SetCookie("sid", "", -1, "/", "", false, true)
}

func (h *MembersAuthHandler) recordFailedLogin(c *gin.Context, email string) {
	// Record failed login attempt in Redis
	key := fmt.Sprintf("failed_login:%s", email)
	h.RedisClient.Incr(c, key)
	h.RedisClient.Expire(c, key, time.Duration(h.Config.LockoutDurationMinutes)*time.Minute)
}

func (h *MembersAuthHandler) recordSuccessfulLogin(c *gin.Context, userID uint) {
	// Clear failed login attempts
	// This would be implemented based on your audit logging requirements
}

func (h *MembersAuthHandler) isAccountLocked(email string) bool {
	key := fmt.Sprintf("failed_login:%s", email)
	ctx := context.Background()
	count, err := h.RedisClient.Get(ctx, key).Int()
	if err != nil {
		return false
	}
	return count >= h.Config.MaxLoginAttempts
}

func (h *MembersAuthHandler) revokeAllUserSessions(userID uint) {
	// Get all user sessions and revoke them
	sessions, err := h.SessionManager.GetUserSessions(userID)
	if err != nil {
		return
	}

	for _, session := range sessions {
		h.SessionManager.RevokeSession(session.SessionID)
	}
}
