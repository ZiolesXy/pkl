package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID uint, role string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role": role,
		"exp": time.Now().Add(duration).Unix(),
		"type": "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
}

func GenerateRefreshToken(userID uint, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
		"type": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
}

func ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	return validateToken(tokenString, os.Getenv("JWT_ACCESS_SECRET"), "access")
}

func ValidateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	return validateToken(tokenString, os.Getenv("JWT_REFRESH_SECRET"), "refresh")
}

func validateToken(tokenString, secret, tokenType string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	if claims["type"] != tokenType {
		return nil, errors.New("invalid token tyype")
	}

	return claims, nil
}