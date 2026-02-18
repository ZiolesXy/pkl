package handlers

import (
	"net/http"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Checkout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDRaw, exist := c.Get("user_id")
		if !exist {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID := userIDRaw.(uint)

		var req request.CheckoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		tx := db.Begin()

		var cart models.Cart
		if err := tx.Where("user_id = ?", userID).
			First(&cart).Error; err != nil {

			tx.Rollback()
			response.ErrorResponse(c, http.StatusNotFound, "cart not found")
			return
		}

		var items []models.CartItem
		if err := tx.Preload("Product").
			Where("cart_id = ? AND id IN ?", cart.ID, req.CartItemIDs).
			Find(&items).Error; err != nil {

			tx.Rollback()
			response.ErrorResponse(c, http.StatusInternalServerError, "failed fetch items")
			return
		}

		if len(items) != len(req.CartItemIDs) {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusBadRequest, "invalid cart items")
			return
		}

		totalPrice := 0.0

		for _, item := range items {

			if item.Product.Stock < item.Quantity {
				tx.Rollback()
				response.ErrorResponse(c, http.StatusBadRequest, "insufficient stock for "+item.Product.Name)
				return
			}

			if err := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {

				tx.Rollback()
				response.ErrorResponse(c, http.StatusInternalServerError, "failed update stock")
				return
			}

			totalPrice += item.Product.Price * float64(item.Quantity)
		}

		checkout := models.Checkout{
			UserID:     userID,
			TotalPrice: totalPrice,
			Status:     "pending",
		}

		if err := tx.Create(&checkout).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusInternalServerError, "failed create checkout")
			return
		}

		for _, item := range items {
			checkoutItem := models.CheckoutItem{
				CheckoutID: checkout.ID,
				ProductID:  item.ProductID,
				Quantity:   item.Quantity,
				Price:      item.Product.Price,
			}

			if err := tx.Create(&checkoutItem).Error; err != nil {
				tx.Rollback()
				response.ErrorResponse(c, http.StatusInternalServerError, "failed create checkout item")
				return
			}
		}

		tx.Where("id IN ?", req.CartItemIDs).
			Delete(&models.CartItem{})

		if err := tx.Commit().Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "transaction failed")
			return
		}

		// 🔥 Reload with relations
		var result models.Checkout
		db.Preload("User").
			Preload("Items").
			Preload("Items.Product").
			First(&result, checkout.ID)

		res := response.BuildCheckoutDetailResponse(result)
		response.SuccessResponse(c, "Checkout created", res)
	}
}

func ApproveCheckout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		var checkout models.Checkout
		if err := db.Preload("User").
			Preload("Items").
			Preload("Items.Product").
			First(&checkout, id).Error; err != nil {

			response.ErrorResponse(c, http.StatusNotFound, "checkout not found")
			return
		}

		if checkout.Status != "pending" {
			if checkout.Status == "rejected" {
				response.ErrorResponse(c, http.StatusBadRequest, "status already rejected cannot change")
				return
			}
			response.ErrorResponse(c, http.StatusBadRequest, "invalid status")
			return
		}

		checkout.Status = "approved"

		if err := db.Save(&checkout).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed update status")
			return
		}

		res := response.BuildCheckoutDetailResponse(checkout)
		response.SuccessResponse(c, "Checkout approved", res)
	}
}

func RejectCheckout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		tx := db.Begin()

		var checkout models.Checkout
		if err := tx.Preload("User").
			Preload("Items").
			Preload("Items.Product").
			First(&checkout, id).Error; err != nil {

			tx.Rollback()
			response.ErrorResponse(c, http.StatusNotFound, "checkout not found")
			return
		}

		if checkout.Status != "pending" {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusBadRequest, "invalid status")
			return
		}

		for _, item := range checkout.Items {
			if err := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {

				tx.Rollback()
				response.ErrorResponse(c, http.StatusInternalServerError, "failed restore stock")
				return
			}
		}

		checkout.Status = "rejected"

		if err := tx.Save(&checkout).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusInternalServerError, "failed update status")
			return
		}

		if err := tx.Commit().Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "transaction failed")
			return
		}

		res := response.BuildCheckoutDetailResponse(checkout)
		response.SuccessResponse(c, "Checkout rejected", res)
	}
}

func GetCheckout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// status := c.Query("status")

		var checkouts []models.Checkout

		query := db.
			Preload("User").
			Preload("Items").
			Preload("Items.Product").
			Order("created_at DESC")

		if err := query.Find(&checkouts).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch checkout")
			return
		}

		res := response.BuildCheckOutListResponse(checkouts)
		response.SuccessListResponse(c, "checkout list fetched", res)
	}
}