package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"focastore/response"
)

func RequireRole(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleAny, ok := c.Get(CtxRoleKey)
		if !ok {
			response.Error(c, http.StatusForbidden, "forbidden")
			c.Abort()
			return
		}

		role, _ := roleAny.(string)
		if role != required {
			response.Error(c, http.StatusForbidden, "forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
