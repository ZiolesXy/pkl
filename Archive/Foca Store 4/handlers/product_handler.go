package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"focastore/models"
	"focastore/request"
	"focastore/response"
)

type ProductHandler struct {
	db *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

func (h *ProductHandler) List(c *gin.Context) {
	var products []models.Product
	if err := h.db.Order("id desc").Find(&products).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}
	response.Success(c, "products", gin.H{"items": products})
}

func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	var p models.Product
	if err := h.db.First(&p, uint(id)).Error; err != nil {
		response.Error(c, http.StatusNotFound, "product not found")
		return
	}
	response.Success(c, "product", p)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req request.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}
	p := models.Product{Name: req.Name, Description: req.Description, Price: req.Price, Stock: req.Stock}
	if err := h.db.Create(&p).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "create product failed")
		return
	}
	response.Created(c, "product created", p)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req request.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	var p models.Product
	if err := h.db.First(&p, uint(id)).Error; err != nil {
		response.Error(c, http.StatusNotFound, "product not found")
		return
	}

	p.Name = req.Name
	p.Description = req.Description
	p.Price = req.Price
	p.Stock = req.Stock
	if err := h.db.Save(&p).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "update product failed")
		return
	}
	response.Success(c, "product updated", p)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.db.Delete(&models.Product{}, uint(id)).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "delete product failed")
		return
	}
	response.Success(c, "product deleted", gin.H{})
}
