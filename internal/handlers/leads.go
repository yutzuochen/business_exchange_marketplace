package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"trade_company/internal/auth"
	"trade_company/internal/config"
	"trade_company/internal/middleware"
	"trade_company/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LeadHandler struct {
	DB           *gorm.DB
	RedisClient  *redis.Client
	Config       *config.Config
	EmailService *auth.EmailService
}

func NewLeadHandler(db *gorm.DB, redisClient *redis.Client, config *config.Config) *LeadHandler {
	emailService := auth.NewEmailService(config)

	return &LeadHandler{
		DB:           db,
		RedisClient:  redisClient,
		Config:       config,
		EmailService: emailService,
	}
}

type contactSellerRequest struct {
	SellerID     uint   `json:"seller_id" binding:"required"`
	ListingID    *uint  `json:"listing_id"`
	Subject      string `json:"subject" binding:"required,max=255"`
	Message      string `json:"message" binding:"required,max=2000"`
	ContactPhone string `json:"contact_phone"`

	// Anti-spam fields
	Honeypot       string `json:"website"`               // Hidden field to catch bots
	FormTime       int64  `json:"form_time"`             // Time when form was rendered
	TurnstileToken string `json:"cf-turnstile-response"` // Cloudflare Turnstile token
}

// ContactSeller handles contact form submissions from buyers to sellers
func (h *LeadHandler) ContactSeller(c *gin.Context) {
	var req contactSellerRequest
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

	// Verify Turnstile token (if enabled)
	if h.Config.AppEnv == "production" && req.TurnstileToken != "" {
		if !h.verifyTurnstileToken(req.TurnstileToken, c.ClientIP()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid security token"})
			return
		}
	}

	// Get sender user ID from session
	senderID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Check if sender is trying to contact themselves
	if senderID == req.SellerID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot contact yourself"})
		return
	}

	// Verify seller exists and is active
	var seller models.User
	if err := h.DB.Where("id = ? AND is_active = ?", req.SellerID, true).First(&seller).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seller not found"})
		return
	}

	// Verify listing exists if provided
	if req.ListingID != nil {
		var listing models.Listing
		if err := h.DB.Where("id = ? AND owner_id = ?", req.ListingID, req.SellerID).First(&listing).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing"})
			return
		}
	}

	// Check rate limiting
	if !h.checkContactRateLimit(senderID, req.SellerID) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many contact requests. Please try again later."})
		return
	}

	// Create lead
	lead := models.Lead{
		SenderID:     senderID,
		ReceiverID:   req.SellerID,
		ListingID:    req.ListingID,
		Subject:      req.Subject,
		Message:      req.Message,
		ContactPhone: req.ContactPhone,
		IsRead:       false,
		IsSpam:       false,
	}

	// Check for spam indicators
	if h.isSpam(lead) {
		lead.IsSpam = true
	}

	if err := h.DB.Create(&lead).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Send email notification to seller
	if err := h.EmailService.SendLeadNotification(&seller, &lead); err != nil {
		// Log error but don't fail the request
	}

	// Record contact for rate limiting
	h.recordContact(senderID, req.SellerID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent successfully",
		"lead_id": lead.ID,
	})
}

// GetUserLeads returns leads for the authenticated user
func (h *LeadHandler) GetUserLeads(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var leads []models.Lead
	if err := h.DB.Where("receiver_id = ?", userID).
		Preload("Sender").
		Preload("Listing").
		Order("created_at DESC").
		Find(&leads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"leads": leads,
	})
}

// MarkLeadAsRead marks a lead as read
func (h *LeadHandler) MarkLeadAsRead(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	leadID := c.Param("id")

	var lead models.Lead
	if err := h.DB.Where("id = ? AND receiver_id = ?", leadID, userID).First(&lead).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"})
		return
	}

	if err := h.DB.Model(&lead).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lead"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lead marked as read",
	})
}

// AdminGetLeads returns all leads for admin users
func (h *LeadHandler) AdminGetLeads(c *gin.Context) {
	// This would check admin role in middleware
	var leads []models.Lead
	if err := h.DB.Preload("Sender").
		Preload("Receiver").
		Preload("Listing").
		Order("created_at DESC").
		Find(&leads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leads"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"leads": leads,
	})
}

// Helper methods
func (h *LeadHandler) checkContactRateLimit(senderID, receiverID uint) bool {
	key := fmt.Sprintf("contact_rate_limit:%d:%d", senderID, receiverID)
	ctx := context.Background()

	count, err := h.RedisClient.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return true // Allow if Redis error
	}

	if count >= h.Config.RateLimitContactSellerPerHour {
		return false
	}

	return true
}

func (h *LeadHandler) recordContact(senderID, receiverID uint) {
	key := fmt.Sprintf("contact_rate_limit:%d:%d", senderID, receiverID)
	ctx := context.Background()

	pipe := h.RedisClient.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Hour)
	pipe.Exec(ctx)
}

func (h *LeadHandler) isSpam(lead models.Lead) bool {
	// Basic spam detection
	spamKeywords := []string{
		"buy now", "click here", "free money", "make money fast",
		"weight loss", "viagra", "casino", "lottery",
	}

	message := lead.Message
	for _, keyword := range spamKeywords {
		if strings.Contains(strings.ToLower(message), keyword) {
			return true
		}
	}

	// Check for excessive links
	linkCount := strings.Count(message, "http")
	if linkCount > 3 {
		return true
	}

	return false
}

func (h *LeadHandler) verifyTurnstileToken(token, ip string) bool {
	// TODO: Implement Cloudflare Turnstile verification
	// For now, return true to allow development
	return true
}
