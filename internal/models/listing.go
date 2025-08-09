package models

import "time"

type Listing struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	PriceUSD    int64  `gorm:"not null"`
	Location    string `gorm:"size:255"`
	OwnerID     uint   `gorm:"index;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Images      []Image
}
