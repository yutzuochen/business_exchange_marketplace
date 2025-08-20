package redisclient

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"trade_company/internal/models"
)

type CacheService struct {
	client *redis.Client
}

func NewCacheService(client *redis.Client) *CacheService {
	return &CacheService{client: client}
}

// Cache keys
const (
	ListingSearchKey = "listing:search:"
	ListingDetailKey = "listing:detail:"
	UserProfileKey   = "user:profile:"
	CategoryListKey  = "category:list"
)

// TTL constants
const (
	SearchResultTTL = 15 * time.Minute
	ListingDetailTTL = 30 * time.Minute
	UserProfileTTL = 1 * time.Hour
	CategoryListTTL = 24 * time.Hour
)

// CacheListingSearch caches search results
func (c *CacheService) CacheListingSearch(query string, filters map[string]interface{}, results []models.Listing) error {
	key := fmt.Sprintf("%s%s", ListingSearchKey, hashQuery(query, filters))
	
	data, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("failed to marshal search results: %w", err)
	}
	
	ctx := context.Background()
	return c.client.Set(ctx, key, data, SearchResultTTL).Err()
}

// GetCachedListingSearch retrieves cached search results
func (c *CacheService) GetCachedListingSearch(query string, filters map[string]interface{}) ([]models.Listing, error) {
	key := fmt.Sprintf("%s%s", ListingSearchKey, hashQuery(query, filters))
	
	ctx := context.Background()
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get cached search results: %w", err)
	}
	
	var results []models.Listing
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached search results: %w", err)
	}
	
	return results, nil
}

// CacheListingDetail caches individual listing details
func (c *CacheService) CacheListingDetail(listingID uint, listing *models.Listing) error {
	key := fmt.Sprintf("%s%d", ListingDetailKey, listingID)
	
	data, err := json.Marshal(listing)
	if err != nil {
		return fmt.Errorf("failed to marshal listing: %w", err)
	}
	
	ctx := context.Background()
	return c.client.Set(ctx, key, data, ListingDetailTTL).Err()
}

// GetCachedListingDetail retrieves cached listing details
func (c *CacheService) GetCachedListingDetail(listingID uint) (*models.Listing, error) {
	key := fmt.Sprintf("%s%d", ListingDetailKey, listingID)
	
	ctx := context.Background()
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get cached listing: %w", err)
	}
	
	var listing models.Listing
	if err := json.Unmarshal(data, &listing); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached listing: %w", err)
	}
	
	return &listing, nil
}

// InvalidateListingCache invalidates all listing-related caches
func (c *CacheService) InvalidateListingCache(listingID uint) error {
	ctx := context.Background()
	
	// Invalidate listing detail cache
	detailKey := fmt.Sprintf("%s%d", ListingDetailKey, listingID)
	if err := c.client.Del(ctx, detailKey).Err(); err != nil {
		return fmt.Errorf("failed to invalidate listing detail cache: %w", err)
	}
	
	// Invalidate all search caches (pattern matching)
	pattern := fmt.Sprintf("%s*", ListingSearchKey)
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get search cache keys: %w", err)
	}
	
	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to invalidate search caches: %w", err)
		}
	}
	
	return nil
}

// InvalidateUserCache invalidates user-related caches
func (c *CacheService) InvalidateUserCache(userID uint) error {
	ctx := context.Background()
	
	// Invalidate user profile cache
	profileKey := fmt.Sprintf("%s%d", UserProfileKey, userID)
	if err := c.client.Del(ctx, profileKey).Err(); err != nil {
		return fmt.Errorf("failed to invalidate user profile cache: %w", err)
	}
	
	return nil
}

// ClearAllCaches clears all caches (use with caution)
func (c *CacheService) ClearAllCaches() error {
	ctx := context.Background()
	return c.client.FlushDB(ctx).Err()
}

// GetCacheStats returns cache statistics
func (c *CacheService) GetCacheStats() (map[string]interface{}, error) {
	ctx := context.Background()
	
	info, err := c.client.Info(ctx, "memory").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache stats: %w", err)
	}
	
	// Parse Redis INFO output for memory usage
	stats := map[string]interface{}{
		"info": info,
	}
	
	return stats, nil
}

// hashQuery creates a hash for the search query and filters
func hashQuery(query string, filters map[string]interface{}) string {
	// Simple hash implementation - in production, use a proper hash function
	return fmt.Sprintf("%s_%v", query, filters)
}
