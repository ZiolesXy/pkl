package response

import (
	"time"
	"voca-store/models"
)

type CheckoutItemResponse struct {
	ID       uint                `json:"id"`
	Quantity int                 `json:"quantity"`
	Price    float64             `json:"price"`
	Product  ProductMiniResponse `json:"product"`
}

type CheckoutDetailResponse struct {
	ID         uint                     `json:"id"`
	User       UserMiniResponse         `json:"user"`
	TotalPrice float64                  `json:"total_price"`
	Status     string                   `json:"status"`
	Coupon     *CouponResponse          `json:"coupon,omitempty"`
	Address    *CheckoutAddressResponse `json:"address,omitempty"`
	Items      []CheckoutItemResponse   `json:"items"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
}

type CheckoutListResponse struct {
	Entries []CheckoutDetailResponse `json:"entries"`
}

func BuildCheckoutDetailResponse(checkout models.Checkout) CheckoutDetailResponse {
	var items []CheckoutItemResponse

	for _, item := range checkout.Items {
		items = append(items, CheckoutItemResponse{
			ID:       item.ID,
			Quantity: item.Quantity,
			Price:    item.Price,
			Product: ProductMiniResponse{
				ID:   item.Product.ID,
				Name: item.Product.Name,
			},
		})
	}

	var coupon *CouponResponse
	if checkout.Coupon != nil {
		c := BuildCouponResponse(*checkout.Coupon)
		coupon = &c
	}

	var addressResp *CheckoutAddressResponse
	if checkout.Address != nil {
		a := BuildCheckoutAddressResponse(*checkout.Address)
		addressResp = &a
	}

	return CheckoutDetailResponse{
		ID: checkout.ID,
		User: UserMiniResponse{
			ID:    checkout.User.ID,
			Name:  checkout.User.Name,
			Email: checkout.User.Email,
		},
		Coupon:     coupon,
		Address:    addressResp,
		TotalPrice: checkout.TotalPrice,
		Status:     checkout.Status,
		Items:      items,
		CreatedAt:  checkout.CreatedAt,
		UpdatedAt:  checkout.UpdatedAt,
	}
}

func BuildCheckOutListResponse(checkouts []models.Checkout) CheckoutListResponse {
	response := []CheckoutDetailResponse{}

	for _, checkout := range checkouts {
		response = append(response, BuildCheckoutDetailResponse(checkout))
	}

	return CheckoutListResponse{
		Entries: response,
	}
}
