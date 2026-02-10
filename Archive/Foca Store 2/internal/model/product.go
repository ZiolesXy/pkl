package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string    `gorm:"type:varchar(100);not null" json:"name"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

type Product struct {
	gorm.Model
	Name        string   `gorm:"type:varchar(255);not null" json:"name"`
	Description string   `gorm:"type:text" json:"description"`
	Price       float64  `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int      `gorm:"type:int;not null" json:"stock"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// DTO untuk Request
type ProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

type CategoryInput struct {
	Name string `json:"name" binding:"required"`
}