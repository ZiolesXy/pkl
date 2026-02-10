package product

import (
	"main/internal/database"
	"main/internal/models"
)

type ProductRepository struct{}

func (r *ProductRepository) Create(p *models.Product) error {
	return database.DB.Create(p).Error
}

func (r *ProductRepository) FindAll(limit, offset int) ([]models.Product, error) {
	var products []models.Product
	err := database.DB.Preload("Category").Limit(limit).Offset(offset).Find(&products).Order("id ASC").Error
	return products, err
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var p models.Product
	err := database.DB.Preload("Category").First(&p, id).Error
	return &p, err
}

func (r *ProductRepository) Update(p *models.Product) error {
	return database.DB.Save(p).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Product{}, id).Error
}