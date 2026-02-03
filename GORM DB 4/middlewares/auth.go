package middlewares

import (
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
				respons.NewJsonResponse("No token", nil),
			)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error){
			return helpers.SECRET, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(
				401,
				respons.NewJsonResponse("Invalid tokrn", nil),
			)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])

		c.Next()
	}
}