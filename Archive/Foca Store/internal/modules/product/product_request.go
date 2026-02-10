package product

type CreateProductRequest struct{
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	Stock      int     `json:"stock" binding:"required,gte=0"`
	CategoryID uint    `json:"category_id" binding:"required"`
}