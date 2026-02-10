package models

type CartItem struct{
	BaseModel
	CartID uint
	ProductID uint
	Product Product
	Quantity int
}