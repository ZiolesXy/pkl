package middleware

import (
	"foca-store/database"
	"foca-store/models"
	"foca-store/response"
	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetUint("user_id")
		var u models.User
		database.DB.First(&u, uid)
		if u.Role != "ADMIN" {
			response.Error(c, 403, "admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}
