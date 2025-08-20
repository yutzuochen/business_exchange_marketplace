package models

import (
	"errors"
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Username     string    `gorm:"uniqueIndex;size:100;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	FirstName   string    `gorm:"size:100" json:"first_name"`
	LastName    string    `gorm:"size:100" json:"last_name"`
	Phone       string    `gorm:"size:20" json:"phone"`
	Role        string    `gorm:"size:32;not null;default:user;index" json:"role"`
	IsActive    bool      `gorm:"default:true;index" json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	
	// Relations
	Listings     []Listing     `gorm:"foreignKey:OwnerID" json:"listings,omitempty"`
	Favorites   []Favorite    `gorm:"foreignKey:UserID" json:"favorites,omitempty"`
	SentMessages     []Message `gorm:"foreignKey:SenderID" json:"sent_messages,omitempty"`
	ReceivedMessages []Message `gorm:"foreignKey:ReceiverID" json:"received_messages,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:BuyerID" json:"transactions,omitempty"`
}

var ErrPlaceholder = errors.New("placeholder")
