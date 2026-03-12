package repository

import (
	"context"
	"fmt"
	"voca-plane/internal/domain/models"

	"gorm.io/gorm"
)

type airportRepository struct {
	db *gorm.DB
}

func NewAirportRepository(db *gorm.DB) AirportRepository {
	return &airportRepository{db: db}
}

func (r *airportRepository) GetAll(ctx context.Context, page, limit int, sortBy, order string) ([]models.Airport, int64, error) {
	var airports []models.Airport
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Airport{})
	query.Count(&total)

	// Airports Whitelist
	allowedColumns := map[string]bool{
		"id":   true,
		"name": true,
		"code": true,
		"city": true,
	}

	if sortBy != "" && allowedColumns[sortBy] {
		if order != "asc" && order != "desc" {
			order = "asc"
		}
		query = query.Order(fmt.Sprintf("%s %s", sortBy, order))
	} else {
		query = query.Order("id ASC")
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&airports).Error
	return airports, total, err
}

func (r *airportRepository) GetByID(ctx context.Context, id uint) (*models.Airport, error) {
	var airport models.Airport
	err := r.db.WithContext(ctx).First(&airport, id).Error
	return &airport, err
}

func (r *airportRepository) Create(ctx context.Context, tx *gorm.DB, airport *models.Airport) error {
	if err := tx.WithContext(ctx).Create(airport).Error; err != nil {
		return err
	}
	return tx.WithContext(ctx).First(airport, airport.ID).Error
}

func (r *airportRepository) Update(ctx context.Context, tx *gorm.DB, airport *models.Airport) error {
	if err := tx.WithContext(ctx).Save(airport).Error; err != nil {
		return err
	}
	return tx.WithContext(ctx).First(airport, airport.ID).Error
}

func (r *airportRepository) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	return tx.WithContext(ctx).Delete(&models.Airport{}, id).Error
}