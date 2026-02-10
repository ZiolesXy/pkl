package service

import (
	"errors"
	"main/internal/model"
	"main/internal/repository"
	"main/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(input model.RegisterInput) error
	Login(input model.LoginInput) (string, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(input model.RegisterInput) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     model.RoleCustomer,
	}

	return s.repo.CreateUser(&user)
}

func (s *authService) Login(input model.LoginInput) (string, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	return utils.GenerateToken(user.ID, user.Role)
}