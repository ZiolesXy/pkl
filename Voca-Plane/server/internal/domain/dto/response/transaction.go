package response

import (
	"time"
	"voca-plane/internal/domain/models"
)

type TransactionPassengerResponse struct {
	FullName    string `json:"full_name"`
	Nationality string `json:"nationality"`
	PassportNo  string `json:"passport_no"`
	SeatNumber  string `json:"seat_number"`
}

type TransactionResponse struct {
	Code          string                         `json:"code"`
	FlightNumber  string                         `json:"flight_number"`
	ClassName     string                         `json:"class_name"`
	TotalPrice    float64                        `json:"total_price"`
	PaymentStatus string                         `json:"payment_status"`
	PaymentURL    string                         `json:"payment_url"`
	PromoCode     *string                        `json:"promo_code,omitempty"`
	Discount      float64                        `json:"discount"`
	ExpiresAt     time.Time                      `json:"expires_at"`
	CreatedAt     time.Time                      `json:"created_at"`
	Passengers    []TransactionPassengerResponse `json:"passengers,omitempty"`
}

func ToTransactionPassengerResponse(p models.TransactionPassenger) TransactionPassengerResponse {
	return TransactionPassengerResponse{
		FullName:    p.FullName,
		Nationality: p.Nationality,
		PassportNo:  p.PassportNo,
		SeatNumber:  p.SeatNumber,
	}
}

func ToTransactionResponse(t models.Transaction) TransactionResponse {
	var passengers []TransactionPassengerResponse
	for _, p := range t.Passengers {
		passengers = append(passengers, ToTransactionPassengerResponse(p))
	}

	return TransactionResponse{
		Code:          t.Code,
		FlightNumber:  t.Flight.FlightNumber,
		ClassName:     t.FlightClass.ClassType,
		TotalPrice:    t.TotalPrice,
		PaymentStatus: t.PaymentStatus,
		PaymentURL:    t.PaymentURL,
		PromoCode:     t.PromoCode,
		Discount:      t.Discount,
		ExpiresAt:     t.ExpiresAt,
		CreatedAt:     t.CreatedAt,
		Passengers:    passengers,
	}
}
