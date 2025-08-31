// Package auth provides JWT (JSON Web Token) authentication functionality
// for the Business Exchange Marketplace application.
//
// This package handles:
// - JWT token generation with user claims
// - JWT token parsing and validation
// - Token expiration management
// - Cross-service authentication compatibility
package auth

import (
	"errors"
	"time"

	"trade_company/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT payload structure containing user information
// and standard JWT registered claims.
//
// Fields:
//   - UserID: Unique identifier for the authenticated user
//   - Email: User's email address for identification
//   - RegisteredClaims: Standard JWT claims (issuer, expiration, etc.)
//
// The "uid" JSON tag ensures compatibility with the auction service
// which expects the user ID field to be named "uid".
type Claims struct {
	UserID uint   `json:"uid"`   // User identifier (compatible with auction service)
	Email  string `json:"email"` // User email address
	jwt.RegisteredClaims          // Standard JWT claims (iss, exp, iat, etc.)
}

// GenerateToken creates a new JWT token for an authenticated user.
//
// This function generates a signed JWT token containing the user's ID and email,
// along with standard claims like issuer, issued time, and expiration time.
//
// Parameters:
//   - cfg: Application configuration containing JWT settings
//   - userID: Unique identifier of the user to authenticate
//   - email: Email address of the user
//
// Returns:
//   - string: Signed JWT token string
//   - error: Any error that occurred during token generation
//
// The token is signed using HMAC-SHA256 algorithm and expires after
// the configured number of minutes (default: 60 minutes).
func GenerateToken(cfg *config.Config, userID uint, email string) (string, error) {
	// Create JWT claims with user information and metadata
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.JWTIssuer,                                                                           // Token issuer (typically service name)
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                          // Token creation time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWTExpireMinutes) * time.Minute)), // Token expiration time
		},
	}
	
	// Create and sign the token using HMAC-SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ParseToken validates and parses a JWT token string, returning the contained claims.
//
// This function verifies the token signature, checks expiration, and extracts
// the user information from the token payload.
//
// Parameters:
//   - cfg: Application configuration containing JWT secret for verification
//   - tokenString: The JWT token string to parse and validate
//
// Returns:
//   - *Claims: Parsed user claims if token is valid
//   - error: Authentication error if token is invalid, expired, or malformed
//
// Common errors:
//   - Token signature verification failure
//   - Token expiration
//   - Malformed token structure
//   - Invalid claims format
func ParseToken(cfg *config.Config, tokenString string) (*Claims, error) {
	// Parse and validate the token with our claims structure
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// Verify the signing method and return the secret key
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// Verify token validity (signature, expiration, etc.)
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	// Extract and validate claims structure
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}
	
	return claims, nil
}
