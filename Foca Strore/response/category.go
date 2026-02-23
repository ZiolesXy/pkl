package response

import (
	"time"
	"voca-store/models"
)

type CategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryListResponse struct {
	Entries []CategoryResponse `json:"entries"`
}

func BuildCategoryResponse(category models.Category) CategoryResponse {
	return CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func BuildCategoryListResponse(categories []models.Category) CategoryListResponse {
	response := []CategoryResponse{}

	for _, c := range categories {
		response = append(response, BuildCategoryResponse(c))
	}

	return CategoryListResponse{
		Entries: response,
	}
}
