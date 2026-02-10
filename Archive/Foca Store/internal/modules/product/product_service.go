package product

import "main/internal/models"

type ProductService struct {
	repo *ProductRepository
}

func NewProductService(r *ProductRepository) *ProductService {
	return &ProductService{r}
}

func (s *ProductService) Create(req CreateProductRequest) error {
	p := models.Product{
		Name: req.Name,
		Price: req.Price,
		Stock: req.Stock,
		CategoryID: req.CategoryID,
	}
	return s.repo.Create(&p)
}