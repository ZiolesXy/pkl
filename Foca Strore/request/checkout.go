package request

type UpdateCheckoutStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending success failed"`
}