// Package auth provides session management functionality for user authentication
// in the Business Exchange Marketplace application.
//
// The SessionManager implements a hybrid session storage approach:
// - Redis for fast session lookups and automatic expiration
// - MySQL for persistent storage and session history
// - Fallback mechanism when Redis is unavailable
package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"trade_company/internal/config"
	"trade_company/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SessionManager handles user session lifecycle management with dual storage.
//
// Architecture:
//   - Primary storage: Redis (fast access, automatic expiration)
//   - Secondary storage: MySQL (persistent, backup, audit trail)
//   - Graceful degradation when Redis is unavailable
//
// Features:
//   - Cryptographically secure session ID generation
//   - Automatic session expiration
//   - IP address and User-Agent tracking for security
//   - Bulk session management (get all user sessions)
//   - Cleanup utilities for expired sessions
type SessionManager struct {
	redisClient *redis.Client  // Redis client for fast session storage
	db          *gorm.DB       // Database connection for persistent storage
	config      *config.Config // Application configuration
}

// NewSessionManager creates a new session manager instance with dependencies.
//
// Parameters:
//   - redisClient: Redis client for caching (can be nil for database-only mode)
//   - db: GORM database connection for persistent storage
//   - config: Application configuration containing session settings
//
// Returns:
//   - *SessionManager: Configured session manager ready for use
func NewSessionManager(redisClient *redis.Client, db *gorm.DB, config *config.Config) *SessionManager {
	return &SessionManager{
		redisClient: redisClient,
		db:          db,
		config:      config,
	}
}

// CreateSession creates a new authenticated session for a user after login.
//
// This method generates a cryptographically secure session ID and stores
// the session in both Redis (for performance) and MySQL (for persistence).
// The session includes security metadata like IP address and User-Agent.
//
// Parameters:
//   - userID: Unique identifier of the authenticated user
//   - ipAddress: Client IP address for security tracking
//   - userAgent: Client browser/app identifier for security tracking
//
// Returns:
//   - *models.UserSession: Created session object with all metadata
//   - error: Any error that occurred during session creation
//
// Security considerations:
//   - Uses crypto/rand for secure session ID generation
//   - Tracks IP and User-Agent for session hijacking detection
//   - Automatic expiration based on configured TTL
func (sm *SessionManager) CreateSession(userID uint, ipAddress, userAgent string) (*models.UserSession, error) {
	// Generate cryptographically secure session ID (64 character hex string)
	sessionID, err := sm.generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Calculate session expiration time based on configuration
	expiresAt := time.Now().Add(time.Duration(sm.config.SessionTTLMinutes) * time.Minute)

	// Create session record for persistent storage
	session := &models.UserSession{
		UserID:    userID,
		SessionID: sessionID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ExpiresAt: expiresAt,
	}

	// Store session in MySQL database for persistence and audit trail
	if err := sm.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session in database: %w", err)
	}

	// Store session in Redis for fast lookup (if Redis is available)
	if sm.redisClient != nil {
		key := fmt.Sprintf("session:%s", sessionID)
		sessionData := map[string]interface{}{
			"user_id":    userID,
			"ip_address": ipAddress,
			"user_agent": userAgent,
			"expires_at": expiresAt.Unix(),
		}

		// Store session data as Redis hash with automatic expiration
		ttl := time.Until(expiresAt)
		if err := sm.redisClient.HMSet(ctx, key, sessionData).Err(); err != nil {
			return nil, fmt.Errorf("failed to store session in Redis: %w", err)
		}
		if err := sm.redisClient.Expire(ctx, key, ttl).Err(); err != nil {
			return nil, fmt.Errorf("failed to set Redis expiration: %w", err)
		}
	}

	return session, nil
}

// GetSession retrieves and validates a session by session ID.
//
// This method implements a tiered lookup strategy:
// 1. Check Redis first for fastest access
// 2. Fallback to MySQL database if Redis miss
// 3. Validate expiration and clean up expired sessions
//
// Parameters:
//   - sessionID: Unique session identifier to retrieve
//
// Returns:
//   - *models.UserSession: Session object if valid and not expired
//   - error: Session not found, expired, or lookup error
//
// Performance characteristics:
//   - Redis hit: ~1ms response time
//   - Database fallback: ~10-50ms response time
//   - Automatic cleanup of expired sessions
func (sm *SessionManager) GetSession(sessionID string) (*models.UserSession, error) {
	// Primary lookup: Try Redis for fastest access
	if sm.redisClient != nil {
		key := fmt.Sprintf("session:%s", sessionID)
		sessionData, err := sm.redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get session from Redis: %w", err)
		}

		if len(sessionData) > 0 {
			// Parse and validate Redis session data
			userID, err := parseUint(sessionData["user_id"])
			if err != nil {
				return nil, fmt.Errorf("invalid user ID in session: %w", err)
			}

			expiresAt := time.Unix(parseInt64(sessionData["expires_at"]), 0)
			if time.Now().After(expiresAt) {
				// Session expired in Redis, clean it up
				sm.RevokeSession(sessionID)
				return nil, fmt.Errorf("session expired")
			}

			// Reconstruct session object from Redis data
			session := &models.UserSession{
				UserID:    userID,
				SessionID: sessionID,
				IPAddress: sessionData["ip_address"],
				UserAgent: sessionData["user_agent"],
				ExpiresAt: expiresAt,
			}

			return session, nil
		}
	}

	// Fallback lookup: Try MySQL database
	var session models.UserSession
	if err := sm.db.Where("session_id = ? AND expires_at > ?", sessionID, time.Now()).First(&session).Error; err != nil {
		return nil, fmt.Errorf("session not found or expired: %w", err)
	}

	return &session, nil
}

// RevokeSession immediately invalidates and removes a session from both stores.
// Used for logout, security incidents, or administrative actions.
func (sm *SessionManager) RevokeSession(sessionID string) error {
	// Remove from Redis cache
	if sm.redisClient != nil {
		key := fmt.Sprintf("session:%s", sessionID)
		if err := sm.redisClient.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("failed to remove session from Redis: %w", err)
		}
	}

	// Remove from MySQL database
	if err := sm.db.Where("session_id = ?", sessionID).Delete(&models.UserSession{}).Error; err != nil {
		return fmt.Errorf("failed to remove session from database: %w", err)
	}

	return nil
}

// GetUserSessions returns all active (non-expired) sessions for a specific user.
// Useful for security dashboards and "log out from all devices" functionality.
func (sm *SessionManager) GetUserSessions(userID uint) ([]models.UserSession, error) {
	var sessions []models.UserSession
	if err := sm.db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	return sessions, nil
}

// CleanupExpiredSessions removes expired sessions from the database.
// Should be called periodically (e.g., via cron job) for housekeeping.
// Redis sessions expire automatically, but database cleanup requires manual intervention.
func (sm *SessionManager) CleanupExpiredSessions() error {
	if err := sm.db.Where("expires_at <= ?", time.Now()).Delete(&models.UserSession{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup expired sessions in database: %w", err)
	}
	return nil
}

// generateSessionID creates a cryptographically secure 64-character hex session ID.
// Uses crypto/rand for secure random number generation to prevent session prediction attacks.
func (sm *SessionManager) generateSessionID() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits of entropy
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil // Returns 64-character hex string
}

// parseUint safely parses a string to uint, handling Redis string data.
func parseUint(s string) (uint, error) {
	var result uint64
	_, err := fmt.Sscanf(s, "%d", &result)
	return uint(result), err
}

// parseInt64 safely parses a string to int64 for timestamp handling.
func parseInt64(s string) int64 {
	var result int64
	fmt.Sscanf(s, "%d", &result)
	return result
}

// ctx provides context for Redis operations (background context is sufficient for session operations).
var ctx = context.Background()
