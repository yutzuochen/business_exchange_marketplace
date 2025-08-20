package models

import "time"

type Listing struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null;index" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	PriceUSD    int64     `gorm:"not null;index" json:"price_usd"`
	Category    string    `gorm:"size:100;index" json:"category"`
	Condition   string    `gorm:"size:50;default:used" json:"condition"`
	Location    string    `gorm:"size:255;index" json:"location"`
	Status      string    `gorm:"size:20;default:active;index" json:"status"`
	OwnerID     uint      `gorm:"index;not null" json:"owner_id"`
	ViewCount   int       `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	Owner   User     `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Images  []Image  `gorm:"foreignKey:ListingID" json:"images,omitempty"`
	Favorites []Favorite `gorm:"foreignKey:ListingID" json:"favorites,omitempty"`
}
