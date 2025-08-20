package models

import "time"

type Message struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SenderID    uint      `gorm:"index;not null" json:"sender_id"`
	ReceiverID  uint      `gorm:"index;not null" json:"receiver_id"`
	ListingID   *uint     `gorm:"index" json:"listing_id,omitempty"`
	Subject     string    `gorm:"size:255" json:"subject"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	IsRead      bool      `gorm:"default:false;index" json:"is_read"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relations
	Sender   User    `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Receiver User    `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
	Listing  *Listing `gorm:"foreignKey:ListingID" json:"listing,omitempty"`
}
