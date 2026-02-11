package models

import "gorm.io/gorm"

type Checkout struct {
	gorm.Model
	UserID uint
	Status string // pending | success | failed
	Total  int
	Items  []CheckoutItem
}