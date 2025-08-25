package models

import (
	"errors"
	"time"
)

type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Email        string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Username     string     `gorm:"uniqueIndex;size:100;not null" json:"username"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	FirstName    string     `gorm:"size:100" json:"first_name"`
	LastName     string     `gorm:"size:100" json:"last_name"`
	Phone        string     `gorm:"size:20" json:"phone"`
	Role         string     `gorm:"size:32;not null;default:user;index" json:"role"`
	IsActive     bool       `gorm:"default:true;index" json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Email verification
	EmailVerifiedAt        *time.Time `gorm:"index" json:"email_verified_at,omitempty"`
	EmailVerificationToken string     `gorm:"size:255" json:"-"`

	// 2FA support
	TwoFactorEnabled bool   `gorm:"default:false" json:"two_factor_enabled"`
	TwoFactorSecret  string `gorm:"size:255" json:"-"`

	// Seller-specific fields
	CompanyName  string `gorm:"size:255" json:"company_name,omitempty"`
	TaxID        string `gorm:"size:20" json:"tax_id,omitempty"` // 統一編號
	ContactPhone string `gorm:"size:20" json:"contact_phone,omitempty"`

	// Notification preferences
	EmailNotifications bool `gorm:"default:true" json:"email_notifications"`
	MarketingEmails    bool `gorm:"default:false" json:"marketing_emails"`

	// Relations
	Listings         []Listing     `gorm:"foreignKey:OwnerID" json:"listings,omitempty"`
	Favorites        []Favorite    `gorm:"foreignKey:UserID" json:"favorites,omitempty"`
	SentMessages     []Message     `gorm:"foreignKey:SenderID" json:"sent_messages,omitempty"`
	ReceivedMessages []Message     `gorm:"foreignKey:ReceiverID" json:"received_messages,omitempty"`
	Transactions     []Transaction `gorm:"foreignKey:BuyerID" json:"transactions,omitempty"`

	// Session management
	Sessions []UserSession `gorm:"foreignKey:UserID" json:"sessions,omitempty"`

	// Lead management
	ReceivedLeads []Lead `gorm:"foreignKey:ReceiverID" json:"received_leads,omitempty"`
}

// UserSession represents user login sessions
type UserSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	SessionID string    `gorm:"size:255;not null;uniqueIndex" json:"session_id"`
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Lead represents contact form submissions from buyers to sellers
type Lead struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SenderID     uint      `gorm:"not null;index" json:"sender_id"`
	ReceiverID   uint      `gorm:"not null;index" json:"receiver_id"`
	ListingID    *uint     `gorm:"index" json:"listing_id,omitempty"`
	Subject      string    `gorm:"size:255;not null" json:"subject"`
	Message      string    `gorm:"type:text;not null" json:"message"`
	ContactPhone string    `gorm:"size:20" json:"contact_phone,omitempty"`
	IsRead       bool      `gorm:"default:false;index" json:"is_read"`
	IsSpam       bool      `gorm:"default:false;index" json:"is_spam"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Sender   User     `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Receiver User     `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
	Listing  *Listing `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
}

// PasswordResetToken represents password reset tokens
type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Token     string    `gorm:"size:255;not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// AuditLog represents security audit events
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id,omitempty"`
	Event     string    `gorm:"size:100;not null" json:"event"`
	Details   string    `gorm:"type:text" json:"details"`
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

var ErrPlaceholder = errors.New("placeholder")
