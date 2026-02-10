package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service *ProductService
	repo    *ProductRepository
}

func NewProductController(s *ProductService, r *ProductRepository) *ProductController {
	return &ProductController{s, r}
}

func (c *ProductController) Create(ctx *gin.Context) {
	var req CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.Create(req); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, gin.H{"message": "product created"})
}

func (c *ProductController) GetAll(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	offset := (page - 1) * limit

	data, _ := c.repo.FindAll(limit, offset)
	ctx.JSON(200, data)
}
