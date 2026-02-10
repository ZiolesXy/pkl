package repository

import (
	"main/internal/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	// Category
	CreateCategory(category *model.Category) error
	FindAllCategories() ([]model.Category, error)
	FindCategoryByID(id uint) (*model.Category, error)
	
	// Product
	CreateProduct(product *model.Product) error
	FindAllProducts() ([]model.Product, error)
	FindProductByID(id uint) (*model.Product, error)
	UpdateProduct(product *model.Product) error
	DeleteProduct(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) CreateCategory(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *productRepository) FindAllCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *productRepository) FindCategoryByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	return &category, err
}

func (r *productRepository) CreateProduct(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindAllProducts() ([]model.Product, error) {
	var products []model.Product
	// Preload untuk mengambil data kategori terkait (JOIN)
	err := r.db.Preload("Category").Find(&products).Error
	return products, err
}

func (r *productRepository) FindProductByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").First(&product, id).Error
	return &product, err
}

func (r *productRepository) UpdateProduct(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) DeleteProduct(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}