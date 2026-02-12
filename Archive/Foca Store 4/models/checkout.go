package models

import "time"

type CheckoutStatus string

const (
	CheckoutPending CheckoutStatus = "pending"
	CheckoutSuccess CheckoutStatus = "success"
	CheckoutFailed  CheckoutStatus = "failed"
)

type Checkout struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	User        User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	TotalAmount int64          `gorm:"not null" json:"total_amount"`
	Status      CheckoutStatus `gorm:"size:20;not null" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
