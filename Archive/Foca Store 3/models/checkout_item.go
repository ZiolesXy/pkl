package models

import "gorm.io/gorm"

type CheckoutItem struct {
	gorm.Model
	CheckoutID uint
	ProductID  uint
	Quantity   int
	Price      int
	Product    Product
}
