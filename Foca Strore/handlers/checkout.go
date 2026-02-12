package handlers

import (
	"errors"
	"net/http"
	// "time"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Checkout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
			return
		}

		// Get user's cart with items
		var cart models.Cart
		if err := db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.ErrorResponse(c, http.StatusNotFound, "Cart not found")
				return
			}
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch cart")
			return
		}

		if len(cart.Items) == 0 {
			response.ErrorResponse(c, http.StatusBadRequest, "Cart is empty")
			return
		}

		// Calculate total price and check stock
		totalPrice := 0.0
		for _, item := range cart.Items {
			if item.Product == nil {
				response.ErrorResponse(c, http.StatusBadRequest, "Product not found in cart item")
				return
			}
			if item.Product.Stock < item.Quantity {
				response.ErrorResponse(c, http.StatusBadRequest, 
					"Insufficient stock for product: " + item.Product.Name)
				return
			}
			totalPrice += item.Product.Price * float64(item.Quantity)
		}

		// Create checkout record
		checkout := models.Checkout{
			UserID:     userID.(uint),
			TotalPrice: totalPrice,
			Status:     "pending",
		}

		if err := db.Create(&checkout).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create checkout")
			return
		}

		// Clear cart items
		if err := db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			// Rollback checkout if failed to clear cart
			db.Delete(&checkout)
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to clear cart")
			return
		}

		// Build checkout response
		checkoutResp := response.BuildCheckoutResponse(
			checkout.ID,
			checkout.UserID,
			"",
			checkout.TotalPrice,
			checkout.Status,
			checkout.CreatedAt,
			checkout.UpdatedAt,
		)

		response.SuccessResponse(c, "Checkout successful", checkoutResp)
	}
}

func UpdateCheckoutStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req request.UpdateCheckoutStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		checkoutID := c.Param("id")
		var checkout models.Checkout
		if err := db.First(&checkout, checkoutID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.ErrorResponse(c, http.StatusNotFound, "Checkout not found")
			} else {
				response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch checkout")
			}
			return
		}

		if err := db.Model(&checkout).Update("status", req.Status).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update checkout status")
			return
		}

		// Reload checkout
		if err := db.First(&checkout, checkout.ID).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to reload checkout")
			return
		}

		// Build checkout response
		checkoutResp := response.BuildCheckoutResponse(
			checkout.ID,
			checkout.UserID,
			"",
			checkout.TotalPrice,
			checkout.Status,
			checkout.CreatedAt,
			checkout.UpdatedAt,
		)

		response.SuccessResponse(c, "Checkout status updated successfully", checkoutResp)
	}
}