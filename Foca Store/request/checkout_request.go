package request

type UpdateCheckoutStatusRequest struct {
	Status string `json:"status"` // success | failed
}