package service

import (
	"context"
	"errors"
	"voca-plane/internal/domain/models"
	"voca-plane/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	user.Password = ""
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint, name, email, password string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if name != "" {
		user.Name = name
	}

	if email != "" && email != user.Email {
		_, err := s.userRepo.FindByEmail(ctx, email)
		if err == nil {
			return errors.New("email already registered")
		}
		user.Email = email
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	return s.userRepo.Update(ctx, user)
}