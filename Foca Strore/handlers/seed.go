package handlers

import (
	"net/http"
	// "voca-store/database"
	// "voca-store/models"
	"voca-store/response"
	"voca-store/seeders"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SeedHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Seed roles
		if err := seeders.SeedRoles(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed roles")
			return
		}

		// Seed admin
		if err := seeders.SeedAdmin(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed admin")
			return
		}

		// Seed users
		if err := seeders.SeedUsers(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed users")
			return
		}

		// Seed products
		if err := seeders.SeedProducts(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed products")
			return
		}

		response.SuccessResponse(c, "Database seeded successfully", nil)
	}
}