package models

import "time"

type Product struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	Slug          string `gorm:"uniqueIndex;not null"`
	Description   string
	ImageURL      string
	ImagePublicID string
	Price         float64 `gorm:"not null"`
	Stock         int
	CategoryID    uint
	Category      *Category
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CartItems     []CartItem `gorm:"foreignKey:ProductID"`
}
