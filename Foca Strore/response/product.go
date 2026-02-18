package response

import (
	"time"
	"voca-store/models"
)

type ProductMiniResponse struct {
	ID uint `json:"id"`
	Name string `json:"name"`
}

type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductListResponse struct {
	Entries []ProductResponse `json:"entries"`
}

func BuildProductResponse(id uint, name, slug, description, imageURL string, price float64, stock int, createdAt, updatedAt time.Time) ProductResponse {
	return ProductResponse{
		ID:          id,
		Name:        name,
		Slug:        slug,
		Description: description,
		ImageURL:    imageURL,
		Price:       price,
		Stock:       stock,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func BuildProductListResponse(products []models.Product) ProductListResponse {
	var responses []ProductResponse

	for _, p := range products {
		responses = append(responses, ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description,
			ImageURL:    p.ImageURL,
			Price:       p.Price,
			Stock:       p.Stock,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return ProductListResponse{
		Entries: responses,
	}
}