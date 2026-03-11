package response

import (
	"time"
	"voca-plane/internal/domain/models"
)

type TransactionItemResponse struct {
	PassengerName string `json:"passenger_name"`
	Nationality   string `json:"nationality"`
	PassportNo    string `json:"passport_no"`
	SeatNumber    string `json:"seat_number"`
	ClassName     string `json:"class_name"`
	Price         float64 `json:"price"`
}

type TransactionResponse struct {
	Code          string                    `json:"code"`
	FlightNumber  string                    `json:"flight_number"`
	TotalPrice    float64                   `json:"total_price"`
	PaymentStatus string                    `json:"payment_status"`
	PaymentURL    string                    `json:"payment_url"`
	PromoCode     *string                   `json:"promo_code,omitempty"`
	Discount      float64                   `json:"discount"`
	ExpiresAt     time.Time                 `json:"expires_at"`
	CreatedAt     time.Time                 `json:"created_at"`
	Items         []TransactionItemResponse `json:"items,omitempty"`
}

func ToTransactionItemResponse(p models.TransactionItem) TransactionItemResponse {
	return TransactionItemResponse{
		PassengerName: p.PassengerName,
		Nationality:   p.Nationality,
		PassportNo:    p.PassportNo,
		SeatNumber:    p.SeatNumber,
		ClassName:     p.FlightClass.ClassType,
		Price:         p.Price,
	}
}

func ToTransactionResponse(t models.Transaction) TransactionResponse {
	var items []TransactionItemResponse
	for _, p := range t.Items {
		items = append(items, ToTransactionItemResponse(p))
	}

	return TransactionResponse{
		Code:          t.Code,
		FlightNumber:  t.Flight.FlightNumber,
		TotalPrice:    t.TotalPrice,
		PaymentStatus: t.PaymentStatus,
		PaymentURL:    t.PaymentURL,
		PromoCode:     t.PromoCode,
		Discount:      t.Discount,
		ExpiresAt:     t.ExpiresAt,
		CreatedAt:     t.CreatedAt,
		Items:         items,
	}
}
