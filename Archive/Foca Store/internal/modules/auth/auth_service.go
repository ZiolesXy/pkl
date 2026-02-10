package auth

import (
	"errors"
	"main/internal/models"
	"time"

	"main/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *AuthRepository
}

func NewAuthService(repo *AuthRepository) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) Register(req RegisterRequest) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.User{
		Name: req.Name,
		Email: req.Email,
		Password: string(hash),
		RoleID: 2,
	}
	return s.repo.CreateUser(&user)
}

func (s *AuthService) Login(req LoginRequest) (string, string, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", "", errors.New("invalid credential")
	}

	access, _ := utils.GenerateAccessToken(user.ID, user.Role.Name, time.Minute*15)
	refresh, _ := utils.GenerateRefreshToken(user.ID, time.Hour*24*7)	

	return access, refresh, nil
}