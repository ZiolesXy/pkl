package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"focastore/middleware"
	"focastore/models"
	"focastore/request"
	"focastore/response"
)

type CheckoutHandler struct {
	db *gorm.DB
}

func NewCheckoutHandler(db *gorm.DB) *CheckoutHandler {
	return &CheckoutHandler{db: db}
}

func (h *CheckoutHandler) Checkout(c *gin.Context) {
	uidAny, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, _ := uidAny.(uint)

	var cart models.Cart
	if err := h.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusBadRequest, "cart is empty")
			return
		}
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}

	var items []models.CartItem
	if err := h.db.Preload("Product").Where("cart_id = ?", cart.ID).Find(&items).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}
	if len(items) == 0 {
		response.Error(c, http.StatusBadRequest, "cart is empty")
		return
	}

	var total int64
	for _, it := range items {
		if it.Product.ID == 0 {
			response.Error(c, http.StatusBadRequest, "product not found")
			return
		}
		if it.Product.Stock < it.Quantity {
			response.Error(c, http.StatusBadRequest, "not enough stock")
			return
		}
		total += it.Product.Price * it.Quantity
	}

	var checkout models.Checkout
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		for _, it := range items {
			res := tx.Model(&models.Product{}).Where("id = ? AND stock >= ?", it.ProductID, it.Quantity).
				Update("stock", gorm.Expr("stock - ?", it.Quantity))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return errors.New("not enough stock")
			}
		}

		checkout = models.Checkout{UserID: userID, TotalAmount: total, Status: models.CheckoutPending}
		if err := tx.Create(&checkout).Error; err != nil {
			return err
		}

		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err.Error() == "not enough stock" {
			response.Error(c, http.StatusBadRequest, "not enough stock")
			return
		}
		response.Error(c, http.StatusInternalServerError, "checkout failed")
		return
	}

	response.Created(c, "checkout created", checkout)
}

func (h *CheckoutHandler) ListMyCheckouts(c *gin.Context) {
	uidAny, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, _ := uidAny.(uint)

	var list []models.Checkout
	if err := h.db.Where("user_id = ?", userID).Order("id desc").Find(&list).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}
	response.Success(c, "checkouts", gin.H{"items": list})
}

func (h *CheckoutHandler) ListAllCheckouts(c *gin.Context) {
	var list []models.Checkout
	if err := h.db.Order("id desc").Find(&list).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}
	response.Success(c, "checkouts", gin.H{"items": list})
}

func (h *CheckoutHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req request.UpdateCheckoutStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	var checkout models.Checkout
	if err := h.db.First(&checkout, uint(id)).Error; err != nil {
		response.Error(c, http.StatusNotFound, "checkout not found")
		return
	}

	checkout.Status = models.CheckoutStatus(req.Status)
	if err := h.db.Save(&checkout).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "update status failed")
		return
	}
	response.Success(c, "status updated", checkout)
}
