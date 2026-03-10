package response

import (
	"time"
	"voca-plane/internal/domain/models"
)

type FlightSeatResponse struct {
	ID          uint   `json:"id"`
	SeatNumber  string `json:"seat_number"`
	IsAvailable bool   `json:"is_available"`
}

type FlightClassResponse struct {
	ID         uint                 `json:"id"`
	ClassType  string               `json:"class_type"`
	Price      float64              `json:"price"`
	TotalSeats int                  `json:"total_seats"`
	Seats      []FlightSeatResponse `json:"seats,omitempty"`
}

type FlightResponse struct {
	ID             uint                  `json:"id"`
	Airline        AirlineResponse       `json:"airline"`
	Origin         AirportResponse       `json:"origin"`
	Destination    AirportResponse       `json:"destination"`
	DepartureTime  time.Time              `json:"departure_time"`
	ArrivalTime    time.Time             `json:"arrival_time"`
	FlightNumber   string                `json:"flight_number"`
	TotalSeats     int                   `json:"total_seats"`
	AvailableSeats int                   `json:"available_seats"`
	TotalRows      int                   `json:"total_rows"`
	TotalColumns   int                   `json:"total_columns"`
	FlightClasses  []FlightClassResponse `json:"classes,omitempty"`
}

type TransactionResponse struct {
	Code          string    `json:"code"`
	FlightNumber  string    `json:"flight_number"`
	ClassType     string    `json:"class_type"`
	TotalPrice    float64   `json:"total_price"`
	PaymentStatus string    `json:"payment_status"`
	PaymentURL    string    `json:"payment_url"`
	ExpiresAt     time.Time `json:"expires_at"`
}

func ToFlightSeatResponse(s models.FlightSeat) FlightSeatResponse {
	return FlightSeatResponse{
		ID:          s.ID,
		SeatNumber:  s.SeatNumber,
		IsAvailable: s.IsAvailable,
	}
}

func ToFlightClassResponse(fc models.FlightClass) FlightClassResponse {
	var seats []FlightSeatResponse
	for _, s := range fc.Seats {
		seats = append(seats, ToFlightSeatResponse(s))
	}

	return FlightClassResponse{
		ID:         fc.ID,
		ClassType:  fc.ClassType,
		Price:      fc.Price,
		TotalSeats: len(fc.Seats),
		Seats:      seats,
	}
}

func ToFlightResponse(f models.Flight) FlightResponse {
	available := 0
	for _, class := range f.FlightClasses {
		for _, seat := range class.Seats {
			if seat.IsAvailable {
				available++
			}
		}
	}

	var classes []FlightClassResponse
	for _, c := range f.FlightClasses {
		classes = append(classes, ToFlightClassResponse(c))
	}

	return FlightResponse{
		ID:             f.ID,
		Airline:        ToAirlineResponse(f.Airline),
		Origin:         ToAirportResponse(f.Origin),
		Destination:    ToAirportResponse(f.Destination),
		DepartureTime:  f.DepartureTime,
		ArrivalTime:    f.ArrivalTime,
		FlightNumber:   f.FlightNumber,
		TotalSeats:     f.TotalSeats,
		AvailableSeats: available,
		TotalRows:      f.TotalRows,
		TotalColumns:   f.TotalColumns,
		FlightClasses:  classes,
	}
}

func ToTransactionResponse(t models.Transaction) TransactionResponse {
	return TransactionResponse{
		Code:          t.Code,
		FlightNumber:  t.Flight.FlightNumber,
		ClassType:     t.FlightClass.ClassType,
		TotalPrice:    t.TotalPrice,
		PaymentStatus: t.PaymentStatus,
		PaymentURL:    t.PaymentURL,
		ExpiresAt:     t.ExpiresAt,
	}
}