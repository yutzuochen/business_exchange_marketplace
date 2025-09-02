// Package handlers provides HTTP request handlers for the Business Exchange Marketplace API.
// This file contains auction proxy handlers that forward requests to the auction service.
package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"trade_company/internal/config"
	"trade_company/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuctionProxyHandler handles proxy requests to the auction service.
// This allows the frontend to use HttpOnly cookies while still accessing auction functionality.
type AuctionProxyHandler struct {
	Cfg *config.Config // Configuration for auction service URL
	Log *zap.Logger    // Logger for proxy requests
}

// NewAuctionProxyHandler creates a new auction proxy handler.
func NewAuctionProxyHandler(cfg *config.Config, log *zap.Logger) *AuctionProxyHandler {
	return &AuctionProxyHandler{
		Cfg: cfg,
		Log: log,
	}
}

// getAuctionServiceURL returns the auction service base URL.
func (h *AuctionProxyHandler) getAuctionServiceURL() string {
	// Default to localhost for development
	return "http://127.0.0.1:8081"
}

// forwardRequest forwards a request to the auction service with proper authentication.
func (h *AuctionProxyHandler) forwardRequest(c *gin.Context, path string) {
	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		h.Log.Warn("Auction proxy request failed - no user ID in context",
			zap.String("ip", c.ClientIP()),
			zap.String("path", path))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userIDValue, ok := userID.(uint)
	if !ok {
		h.Log.Error("Auction proxy request failed - invalid user ID type in context",
			zap.String("ip", c.ClientIP()),
			zap.String("path", path),
			zap.Any("user_id_value", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Get the JWT token from the request context (set by JWT middleware)
	token, exists := c.Get("jwt_token")
	if !exists {
		h.Log.Warn("Auction proxy request failed - no JWT token in context",
			zap.String("ip", c.ClientIP()),
			zap.String("path", path),
			zap.Uint("user_id", userIDValue))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	tokenString, ok := token.(string)
	if !ok {
		h.Log.Error("Auction proxy request failed - invalid JWT token type in context",
			zap.String("ip", c.ClientIP()),
			zap.String("path", path),
			zap.Uint("user_id", userIDValue))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Build the target URL
	targetURL := h.getAuctionServiceURL() + path

	// Read the request body if present
	var bodyBytes []byte
	if c.Request.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(c.Request.Body)
		if err != nil {
			h.Log.Error("Auction proxy request failed - failed to read request body",
				zap.String("ip", c.ClientIP()),
				zap.String("path", path),
				zap.Uint("user_id", userIDValue),
				logger.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
			return
		}
	}

	// Create the request to the auction service
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		h.Log.Error("Auction proxy request failed - failed to create request",
			zap.String("ip", c.ClientIP()),
			zap.String("path", path),
			zap.Uint("user_id", userIDValue),
			logger.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	// Copy headers from the original request
	for key, values := range c.Request.Header {
		// Skip headers that shouldn't be forwarded
		if key == "Host" || key == "Content-Length" {
			continue
		}
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Set the Authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Make the request to the auction service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.Log.Error("Auction proxy request failed - failed to forward request",
			zap.String("ip", c.ClientIP()),
			zap.String("path", path),
			zap.Uint("user_id", userIDValue),
			logger.Err(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to connect to auction service"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.Log.Error("Auction proxy request failed - failed to read response body",
			zap.String("ip", c.ClientIP()),
			zap.String("path", path),
			zap.Uint("user_id", userIDValue),
			logger.Err(err))
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to read response"})
		return
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Set the response status and body
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)

	h.Log.Info("Auction proxy request completed",
		zap.String("ip", c.ClientIP()),
		zap.String("path", path),
		zap.Uint("user_id", userIDValue),
		zap.Int("status_code", resp.StatusCode))
}

// GetAuctions proxies GET /api/v1/auctions requests to the auction service.
func (h *AuctionProxyHandler) GetAuctions(c *gin.Context) {
	path := "/api/v1/auctions"
	if c.Request.URL.RawQuery != "" {
		path += "?" + c.Request.URL.RawQuery
	}
	h.forwardRequest(c, path)
}

// GetAuction proxies GET /api/v1/auctions/:id requests to the auction service.
func (h *AuctionProxyHandler) GetAuction(c *gin.Context) {
	auctionID := c.Param("id")
	path := fmt.Sprintf("/api/v1/auctions/%s", auctionID)
	h.forwardRequest(c, path)
}

// CreateAuction proxies POST /api/v1/auctions requests to the auction service.
func (h *AuctionProxyHandler) CreateAuction(c *gin.Context) {
	h.forwardRequest(c, "/api/v1/auctions")
}

// ActivateAuction proxies POST /api/v1/auctions/:id:activate requests to the auction service.
func (h *AuctionProxyHandler) ActivateAuction(c *gin.Context) {
	auctionID := c.Param("id")
	path := fmt.Sprintf("/api/v1/auctions/%s:activate", auctionID)
	h.forwardRequest(c, path)
}

// PlaceBid proxies POST /api/v1/auctions/:id/bids requests to the auction service.
func (h *AuctionProxyHandler) PlaceBid(c *gin.Context) {
	auctionID := c.Param("id")
	path := fmt.Sprintf("/api/v1/auctions/%s/bids", auctionID)
	h.forwardRequest(c, path)
}

// GetMyBids proxies GET /api/v1/auctions/:id/my-bids requests to the auction service.
func (h *AuctionProxyHandler) GetMyBids(c *gin.Context) {
	auctionID := c.Param("id")
	path := fmt.Sprintf("/api/v1/auctions/%s/my-bids", auctionID)
	h.forwardRequest(c, path)
}

// GetAuctionResults proxies GET /api/v1/auctions/:id/results requests to the auction service.
func (h *AuctionProxyHandler) GetAuctionResults(c *gin.Context) {
	auctionID := c.Param("id")
	path := fmt.Sprintf("/api/v1/auctions/%s/results", auctionID)
	h.forwardRequest(c, path)
}

// WebSocketProxy handles WebSocket connections to the auction service.
// This creates a WebSocket connection with the auction service and forwards messages.
func (h *AuctionProxyHandler) WebSocketProxy(c *gin.Context) {
	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		h.Log.Warn("WebSocket proxy request failed - no user ID in context",
			zap.String("ip", c.ClientIP()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userIDValue, ok := userID.(uint)
	if !ok {
		h.Log.Error("WebSocket proxy request failed - invalid user ID type in context",
			zap.String("ip", c.ClientIP()),
			zap.Any("user_id_value", userID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Get the JWT token from the request context
	token, exists := c.Get("jwt_token")
	if !exists {
		h.Log.Warn("WebSocket proxy request failed - no JWT token in context",
			zap.String("ip", c.ClientIP()),
			zap.Uint("user_id", userIDValue))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	tokenString, ok := token.(string)
	if !ok {
		h.Log.Error("WebSocket proxy request failed - invalid JWT token type in context",
			zap.String("ip", c.ClientIP()),
			zap.Uint("user_id", userIDValue))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	auctionID := c.Param("id")

	// For WebSocket, we need to return the WebSocket URL with the token
	// The frontend will connect directly to the auction service with this URL
	// Use the test endpoint that accepts query parameter tokens
	wsURL := fmt.Sprintf("ws://127.0.0.1:8081/ws/test/%s?token=%s", auctionID, tokenString)

	h.Log.Info("WebSocket proxy URL generated",
		zap.String("ip", c.ClientIP()),
		zap.String("auction_id", auctionID),
		zap.Uint("user_id", userIDValue))

	c.JSON(http.StatusOK, gin.H{
		"ws_url": wsURL,
	})
}
