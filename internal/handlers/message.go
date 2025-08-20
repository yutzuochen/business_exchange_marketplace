package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"trade_company/internal/models"
)

type MessageHandler struct {
	DB *gorm.DB
}

// List returns the current user's messages
func (h *MessageHandler) List(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var messages []models.Message
	if err := h.DB.Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Preload("Sender").
		Preload("Receiver").
		Preload("Listing").
		Order("created_at desc").
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

// Get returns a specific message
func (h *MessageHandler) Get(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message models.Message
	if err := h.DB.Where("id = ? AND (sender_id = ? OR receiver_id = ?)", messageID, userID, userID).
		Preload("Sender").
		Preload("Receiver").
		Preload("Listing").
		First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

// Create creates a new message
func (h *MessageHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input struct {
		ReceiverID uint   `json:"receiver_id" binding:"required"`
		ListingID  *uint  `json:"listing_id"`
		Subject    string `json:"subject"`
		Content    string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if len(input.Content) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message content is required"})
		return
	}

	// Check if receiver exists
	var receiver models.User
	if err := h.DB.First(&receiver, input.ReceiverID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receiver not found"})
		return
	}

	// Check if listing exists (if provided)
	if input.ListingID != nil {
		var listing models.Listing
		if err := h.DB.First(&listing, *input.ListingID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
			return
		}
	}

	// Create message
	message := models.Message{
		SenderID:   userID.(uint),
		ReceiverID: input.ReceiverID,
		ListingID:  input.ListingID,
		Subject:    input.Subject,
		Content:    input.Content,
		IsRead:     false,
	}

	if err := h.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Message sent successfully",
		"data":    message,
	})
}

// MarkAsRead marks a message as read
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	messageIDStr := c.Param("id")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message models.Message
	if err := h.DB.Where("id = ? AND receiver_id = ?", messageID, userID).First(&message).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Mark as read
	message.IsRead = true
	if err := h.DB.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark message as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message marked as read",
		"data":    message,
	})
}
