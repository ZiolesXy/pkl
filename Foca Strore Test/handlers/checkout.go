package handlers

import (
	"net/http"
	"time"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Checkout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// AUTH
		userIDRaw, exist := c.Get("user_id")
		if !exist {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		userID := userIDRaw.(uint)

		// REQUEST
		var req request.CheckoutRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		// TRANSACTION START
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// GET CART
		var cart models.Cart

		if err := tx.
			Where("user_id = ?", userID).
			First(&cart).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusNotFound, "cart not found")
			return
		}

		// GET CART ITEMS + LOCK
		var items []models.CartItem

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Product").
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

		// CALCULATE TOTAL + UPDATE STOCK
		totalPrice := 0.0

		for _, item := range items {
			if item.Product.Stock < item.Quantity {
				tx.Rollback()
				response.ErrorResponse(
					c,
					http.StatusBadRequest,
					"insufficient stock for "+item.Product.Name,
				)
				return
			}

			err := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock - ?", item.Quantity)).
				Error

			if err != nil {
				tx.Rollback()
				response.ErrorResponse(c, 500, "failed update stock")
				return
			}

			totalPrice += item.Product.Price * float64(item.Quantity)

		}

		// APPLY COUPON
		var coupon models.Coupon
		var userCoupon models.UserCoupon

		if req.CouponCode != nil && *req.CouponCode != "" {

			err := tx.
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("code = ?", *req.CouponCode).
				First(&coupon).Error

			if err != nil {
				tx.Rollback()
				response.ErrorResponse(c, 400, "invalid coupon")
				return
			}

			// VALIDATE COUPON OWNERSHIP
			if err := tx.
				Where("user_id = ? AND coupon_id = ?", userID, coupon.ID).
				First(&userCoupon).Error; err != nil {
				tx.Rollback()
				response.ErrorResponse(c, 400, "coupon not claimed")
				return
			}

			if userCoupon.UsedAt != nil {
				tx.Rollback()
				response.ErrorResponse(c, 400, "coupon already used")
				return
			}

			// APPLY DISCOUNT
			if coupon.Type == "percentage" {
				totalPrice -= totalPrice * coupon.Value / 100
			}

			if coupon.Type == "fixed" {
				totalPrice -= coupon.Value
			}

			if totalPrice < 0 {
				totalPrice = 0
			}

			// SAFE INCREMENT QUOTA
			result := tx.Model(&models.Coupon{}).
				Where("id = ? AND used_count < quota", coupon.ID).
				Update("used_count", gorm.Expr("used_count + 1"))

			if result.RowsAffected == 0 {
				tx.Rollback()
				response.ErrorResponse(c, 400, "coupon quota exceeded")
				return
			}

			// MARK COUPON AS USED
			now := time.Now()
			if err := tx.Model(&userCoupon).Update("used_at", now).Error; err != nil {
				tx.Rollback()
				response.ErrorResponse(c, 500, "failed to mark coupon as used")
				return
			}
		}

		// LOOKUP ADDRESS
		var address models.Address
		if err := tx.
			Where("uid = ? AND user_id = ?", req.AddressUID, userID).
			First(&address).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, http.StatusNotFound, "address not found")
			return
		}

		// CREATE CHECKOUT
		checkout := models.Checkout{
			UserID:     userID,
			AddressID:  &address.ID,
			TotalPrice: totalPrice,
			Status:     "pending",
		}

		if coupon.ID != 0 {
			checkout.CouponID = &coupon.ID
		}

		if err := tx.Create(&checkout).Error; err != nil {
			tx.Rollback()
			response.ErrorResponse(c, 500, "failed create checkout")
			return
		}

		// CREATE CHECKOUT ITEMS
		for _, item := range items {

			checkoutItem := models.CheckoutItem{
				CheckoutID: checkout.ID,
				ProductID:  item.ProductID,
				Quantity:   item.Quantity,
				Price:      item.Product.Price,
			}

			if err := tx.Create(&checkoutItem).Error; err != nil {
				tx.Rollback()
				response.ErrorResponse(c, 500, "failed create checkout item")
				return
			}

		}

		// DELETE CART ITEMS
		if err := tx.
			Where("id IN ?", req.CartItemIDs).
			Delete(&models.CartItem{}).Error; err != nil {
				tx.Rollback()
				response.ErrorResponse(c, 500, "failed delete cart")
				return
			}

		// COMMIT
		if err := tx.Commit().Error; err != nil {
			response.ErrorResponse(c, 500, "transaction failed")
			return
		}

		// RELOAD RESULT
		var result models.Checkout

		db.
			Preload("User").
			Preload("Coupon").
			Preload("Address").
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
			Preload("Address").
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
			Preload("Address").
			Preload("Items").
			Preload("Items.Product").
			Preload("Coupon").
			Preload("UserCoupon").
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

		checkouts := []models.Checkout{}

		query := db.
			Preload("User").
			Preload("Coupon").
			Preload("Address").
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

func GetMyCheckout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		userID := userIDRaw.(uint)

		var checkouts []models.Checkout

		if err := db.
			Preload("User").
			Preload("Coupon").
			Preload("Address").
			Preload("Items").
			Preload("Items.Product").
			Where("user_id = ?", userID).
			Order("created_at DESC").
			Find(&checkouts).Error; err != nil {

			response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch checkout")
			return
		}

		res := response.BuildCheckOutListResponse(checkouts)

		response.SuccessListResponse(c, "your checkout list fetched", res)
	}
}
