package models

import "time"

type Coupon struct {
	ID    uint    `gorm:"primaryKey"`
	Code  string  `gorm:"uniqueIndex;not null"`
	Type  string  `gorm:"type:varchar(20);not null"`
	Value float64 `gorm:"not null"`
	Quota int     `gorm:"not null"`

	UsedCount int `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
