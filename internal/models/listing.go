package models

import "time"

type Listing struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Title             string    `gorm:"size:255;not null;index" json:"title"`
	Description       string    `gorm:"type:text" json:"description"`
	Price             int64     `gorm:"not null;index" json:"price"`
	Category          string    `gorm:"size:100;index" json:"category"`
	Condition         string    `gorm:"size:50;default:used" json:"condition"`
	Location          string    `gorm:"size:255;index" json:"location"`
	Status            string    `gorm:"size:20;default:active;index" json:"status"`
	OwnerID           uint      `gorm:"index;not null" json:"owner_id"`
	ViewCount         int       `gorm:"default:0" json:"view_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	BrandStory        string    `gorm:"type:text" json:"brand_story,omitempty"`
	Rent              int64     `gorm:"index" json:"rent,omitempty"`
	Floor             int       `json:"floor,omitempty"`
	Equipment         string    `gorm:"type:text" json:"equipment,omitempty"`
	Decoration        string    `gorm:"size:100" json:"decoration,omitempty"`
	AnnualRevenue     int64     `json:"annual_revenue,omitempty"`
	GrossProfitRate   float64   `json:"gross_profit_rate,omitempty"`
	FastestMovingDate time.Time `json:"fastest_moving_date,omitempty"`
	PhoneNumber       string    `gorm:"size:20" json:"phone_number,omitempty"`
	SquareMeters      float64   `json:"square_meters,omitempty"`
	Industry          string    `gorm:"size:100;index" json:"industry,omitempty"`
	Deposit           int64     `json:"deposit,omitempty"`
	// Relations
	Owner     User       `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Images    []Image    `gorm:"foreignKey:ListingID" json:"images,omitempty"`
	Favorites []Favorite `gorm:"foreignKey:ListingID" json:"favorites,omitempty"`
}
