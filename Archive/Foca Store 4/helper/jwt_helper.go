package helper

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	TokenAccess  TokenType = "access"
	TokenRefresh TokenType = "refresh"
)

type CustomClaims struct {
	UserID uint      `json:"user_id"`
	Role   string    `json:"role"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

func secret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return nil, errors.New("JWT_SECRET is required")
	}
	return []byte(s), nil
}

func GenerateToken(userID uint, role string, tokenType TokenType, ttl time.Duration) (string, error) {
	sec, err := secret()
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(sec)
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	sec, err := secret()
	if err != nil {
		return nil, err
	}

	tok, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return sec, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(*CustomClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
