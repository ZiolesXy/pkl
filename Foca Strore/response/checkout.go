package response

import "time"

type CheckoutResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	UserName   string    `json:"user_name,omitempty"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CheckoutListResponse struct {
	Entries []CheckoutResponse `json:"entries"`
}

func BuildCheckoutResponse(id, userID uint, userName string, totalPrice float64, status string, createdAt, updatedAt time.Time) CheckoutResponse {
	return CheckoutResponse{
		ID:         id,
		UserID:     userID,
		UserName:   userName,
		TotalPrice: totalPrice,
		Status:     status,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func BuildCheckoutListResponse(checkout []CheckoutResponse) CheckoutListResponse {
	return CheckoutListResponse{
		Entries: checkout,
	}
}