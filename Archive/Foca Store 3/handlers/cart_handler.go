package handlers

import (
	"foca-store/database"
	"foca-store/models"
	"foca-store/request"
	"foca-store/response"
	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	uid := c.GetUint("user_id")
	var req request.CartRequest
	c.ShouldBindJSON(&req)

	var cart models.Cart
	database.DB.FirstOrCreate(&cart, models.Cart{UserID: uid})

	item := models.CartItem{
		CartID: cart.ID,
		ProductID: req.ProductID,
		Quantity: req.Quantity,
	}
	database.DB.Create(&item)
	response.Success(c, "added to cart", item)
}
