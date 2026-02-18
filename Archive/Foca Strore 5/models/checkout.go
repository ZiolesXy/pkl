package models

import "time"

type Checkout struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null"`
	User       *User      `gorm:"foreignKey:UserID"`
	TotalPrice float64   `gorm:"not null"`
	Status     string    `gorm:"not null;type:varchar(20);default:'pending'"`
	Items []CheckoutItem
	CreatedAt  time.Time 
	UpdatedAt  time.Time
}