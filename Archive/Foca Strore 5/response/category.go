package response

import "time"

type CategoryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryListResponse struct {
	Entries []CategoryResponse `json:"entries"`
}

func BuildCategoryResponse(
	id uint,
	name string,
	createdAt, updatedAt time.Time,
) CategoryResponse {
	return CategoryResponse{
		ID: id,
		Name: name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func BuildCategoryListResponse(categories []CategoryResponse) CategoryListResponse{
	return CategoryListResponse{
		Entries: categories,
	}
}