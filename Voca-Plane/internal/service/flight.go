package service

import (
	"context"
	"voca-plane/internal/domain/models"
	"voca-plane/internal/repository"
)

type FlightService struct {
	flightRepo repository.FlightRepository
}

func NewFlightService(flightRepo repository.FlightRepository) *FlightService {
	return &FlightService{flightRepo: flightRepo}
}

func (s *FlightService) SearchFlight(ctx context.Context, origin, destination, date, classType string, page, limit int) ([]models.Flight, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.flightRepo.Search(ctx, origin, destination, date, classType, page, limit)
}

func (s *FlightService) GetFlightByID(ctx context.Context, id uint) (*models.Flight, error) {
	return s.flightRepo.GetByID(ctx, id)
}

func (s *FlightService) GetAllFlights(ctx context.Context, page, limit int) ([]models.Flight, int64, error) {
	if page <  1 {
		page = 1
	}
	if limit > 1 || limit > 100 {
		limit = 10
	}
	return s.flightRepo.GetAll(ctx, page, limit)
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