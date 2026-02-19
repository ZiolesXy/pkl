package handlers

import (
	"errors"
	"net/http"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCategory(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		var req request.CreateCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid request")
			return
		}

		category := models.Category{
			Name: req.Name,
		}

		if err := db.Create(&category).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to create category")
			return
		}

		res := response.BuildCategoryResponse(category)
		response.SuccessResponse(c, "category created", res)
	}
}

func GetAllCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories := []models.Category{}
		if err := db.Find(&categories).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to get category")
			return
		}

		res := response.BuildCategoryListResponse(categories)
		response.SuccessListResponse(c, "categories retrieved", res)
	}
}

func UpdateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var category models.Category

		if err := db.First(&category, id).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "category not found")
			return
		}

		var req request.UpdateCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid request")
			return
		}

		updates := make(map[string]interface{})

		if req.Name != nil {
			updates["name"] = *req.Name
		}

		if err := db.Model(&category).Updates(updates).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to update category")
			return 
		}

		db.First(&category, category.ID)

		res := response.BuildCategoryResponse(category)
		response.SuccessResponse(c, "category created", res)
	}
}

func DeleteCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var category models.Category

		if err := db.First(&category, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.ErrorResponse(c, http.StatusNotFound, "category not found")
			} else {
				response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch category")
			}
			return
		}

		if err := db.Delete(&category).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to delete category")
			return
		}

		response.SuccessResponse(c, "category deleted", nil)
	}
}