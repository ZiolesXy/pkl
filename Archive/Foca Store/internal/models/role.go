package models

type Role struct {
	BaseModel
	Name string `gorm:"unique;not null"`
}