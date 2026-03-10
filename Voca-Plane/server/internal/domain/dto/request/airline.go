package request

type CreateAirlineRequest struct {
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code" binding:"required"`
	LogoURL string `json:"logo_url"`
}

type UpdateAirlineRequest struct {
	Name    *string `json:"name,omitempty"`
	Code    *string `json:"code,omitempty"`
	LogoURL *string `json:"logo_url,omitempty"`
}