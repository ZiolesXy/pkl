package middleware

import (
	"foca-store/response"
	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {

		role := c.GetString("role")

		if role != "Admin" {
			response.Error(c, 403, "forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}

