package request

type ProductCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int64  `json:"price" binding:"required,gt=0"`
	Stock       int64  `json:"stock" binding:"required,gte=0"`
}

type ProductUpdateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int64  `json:"price" binding:"required,gt=0"`
	Stock       int64  `json:"stock" binding:"required,gte=0"`
}
