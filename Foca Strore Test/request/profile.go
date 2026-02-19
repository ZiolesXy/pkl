package request

type UpdateProfileRequest struct {
	Name *string `json:"name"`
}

type UpdateProfileWithImageRequest struct {
	Name *string `json:"name"`
}
