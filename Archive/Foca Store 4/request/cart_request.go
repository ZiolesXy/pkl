package request

type AddCartItemRequest struct {
	ProductID uint  `json:"product_id" binding:"required"`
	Quantity  int64 `json:"quantity" binding:"required,gt=0"`
}

type UpdateCartItemRequest struct {
	Quantity int64 `json:"quantity" binding:"required,gt=0"`
}

type UpdateCheckoutStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending success failed"`
}
