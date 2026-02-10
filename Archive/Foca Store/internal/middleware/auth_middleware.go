package middleware

import (
	"net/http"
	"strings"
	"main/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		claims, err := utils.ValidateAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		role := claims["role"].(string)

		if len(requiredRoles) > 0 {
			allowed := false
			for _, r := range requiredRoles {
				if r == role {
					allowed = true
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
		}

		c.Set("user_id", claims["user_id"])
		c.Set("role", role)

		c.Next()
	}
}