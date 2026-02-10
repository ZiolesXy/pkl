package cart

import (
	"main/internal/database"
	"main/internal/models"
)

type CartRepository struct{}

func (r *CartRepository) GetByUser(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := database.DB.Preload("Items.Product").
		Where("user_id = ?", userID).
		FirstOrCreate(&cart, models.Cart{UserID: userID}).Error
	return &cart, err
}
