package handlers

import (
	"net/http"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.CreateCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
			return
		}

		category := models.Category{
			Name: req.Name,
		}

		if err := db.Create(&category).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to create category")
			return
		}

		resp := response.BuildCategoryResponse(
			category.ID,
			category.Name,
			category.CreatedAt,
			category.UpdatedAt,
		)

		response.SuccessResponse(c, "category created succesfully", resp)
	}
}

func GetAllCategory(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		var category []models.Category
		if err := db.Order("id ASC").Find(&category).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch category")
			return
		}

		var responses []response.CategoryResponse
		for _, cat := range category {
			responses = append(responses, response.BuildCategoryResponse(
				cat.ID,
				cat.Name,
				cat.CreatedAt,
				cat.UpdatedAt,
			))
		}

		listResp := response.BuildCategoryListResponse(responses)
		response.SuccessListResponse(c, "category retrived successfully", listResp)
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
			response.ErrorResponse(c, http.StatusBadRequest, "invalid body request")
			return
		}

		if req.Name != nil {
			category.Name = *req.Name
		}

		if err := db.Save(&category).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to update category")
			return
		}

		resp := response.BuildCategoryResponse(
			category.ID,
			category.Name,
			category.CreatedAt,
			category.UpdatedAt,
		)

		response.SuccessResponse(c, "category updated succesfully", resp)
	}
}