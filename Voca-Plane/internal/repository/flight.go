package repository

import (
	"context"
	"time"
	"voca-plane/internal/domain/models"

	"gorm.io/gorm"
)

type flightRepository struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) FlightRepository {
	return &flightRepository{db: db}
}

func (r *flightRepository) Search(ctx context.Context, origin, destination, date, classType string, page, limit int) ([]models.Flight, int64, error) {
	var flights []models.Flight
	var total int64

	layout := "2006-01-02"
	parseDate, err := time.Parse(layout, date)
	if err != nil {
		return nil, 0, err
	}
	startOfDay := time.Date(parseDate.Year(), parseDate.Month(), parseDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24*time.Hour - time.Second)

	query := r.db.WithContext(ctx).Model(&models.Flight{}).
		Joins("JOIN airports AS origin ON flights.origin_id = origin.id").
		Joins("JOIN airports AS dest ON flights.destination_id = dest.id").
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses.Seats").
		Where("origin.code = ? AND dest.code = ?", origin, destination).
		Where("flights.departure_time BETWEEN ? AND ?", startOfDay, endOfDay)
	
	if classType != "" {
		query = query.Joins("JOIN flight_classes ON flight_classes.flight_id = flights.id").
			Where("flight_classes.class_type = ?", classType)
	}

	query.Count(&total)
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&flights).Error
	return flights, total, err
}

func (r *flightRepository) GetByID(ctx context.Context, id uint) (*models.Flight, error) {
	var flight models.Flight
	err := r.db.WithContext(ctx).
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses.Seats").
		First(&flight, id).Error
	return &flight, err
}

func (r *flightRepository) GetClassByID(ctx context.Context, id uint) (*models.FlightClass, error) {
	var class models.FlightClass
	err := r.db.WithContext(ctx).Preload("Seats").First(&class, id).Error
	return &class, err
}

func (r *flightRepository) GetSeat(ctx context.Context, classID uint, seatNumber string) (*models.FlightSeat, error) {
	var seat models.FlightSeat
	err := r.db.WithContext(ctx).Where("flight_class_id = ? AND seat_number = ?", classID, seatNumber).First(&seat).Error
	return &seat, err
}

func (r *flightRepository) UpdateSeatAvailability(ctx context.Context, tx *gorm.DB, seatID uint, available bool) error {
	return tx.WithContext(ctx).Model(&models.FlightSeat{}).Where("id = ?", seatID).Update("is_available", available).Error
}

func (r *flightRepository) GetAll(ctx context.Context, page, limit int) ([]models.Flight, int64, error) {
	var flights []models.Flight
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Flight{}).
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses").
		Preload("FlightClasses.Seats")
	
	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("departure_time ASC").Find(&flights).Error
	return flights, total, err
}

func (r *flightRepository) Create(ctx context.Context, tx *gorm.DB, flight *models.Flight) error {
	if err := tx.WithContext(ctx).Create(flight).Error; err != nil {
		return err
	}

	return tx.WithContext(ctx).
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses.Seats").
		First(flight, flight.ID).Error
}

func (r *flightRepository) Update(ctx context.Context, tx *gorm.DB, flight *models.Flight) error {
	if err := tx.WithContext(ctx).Save(flight).Error; err != nil {
		return err
	}

	return tx.WithContext(ctx).
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses.Seats").
		First(flight, flight.ID).Error
}

func (r *flightRepository) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	return tx.WithContext(ctx).Delete(&models.Flight{}, id).Error
}