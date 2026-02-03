package middlewares

import (
	"main/respons"

	"github.com/gin-gonic/gin"
)

func OnlyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetFloat64("role")
		if role != 1 {
			c.AbortWithStatusJSON(
				403,
				respons.NewJsonResponse("Forbidden", nil),
			)
			return
		}
		c.Next()
	}
}