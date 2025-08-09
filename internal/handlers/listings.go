package handlers

import (
	"net/http"
	"strconv"

	"trade_company/internal/middleware"
	"trade_company/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListingsHandler struct {
	DB *gorm.DB
}

type listingRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	PriceUSD    int64  `json:"price_usd" binding:"required"`
	Location    string `json:"location"`
}

func (h *ListingsHandler) Create(c *gin.Context) {
	var req listingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDVal, _ := c.Get(middleware.ContextUserID)
	ownerID := uint(userIDVal.(uint))
	listing := models.Listing{
		Title:       req.Title,
		Description: req.Description,
		PriceUSD:    req.PriceUSD,
		Location:    req.Location,
		OwnerID:     ownerID,
	}
	if err := h.DB.Create(&listing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, listing)
}

func (h *ListingsHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	var listing models.Listing
	if err := h.DB.First(&listing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, listing)
}

func (h *ListingsHandler) List(c *gin.Context) {
	var listings []models.Listing
	if err := h.DB.Order("id desc").Limit(50).Find(&listings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list failed"})
		return
	}
	c.JSON(http.StatusOK, listings)
}
