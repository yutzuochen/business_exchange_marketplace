package models

import "time"

type Message struct {
	ID         uint   `gorm:"primaryKey"`
	FromUserID uint   `gorm:"index;not null"`
	ToUserID   uint   `gorm:"index;not null"`
	ListingID  uint   `gorm:"index"`
	Body       string `gorm:"type:text;not null"`
	CreatedAt  time.Time
}
