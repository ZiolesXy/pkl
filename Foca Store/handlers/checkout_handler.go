package handlers

import (
	"gorm.io/gorm"
	"foca-store/database"
	"foca-store/models"
	"foca-store/request"
	"foca-store/response"
	"github.com/gin-gonic/gin"
)

func CreateCheckout(c *gin.Context) {
	uid := c.GetUint("user_id")

	var cart models.Cart
	database.DB.Preload("Items.Product").Where("user_id = ?", uid).First(&cart)
	if len(cart.Items) == 0 {
		response.Error(c, 400, "cart empty")
		return
	}

	total := 0
	for _, i := range cart.Items {
		total += i.Product.Price * i.Quantity
	}

	tx := database.DB.Begin()

	co := models.Checkout{UserID: uid, Status: "pending", Total: total}
	tx.Create(&co)

	for _, i := range cart.Items {
		tx.Create(&models.CheckoutItem{
			CheckoutID: co.ID,
			ProductID: i.ProductID,
			Quantity: i.Quantity,
			Price: i.Product.Price,
		})
	}

	tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	tx.Commit()

	response.Success(c, "checkout created", co)
}

func UpdateCheckoutStatus(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateCheckoutStatusRequest
	c.ShouldBindJSON(&req)

	var co models.Checkout
	database.DB.Preload("Items").First(&co, id)

	tx := database.DB.Begin()
	if req.Status == "success" {
		for _, i := range co.Items {
			tx.Model(&models.Product{}).
				Where("id = ?", i.ProductID).
				Update("stock", gorm.Expr("stock - ?", i.Quantity))
		}
	}
	tx.Model(&co).Update("status", req.Status)
	tx.Commit()

	response.Success(c, "checkout updated", co)
}
