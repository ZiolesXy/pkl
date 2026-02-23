package handlers

import (
	"net/http"
	"time"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ClaimCoupon(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID := userIDRaw.(uint)

		var req request.ClaimCouponRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		// Find coupon by code
		var coupon models.Coupon
		if err := db.Where("code = ?", req.CouponCode).First(&coupon).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "coupon not found")
			return
		}

		// Validate active
		if coupon.IsActive != nil && !*coupon.IsActive {
			response.ErrorResponse(c, http.StatusBadRequest, "coupon is not active")
			return
		}

		// Validate not expired
		if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(time.Now()) {
			response.ErrorResponse(c, http.StatusBadRequest, "coupon has expired")
			return
		}

		// Check duplicate claim
		var existing models.UserCoupon
		err := db.Where("user_id = ? AND coupon_id = ?", userID, coupon.ID).First(&existing).Error
		if err == nil {
			response.ErrorResponse(c, http.StatusConflict, "coupon already claimed")
			return
		}

		// Create claim
		userCoupon := models.UserCoupon{
			UserID:   userID,
			CouponID: coupon.ID,
		}

		if err := db.Create(&userCoupon).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to claim coupon")
			return
		}

		// Reload with coupon data
		db.Preload("Coupon").First(&userCoupon, userCoupon.ID)

		res := response.BuildUserCouponResponse(userCoupon)
		response.SuccessResponse(c, "Coupon claimed", res)
	}
}

func GetMyCoupons(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID := userIDRaw.(uint)

		var userCoupons []models.UserCoupon
		if err := db.
			Preload("Coupon").
			Where("user_id = ?", userID).
			Order("created_at DESC").
			Find(&userCoupons).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch coupons")
			return
		}

		res := response.BuildUserCouponListResponse(userCoupons)
		response.SuccessListResponse(c, "Your coupons fetched", res)
	}
}

func RemoveCoupon(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("user_id")
		if !exists {
			response.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userID := userIDRaw.(uint)

		id := c.Param("id")

		var userCoupon models.UserCoupon
		if err := db.Where("id = ? AND user_id = ?", id, userID).First(&userCoupon).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "coupon not found")
			return
		}

		// Cannot remove if already used
		if userCoupon.UsedAt != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "cannot remove used coupon")
			return
		}

		if err := db.Delete(&userCoupon).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "failed to remove coupon")
			return
		}

		response.SuccessResponse(c, "Coupon removed", nil)
	}
}
