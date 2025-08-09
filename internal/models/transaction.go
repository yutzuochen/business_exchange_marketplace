package models

import "time"

type Transaction struct {
	ID        uint   `gorm:"primaryKey"`
	BuyerID   uint   `gorm:"index;not null"`
	SellerID  uint   `gorm:"index;not null"`
	ListingID uint   `gorm:"index;not null"`
	AmountUSD int64  `gorm:"not null"`
	Status    string `gorm:"size:32;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
