package models

import "time"

type Transaction struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	ListingID     uint       `gorm:"index;not null" json:"listing_id"`
	BuyerID       uint       `gorm:"index;not null" json:"buyer_id"`
	SellerID      uint       `gorm:"index;not null" json:"seller_id"`
	Amount        int64      `gorm:"not null" json:"amount"`
	Status        string     `gorm:"size:20;default:pending;index" json:"status"`
	PaymentMethod string     `gorm:"size:50" json:"payment_method"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Relations
	Listing Listing `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
	Buyer   User    `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Seller  User    `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
}
