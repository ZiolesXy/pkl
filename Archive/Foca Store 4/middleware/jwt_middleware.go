package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"focastore/helper"
	"focastore/response"
)

const (
	CtxUserIDKey = "user_id"
	CtxRoleKey   = "role"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			response.Error(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(c, http.StatusUnauthorized, "invalid authorization header")
			c.Abort()
			return
		}

		claims, err := helper.ParseToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}
		if claims.Type != helper.TokenAccess {
			response.Error(c, http.StatusUnauthorized, "invalid token type")
			c.Abort()
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxRoleKey, claims.Role)
		c.Next()
	}
}
