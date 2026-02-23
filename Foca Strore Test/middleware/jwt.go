package middleware

import (
	"net/http"
	"voca-store/helper"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JWTAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorResponse(c, http.StatusUnauthorized, "Missing authorization header")
			c.Abort()
			return
		}

		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		claims, err := helper.ValidateAccessToken(tokenString)
		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
