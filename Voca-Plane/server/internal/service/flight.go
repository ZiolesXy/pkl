package service

import (
	"context"
	"voca-plane/internal/domain/dto/response"
	"voca-plane/internal/domain/models"
	"voca-plane/internal/repository"
)

type FlightService struct {
	flightRepo repository.FlightRepository
}

func NewFlightService(flightRepo repository.FlightRepository) *FlightService {
	return &FlightService{flightRepo: flightRepo}
}

func (s *FlightService) SearchFlight(ctx context.Context, origin, destination, date, classType string, page, limit int) ([]response.FlightResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	flights, total, err := s.flightRepo.Search(ctx, origin, destination, date, classType, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var flightResponses []response.FlightResponse
	for _, f := range flights {
		flightResponses = append(flightResponses, response.ToFlightResponse(f))
	}

	return flightResponses, total, nil
}

func (s *FlightService) GetFlightByID(ctx context.Context, id uint) (*response.FlightResponse, error) {
	flight, err := s.flightRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := response.ToFlightResponse(*flight)
	return &response, nil
}

func (s *FlightService) GetAllFlights(ctx context.Context, page, limit int) ([]response.FlightResponse, int64, error) {
	if page <  1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	flights, total, err := s.flightRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var flightResponses []response.FlightResponse
	for _, f := range flights {
		flightResponses = append(flightResponses, response.ToFlightResponse(f))
	}

	return flightResponses, total, nil
}

func (s *FlightService) CreateFlight(ctx context.Context, flight *models.Flight) error {
	return s.flightRepo.Create(ctx, nil, flight)
}

func (s *FlightService) UpdateFlight(ctx context.Context, flight *models.Flight) error {
	return s.flightRepo.Update(ctx, nil, flight)
}

func (s *FlightService) DeleteFlight(ctx context.Context, id uint) error {
	return s.flightRepo.Delete(ctx, nil, id)
}