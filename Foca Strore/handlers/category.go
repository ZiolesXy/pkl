package handlers

import (
	"errors"
	"net/http"
	"voca-store/helper"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryWithCount struct {
	models.Category
	ProductCount int64
}

func CreateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req request.CreateCategoryRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid request")
			return
		}

		slug, err := helper.GenerateUniqueCategorySlug(db, req.Name)

		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed generate slug")
			return
		}

		category := models.Category{
			Name: req.Name,
			Slug: slug,
		}

		if err := db.Create(&category).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to create category")
			return
		}

		res := response.BuildCategoryResponse(category, 0)

		response.SuccessResponse(c, "category created", res)
	}
}

func GetAllCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rows []CategoryWithCount
		err := db.
			Table("categories").
			Select(`
				categories.*,
				COUNT(products.id) as product_count
			`).
			Joins(`
				LEFT JOIN products
				ON products.category_id = categories.id
			`).
			Group("categories.id").
			Order("categories.id ASC").
			Scan(&rows).Error

		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to get category")
			return
		}

		var categories []models.Category
		countMap := map[uint]int64{}

		for _, r := range rows {
			categories = append(categories, r.Category)
			countMap[r.ID] = r.ProductCount
		}

		res := response.BuildCategoryListResponse(
			categories,
			countMap,
		)

		response.SuccessListResponse(c, "categories retrivied", res)
	}
}

func GetCategoryBySlug(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")

		var category models.Category

		err := db.
			Where("slug = ?", slug).
			First(&category).
			Error

		if err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "category not found")
			return
		}

		var products []models.Product

		err = db.
			Preload("Category").
			Where("category_id = ?", category.ID).
			Order("id ASC").
			Find(&products).
			Error
		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch products")
			return
		}

		res := response.BuildCategoryDetailResponse(
			category,
			products,
		)

		response.SuccessResponse(c, "category retrieved", res)
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
			slug, err := helper.GenerateUniqueCategorySlug(db, *req.Name)

			if err != nil {
				response.ErrorResponse(c, http.StatusInternalServerError, "failed generate slug")
				return
			}
			updates["name"] = *req.Name
			updates["slug"] = slug
		}

		if err := db.Model(&category).Updates(updates).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to update category")
			return
		}

		db.First(&category, category.ID)

		var row CategoryWithCount

		db.Table("categories").
		Select(`
			categories.*,
			COUNT(products.id) as product_count
		`).
		Joins(`
			LEFT JOIN products
			ON products.category_id = categories.id
		`).
		Where("categories.id = ?", category.ID).
		Group("categories.id").
		Scan(&row)

		res := response.BuildCategoryResponse(
			row.Category,
			row.ProductCount,
		)
		response.SuccessResponse(c, "category updated", res)
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
