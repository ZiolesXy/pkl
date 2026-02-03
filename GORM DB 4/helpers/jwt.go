package helpers

import (
	"main/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET = []byte("RAHASIA_SUPER")

func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role": user.RoleID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET)
}