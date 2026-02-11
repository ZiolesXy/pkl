package handlers

import (
	"foca-store/database"
	"foca-store/models"
	"foca-store/request"
	"foca-store/response"
	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var req request.ProductRequest
	c.ShouldBindJSON(&req)

	p := models.Product{
		Name: req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	database.DB.Create(&p)
	response.Success(c, "product created", p)
}

func GetProducts(c *gin.Context) {
	var p []models.Product
	database.DB.Find(&p)
	response.SuccessList(c, "product list", p)
}