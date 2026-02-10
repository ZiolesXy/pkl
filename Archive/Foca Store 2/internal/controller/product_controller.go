package controller

import (
	"main/internal/model"
	"main/internal/service"
	"main/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service service.ProductService
}

func NewProductController(service service.ProductService) *ProductController {
	return &ProductController{service}
}

func (ctrl *ProductController) CreateCategory(c *gin.Context) {
	var input model.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}
	if err := ctrl.service.CreateCategory(input); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Category created", nil)
}

func (ctrl *ProductController) GetAllCategories(c *gin.Context) {
	categories, err := ctrl.service.GetAllCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch categories", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Categories retrieved", categories)
}

func (ctrl *ProductController) CreateProduct(c *gin.Context) {
	var input model.ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}
	if err := ctrl.service.CreateProduct(input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create product", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Product created", nil)
}

func (ctrl *ProductController) GetAllProducts(c *gin.Context) {
	products, err := ctrl.service.GetAllProducts()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch products", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Products retrieved", products)
}