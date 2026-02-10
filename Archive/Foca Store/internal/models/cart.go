package models

type Cart struct {
	BaseModel
	UserID uint
	Items []CartItem
}