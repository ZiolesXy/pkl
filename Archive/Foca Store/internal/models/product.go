package models

type Product struct {
	BaseModel
	Name       string  `gorm:"not null"`
	Price      float64 `gorm:"not null"`
	Stock      int     `gorm:"not null"`
	CategoryID uint
	Category Category
}