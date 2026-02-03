package middlewares

import (
	"errors"
	"main/helpers"
	"main/respons"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.AbortWithStatusJSON(
				401,
				respons.NewJsonResponse("No token provided", nil),
			)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// ✅ pastikan algoritma HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return helpers.ACCESS_SECRET, nil
		})

		// 🔥 BEDAKAN ERROR JWT
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(
					401,
					respons.NewJsonResponse("Token expired", nil),
				)
				return
			}

			c.AbortWithStatusJSON(
				401,
				respons.NewJsonResponse("Invalid token", nil),
			)
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(
				401,
				respons.NewJsonResponse("Invalid token", nil),
			)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(
				401,
				respons.NewJsonResponse("Invalid token claims", nil),
			)
			return
		}

		// ✅ pastikan token access
		if claims["type"] != "access" {
			c.AbortWithStatusJSON(
				401,
				respons.NewJsonResponse("Invalid token type", nil),
			)
			return
		}

		// ✅ simpan data ke context
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])

		c.Next()
	}
}
