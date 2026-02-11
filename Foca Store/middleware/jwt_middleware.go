package middleware

import (
	"foca-store/helper"
	"foca-store/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)
		token, err := helper.ParseToken(t)
		if err != nil || !token.Valid {
			response.Error(c, 401, "unauthorized")
			c.Abort()
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Next()
	}
}
