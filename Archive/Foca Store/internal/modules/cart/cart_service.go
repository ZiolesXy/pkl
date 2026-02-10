package cart

import (
	"errors"

	"main/internal/database"
	"main/internal/models"
)

type CartService struct {
	repo *CartRepository
}

func NewCartService(r *CartRepository) *CartService {
	return &CartService{r}
}

func (s *CartService) AddItem(userID, productID uint, qty int) error {
	cart, _ := s.repo.GetByUser(userID)

	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		return errors.New("product not found")
	}

	item := models.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  qty,
	}
	return database.DB.Create(&item).Error
}
