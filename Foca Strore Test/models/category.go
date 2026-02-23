package models

import "time"

type Category struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"unique;not null"`
    Slug      string    `gorm:"unique;not null"`
    Products  Product `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    CreatedAt time.Time
    UpdatedAt time.Time
}