package models

import "time"

type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	ListingID uint      `gorm:"index;not null" json:"listing_id"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relations
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Listing Listing `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
}
