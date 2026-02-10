package utils

import (
	"main/pkg/config"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint, role string) (string, error) {
	secret := config.GetEnv("JWT_SECRET", "secret")
	expiryHours, _ := strconv.Atoi(config.GetEnv("JWT_EXPIRATION_HOUR", "24"))

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}