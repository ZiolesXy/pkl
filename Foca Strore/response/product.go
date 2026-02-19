package response

import (
	"time"
	"voca-store/models"
)

type CategoryMiniResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProductMiniResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProductResponse struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Slug        string                `json:"slug"`
	Description string                `json:"description,omitempty"`
	ImageURL    string                `json:"image_url,omitempty"`
	Price       float64               `json:"price"`
	Stock       int                   `json:"stock"`
	Category    *CategoryMiniResponse `json:"category,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type ProductListResponse struct {
	Entries []ProductResponse `json:"entries"`
}

func BuildProductResponse(product models.Product) ProductResponse {
	var CategoryResp *CategoryMiniResponse
	if product.Category != nil {
		CategoryResp = &CategoryMiniResponse{
			ID:   product.Category.ID,
			Name: product.Category.Name,
		}
	}

	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		ImageURL:    product.ImageURL,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    CategoryResp,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func BuildProductListResponse(products []models.Product) ProductListResponse {
	responses := []ProductResponse{}

	for _, p := range products {
		var CategoryResp *CategoryMiniResponse
		if p.Category != nil {
			CategoryResp = &CategoryMiniResponse{
				ID:   p.Category.ID,
				Name: p.Category.Name,
			}
		}

		responses = append(responses, ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description,
			ImageURL:    p.ImageURL,
			Price:       p.Price,
			Stock:       p.Stock,
			Category:    CategoryResp,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return ProductListResponse{
		Entries: responses,
	}
}
