package auth

import (
	"main/internal/database"
	"main/internal/models"
)

type AuthRepository struct{}

func (r *AuthRepository) CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func (r *AuthRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("Role").Where("email = ?", email).Error
	return &user, err
}