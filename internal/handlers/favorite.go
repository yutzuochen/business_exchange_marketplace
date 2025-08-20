package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"trade_company/internal/models"
)

type FavoriteHandler struct {
	DB *gorm.DB
}

// List returns the current user's favorites
func (h *FavoriteHandler) List(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var favorites []models.Favorite
	if err := h.DB.Where("user_id = ?", userID).
		Preload("Listing").
		Preload("Listing.Images").
		Preload("Listing.Owner").
		Order("created_at desc").
		Find(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"favorites": favorites,
	})
}

// Add adds a listing to user's favorites
func (h *FavoriteHandler) Add(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input struct {
		ListingID uint `json:"listing_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if listing exists
	var listing models.Listing
	if err := h.DB.First(&listing, input.ListingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	// Check if already favorited
	var existingFavorite models.Favorite
	if err := h.DB.Where("user_id = ? AND listing_id = ?", userID, input.ListingID).
		First(&existingFavorite).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Listing already in favorites"})
		return
	}

	// Create favorite
	favorite := models.Favorite{
		UserID:    userID.(uint),
		ListingID: input.ListingID,
	}

	if err := h.DB.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to favorites"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Added to favorites successfully",
		"favorite": favorite,
	})
}

// Remove removes a listing from user's favorites
func (h *FavoriteHandler) Remove(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	favoriteIDStr := c.Param("id")
	favoriteID, err := strconv.ParseUint(favoriteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid favorite ID"})
		return
	}

	var favorite models.Favorite
	if err := h.DB.Where("id = ? AND user_id = ?", favoriteID, userID).First(&favorite).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite not found"})
		return
	}

	if err := h.DB.Delete(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from favorites successfully"})
}
