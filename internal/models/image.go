package models

import "time"

type Image struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ListingID uint      `gorm:"index;not null" json:"listing_id"`
	Filename  string    `gorm:"size:255;not null" json:"filename"`
	URL       string    `gorm:"size:500;not null" json:"url"`
	AltText   string    `gorm:"size:255" json:"alt_text"`
	Order     int       `gorm:"default:0" json:"order"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relations
	Listing Listing `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
}
