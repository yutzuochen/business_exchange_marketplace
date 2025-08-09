package models

import "time"

type Favorite struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index;not null"`
	ListingID uint `gorm:"index;not null"`
	CreatedAt time.Time
}
