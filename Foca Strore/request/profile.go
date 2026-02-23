package request

type UpdateProfileRequest struct {
	Name       *string `json:"name"`
	Phone      *string `json:"phone"`
	Address    *string `json:"address"`
	PostalCode *string `json:"postal_code"`
}

type UpdateProfileWithImageRequest struct {
	Name       *string `json:"name"`
	Phone      *string `json:"phone"`
	Address    *string `json:"address"`
	PostalCode *string `json:"postal_code"`
}
