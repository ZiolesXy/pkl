package repository

import (
	"main/internal/model"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}