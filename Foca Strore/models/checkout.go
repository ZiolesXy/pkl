package models

import (
	"time"

	"gorm.io/gorm"
)

type Checkout struct {
	ID     uint   `gorm:"primaryKey"`
	UID    string `gorm:"uniqueIndex;not null"`
	UserID uint   `gorm:"not null"`
	User   *User  `gorm:"foreignKey:UserID"`

	CouponID *uint
	Coupon   *Coupon `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	AddressID *uint
	Address   *Address `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	TotalPrice float64 `gorm:"not null"`
	Status     string  `gorm:"not null;type:varchar(20);default:'pending'"`

	Items []CheckoutItem `gorm:"foreignKey:CheckoutID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}