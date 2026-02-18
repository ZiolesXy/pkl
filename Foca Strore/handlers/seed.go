package handlers

import (
	"net/http"

	"voca-store/response"
	"voca-store/seeders"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ================= ROLE =================

func SeedRoleHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := seeders.SeedRoles(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed roles")
			return
		}

		response.SuccessResponse(c, "Roles seeded successfully", nil)
	}
}

// ================= ADMIN =================

func SeedAdminHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := seeders.SeedAdmin(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed admin")
			return
		}

		response.SuccessResponse(c, "Admin seeded successfully", nil)
	}
}

// ================= USERS =================

func SeedUsersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := seeders.SeedUsers(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed users")
			return
		}

		response.SuccessResponse(c, "Users seeded successfully", nil)
	}
}

// ================= PRODUCTS =================

func SeedProductsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := seeders.SeedProducts(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed products")
			return
		}

		response.SuccessResponse(c, "Products seeded successfully", nil)
	}
}

// ================= ALL =================

func SeedAllHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := seeders.SeedRoles(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed roles")
			return
		}

		if err := seeders.SeedAdmin(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed admin")
			return
		}

		if err := seeders.SeedUsers(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed users")
			return
		}

		if err := seeders.SeedProducts(db); err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to seed products")
			return
		}

		response.SuccessResponse(c, "All data seeded successfully", nil)
	}
}
