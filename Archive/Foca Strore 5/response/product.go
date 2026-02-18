package response

import "time"

type CategorySimpleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProductResponse struct {
	ID          uint                   `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	ImageURL    string                 `json:"image_url,omitempty"`
	Price       float64                `json:"price"`
	Stock       int                    `json:"stock"`
	Category    CategorySimpleResponse `json:"category"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ProductListResponse struct {
	Entries []ProductResponse `json:"entries"`
}

func BuildProductResponse(id uint, name, description, imageURL string, price float64, stock int, categoryID uint, categoryName string, createdAt, updatedAt time.Time) ProductResponse {
	return ProductResponse{
		ID:          id,
		Name:        name,
		Description: description,
		ImageURL:    imageURL,
		Price:       price,
		Stock:       stock,
		Category: CategorySimpleResponse{
			ID: categoryID,
			Name: categoryName,
		},
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func BuildProductListResponse(products []ProductResponse) ProductListResponse {
	return ProductListResponse{
		Entries: products,
	}
}
