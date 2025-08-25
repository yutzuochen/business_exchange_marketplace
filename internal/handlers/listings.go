package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"trade_company/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListingsHandler struct {
	DB *gorm.DB
}

func (h *ListingsHandler) checkDB(c *gin.Context) bool {
	if h.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database not available"})
		return false
	}

	// Check if database connection is alive
	sqlDB, err := h.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database connection error"})
		return false
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database ping failed"})
		return false
	}

	return true
}

type listingRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Price       int64  `json:"price" binding:"required"`
	Category    string `json:"category"`
	Condition   string `json:"condition"`
	Location    string `json:"location"`
}

type listingUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Price       *int64  `json:"price"`
	Category    *string `json:"category"`
	Condition   *string `json:"condition"`
	Location    *string `json:"location"`
	Status      *string `json:"status"`
}

func (h *ListingsHandler) Create(c *gin.Context) {
	if !h.checkDB(c) {
		return
	}

	var req listingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ownerID := userID.(uint)
	listing := models.Listing{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Condition:   req.Condition,
		Location:    req.Location,
		OwnerID:     ownerID,
		Status:      "活躍",
	}

	if err := h.DB.Create(&listing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create listing"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Listing created successfully",
		"listing": listing,
	})
}

func (h *ListingsHandler) Get(c *gin.Context) {
	if !h.checkDB(c) {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	var listing models.Listing
	if err := h.DB.Preload("Images").
		Preload("Owner").
		First(&listing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	// Increment view count
	h.DB.Model(&listing).Update("view_count", listing.ViewCount+1)

	// Add price range to listing
	low := int64(float64(listing.Price) * 0.85)
	high := int64(float64(listing.Price) * 1.15)

	listingWithRange := gin.H{
		"id":                  listing.ID,
		"title":               listing.Title,
		"description":         listing.Description,
		"price":               listing.Price,
		"category":            listing.Category,
		"condition":           listing.Condition,
		"location":            listing.Location,
		"status":              listing.Status,
		"owner_id":            listing.OwnerID,
		"view_count":          listing.ViewCount,
		"created_at":          listing.CreatedAt,
		"updated_at":          listing.UpdatedAt,
		"brand_story":         listing.BrandStory,
		"rent":                listing.Rent,
		"floor":               listing.Floor,
		"equipment":           listing.Equipment,
		"decoration":          listing.Decoration,
		"annual_revenue":      listing.AnnualRevenue,
		"gross_profit_rate":   listing.GrossProfitRate,
		"fastest_moving_date": listing.FastestMovingDate,
		"phone_number":        listing.PhoneNumber,
		"square_meters":       listing.SquareMeters,
		"industry":            listing.Industry,
		"deposit":             listing.Deposit,
		"owner":               listing.Owner,
		"images":              listing.Images,
		"price_range": gin.H{
			"low":  low,
			"high": high,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"listing": listingWithRange,
	})
}

func (h *ListingsHandler) List(c *gin.Context) {
	if !h.checkDB(c) {
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	category := c.Query("category")
	location := c.Query("location")
	minPrice, _ := strconv.ParseInt(c.Query("min_price"), 10, 64)
	maxPrice, _ := strconv.ParseInt(c.Query("max_price"), 10, 64)
	condition := c.Query("condition")

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Build query
	query := h.DB.Model(&models.Listing{}).Where("status = ?", "活躍")

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}
	if minPrice > 0 {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("price <= ?", maxPrice)
	}
	if condition != "" {
		query = query.Where("condition = ?", condition)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get listings with pagination
	var listings []models.Listing
	if err := query.Preload("Images").
		Preload("Owner").
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(&listings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch listings"})
		return
	}

	// Add price ranges to listings
	listingsWithRanges := make([]gin.H, len(listings))
	for i, listing := range listings {
		low := int64(float64(listing.Price) * 0.85)
		high := int64(float64(listing.Price) * 1.15)

		listingsWithRanges[i] = gin.H{
			"id":                  listing.ID,
			"title":               listing.Title,
			"description":         listing.Description,
			"price":               listing.Price,
			"category":            listing.Category,
			"condition":           listing.Condition,
			"location":            listing.Location,
			"status":              listing.Status,
			"owner_id":            listing.OwnerID,
			"view_count":          listing.ViewCount,
			"created_at":          listing.CreatedAt,
			"updated_at":          listing.UpdatedAt,
			"brand_story":         listing.BrandStory,
			"rent":                listing.Rent,
			"floor":               listing.Floor,
			"equipment":           listing.Equipment,
			"decoration":          listing.Decoration,
			"annual_revenue":      listing.AnnualRevenue,
			"gross_profit_rate":   listing.GrossProfitRate,
			"fastest_moving_date": listing.FastestMovingDate,
			"phone_number":        listing.PhoneNumber,
			"square_meters":       listing.SquareMeters,
			"industry":            listing.Industry,
			"deposit":             listing.Deposit,
			"owner":               listing.Owner,
			"images":              listing.Images,
			"price_range": gin.H{
				"low":  low,
				"high": high,
			},
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"listings": listingsWithRanges,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (int(total) + limit - 1) / limit,
		},
	})
}

func (h *ListingsHandler) Update(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	var req listingUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if listing exists and user owns it
	var listing models.Listing
	if err := h.DB.Where("id = ? AND owner_id = ?", id, userID).First(&listing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found or access denied"})
		return
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Condition != nil {
		updates["condition"] = *req.Condition
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := h.DB.Model(&listing).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update listing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Listing updated successfully",
		"listing": listing,
	})
}

func (h *ListingsHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	// Check if listing exists and user owns it
	var listing models.Listing
	if err := h.DB.Where("id = ? AND owner_id = ?", id, userID).First(&listing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found or access denied"})
		return
	}

	// Soft delete by setting status to deleted
	if err := h.DB.Model(&listing).Update("status", "deleted").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete listing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing deleted successfully"})
}

func (h *ListingsHandler) UploadImages(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	// Check if listing exists and user owns it
	var listing models.Listing
	if err := h.DB.Where("id = ? AND owner_id = ?", id, userID).First(&listing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found or access denied"})
		return
	}

	// Handle file upload
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No images provided"})
		return
	}

	var uploadedImages []models.Image
	for i, file := range files {
		// Validate file type
		if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
			continue
		}

		// Generate filename
		filename := fmt.Sprintf("listing_%d_%d_%s", listing.ID, i, file.Filename)
		filepath := fmt.Sprintf("./uploads/%s", filename)

		// Save file
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			continue
		}

		// Create image record
		image := models.Image{
			ListingID: listing.ID,
			Filename:  filename,
			URL:       fmt.Sprintf("/uploads/%s", filename),
			Order:     i,
			IsPrimary: i == 0, // First image is primary
		}

		if err := h.DB.Create(&image).Error; err != nil {
			continue
		}

		uploadedImages = append(uploadedImages, image)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Uploaded %d images successfully", len(uploadedImages)),
		"images":  uploadedImages,
	})
}

func (h *ListingsHandler) GetCategories(c *gin.Context) {
	if !h.checkDB(c) {
		return
	}

	var categories []string
	h.DB.Model(&models.Listing{}).
		Where("status = ?", "active").
		Distinct().
		Pluck("category", &categories)

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}
