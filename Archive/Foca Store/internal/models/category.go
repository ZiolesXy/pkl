package models

type Category struct {
	BaseModel
	Name string `gorm:"unique;not null"`
	Product []Product
}