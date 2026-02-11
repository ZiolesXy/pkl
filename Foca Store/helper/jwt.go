package helper

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("ADMIN_PETIR")

func GenerateToken(uid uint, dur time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": uid,
		"exp": time.Now().Add(dur).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func ParseToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}
