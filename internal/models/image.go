package models

import "time"

type Image struct {
	ID        uint   `gorm:"primaryKey"`
	ListingID uint   `gorm:"index;not null"`
	URL       string `gorm:"size:512;not null"`
	CreatedAt time.Time
}
