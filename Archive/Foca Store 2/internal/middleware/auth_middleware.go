package middleware

import (
	"fmt"
	"main/pkg/config"
	"main/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware memvalidasi JWT Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is required", nil)
			c.Abort()
			return
		}

		// 2. Cek format Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization format", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Parse dan Validasi Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma signing sesuai
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.GetEnv("JWT_SECRET", "secret")), nil
		})

		if err != nil || !token.Valid {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", err.Error())
			c.Abort()
			return
		}

		// 4. Ambil claims (user_id dan role)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims", nil)
			c.Abort()
			return
		}

		// 5. Simpan data ke context Gin agar bisa dipakai di Controller
		c.Set("userID", uint(claims["user_id"].(float64)))
		c.Set("role", claims["role"].(string))

		c.Next()
	}
}

// RoleMiddleware membatasi akses berdasarkan role tertentu
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "Role not found in context", nil)
			c.Abort()
			return
		}

		isAuthorized := false
		for _, role := range roles {
			if role == userRole {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			utils.ErrorResponse(c, http.StatusForbidden, "You don't have permission to access this resource", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}