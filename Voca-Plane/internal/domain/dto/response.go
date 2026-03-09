package dto

import (
	"time"
	"voca-plane/internal/domain/models"
)

type AirlineResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	LogoURL string `json:"logo_url"`
}

type AirportResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}

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
	ID            uint                  `json:"id"`
	Airline       AirlineResponse       `json:"airline"`
	Origin        AirportResponse       `json:"origin"`
	Destination   AirportResponse       `json:"destination"`
	DepartureTime time.Time             `json:"departure_time"`
	ArrivalTime   time.Time             `json:"arrival_time"`
	FlightNumber  string                `json:"flight_number"`
	TotalSeats    int                   `json:"total_seats"`
	TotalRows     int                   `json:"total_rows"`
	TotalColumns  int                   `json:"total_columns"`
	FlightClasses []FlightClassResponse `json:"classes,omitempty"`
}

func ToAirlineResponse(a models.Airline) AirlineResponse {
	return AirlineResponse{
		ID:      a.ID,
		Name:    a.Name,
		Code:    a.Code,
		LogoURL: a.LogoURL,
	}
}

func ToAirportResponse(a models.Airport) AirportResponse {
	return AirportResponse{
		ID:   a.ID,
		Code: a.Code,
		Name: a.Name,
		City: a.City,
	}
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
	var classes []FlightClassResponse
	for _, c := range f.FlightClasses {
		classes = append(classes, ToFlightClassResponse(c))
	}

	return FlightResponse{
		ID:            f.ID,
		Airline:       ToAirlineResponse(f.Airline),
		Origin:        ToAirportResponse(f.Origin),
		Destination:   ToAirportResponse(f.Destination),
		DepartureTime: f.DepartureTime,
		ArrivalTime:   f.ArrivalTime,
		FlightNumber:  f.FlightNumber,
		TotalSeats:    f.TotalSeats,
		TotalRows:     f.TotalRows,
		TotalColumns:  f.TotalColumns,
		FlightClasses: classes,
	}
}
