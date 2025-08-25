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

type SessionManager struct {
	redisClient *redis.Client
	db          *gorm.DB
	config      *config.Config
}

func NewSessionManager(redisClient *redis.Client, db *gorm.DB, config *config.Config) *SessionManager {
	return &SessionManager{
		redisClient: redisClient,
		db:          db,
		config:      config,
	}
}

// CreateSession creates a new session for a user
func (sm *SessionManager) CreateSession(userID uint, ipAddress, userAgent string) (*models.UserSession, error) {
	// Generate unique session ID
	sessionID, err := sm.generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(sm.config.SessionTTLMinutes) * time.Minute)

	// Create session in database
	session := &models.UserSession{
		UserID:    userID,
		SessionID: sessionID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ExpiresAt: expiresAt,
	}

	if err := sm.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session in database: %w", err)
	}

	// Store session in Redis for fast access
	key := fmt.Sprintf("session:%s", sessionID)
	sessionData := map[string]interface{}{
		"user_id":    userID,
		"ip_address": ipAddress,
		"user_agent": userAgent,
		"expires_at": expiresAt.Unix(),
	}

	// Set expiration in Redis
	ttl := time.Until(expiresAt)
	if err := sm.redisClient.HMSet(ctx, key, sessionData).Err(); err != nil {
		return nil, fmt.Errorf("failed to store session in Redis: %w", err)
	}
	if err := sm.redisClient.Expire(ctx, key, ttl).Err(); err != nil {
		return nil, fmt.Errorf("failed to set Redis expiration: %w", err)
	}

	return session, nil
}

// GetSession retrieves a session by session ID
func (sm *SessionManager) GetSession(sessionID string) (*models.UserSession, error) {
	// Try Redis first
	key := fmt.Sprintf("session:%s", sessionID)
	sessionData, err := sm.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session from Redis: %w", err)
	}

	if len(sessionData) == 0 {
		// Try database as fallback
		var session models.UserSession
		if err := sm.db.Where("session_id = ? AND expires_at > ?", sessionID, time.Now()).First(&session).Error; err != nil {
			return nil, fmt.Errorf("session not found or expired: %w", err)
		}
		return &session, nil
	}

	// Parse session data from Redis
	userID, err := parseUint(sessionData["user_id"])
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in session: %w", err)
	}

	expiresAt := time.Unix(parseInt64(sessionData["expires_at"]), 0)
	if time.Now().After(expiresAt) {
		// Session expired, remove it
		sm.RevokeSession(sessionID)
		return nil, fmt.Errorf("session expired")
	}

	session := &models.UserSession{
		UserID:    userID,
		SessionID: sessionID,
		IPAddress: sessionData["ip_address"],
		UserAgent: sessionData["user_agent"],
		ExpiresAt: expiresAt,
	}

	return session, nil
}

// RevokeSession removes a session
func (sm *SessionManager) RevokeSession(sessionID string) error {
	// Remove from Redis
	key := fmt.Sprintf("session:%s", sessionID)
	if err := sm.redisClient.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to remove session from Redis: %w", err)
	}

	// Remove from database
	if err := sm.db.Where("session_id = ?", sessionID).Delete(&models.UserSession{}).Error; err != nil {
		return fmt.Errorf("failed to remove session from database: %w", err)
	}

	return nil
}

// GetUserSessions returns all active sessions for a user
func (sm *SessionManager) GetUserSessions(userID uint) ([]models.UserSession, error) {
	var sessions []models.UserSession
	if err := sm.db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	return sessions, nil
}

// CleanupExpiredSessions removes expired sessions
func (sm *SessionManager) CleanupExpiredSessions() error {
	// Clean up database
	if err := sm.db.Where("expires_at <= ?", time.Now()).Delete(&models.UserSession{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup expired sessions in database: %w", err)
	}

	// Redis will automatically expire keys, but we can also clean up manually if needed
	return nil
}

// generateSessionID generates a cryptographically secure random session ID
func (sm *SessionManager) generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Helper functions for parsing Redis data
func parseUint(s string) (uint, error) {
	var result uint64
	_, err := fmt.Sscanf(s, "%d", &result)
	return uint(result), err
}

func parseInt64(s string) int64 {
	var result int64
	fmt.Sscanf(s, "%d", &result)
	return result
}

// Context for Redis operations
var ctx = context.Background()
