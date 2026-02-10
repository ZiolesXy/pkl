package service

import (
	"errors"
	"main/internal/model"
	"main/internal/repository"
)

type ProductService interface {
	CreateCategory(input model.CategoryInput) error
	GetAllCategories() ([]model.Category, error)
	CreateProduct(input model.ProductInput) error
	GetAllProducts() ([]model.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo}
}

func (s *productService) CreateCategory(input model.CategoryInput) error {
	category := model.Category{Name: input.Name}
	return s.repo.CreateCategory(&category)
}

func (s *productService) GetAllCategories() ([]model.Category, error) {
	return s.repo.FindAllCategories()
}

func (s *productService) CreateProduct(input model.ProductInput) error {
	// Validasi apakah kategori ada
	_, err := s.repo.FindCategoryByID(input.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	product := model.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		CategoryID:  input.CategoryID,
	}
	return s.repo.CreateProduct(&product)
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	return s.repo.FindAllProducts()
}