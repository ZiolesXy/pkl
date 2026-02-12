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

type CartHandler struct {
	db *gorm.DB
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{db: db}
}

func (h *CartHandler) getOrCreateCart(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := h.db.Where("user_id = ?", userID).First(&cart).Error
	if err == nil {
		return &cart, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	cart = models.Cart{UserID: userID}
	if err := h.db.Create(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (h *CartHandler) GetCart(c *gin.Context) {
	uidAny, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, _ := uidAny.(uint)

	cart, err := h.getOrCreateCart(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}

	var full models.Cart
	if err := h.db.Preload("Items.Product").First(&full, cart.ID).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}

	response.Success(c, "cart", full)
}

func (h *CartHandler) AddItem(c *gin.Context) {
	uidAny, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, _ := uidAny.(uint)

	var req request.AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	cart, err := h.getOrCreateCart(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}

	var product models.Product
	if err := h.db.First(&product, req.ProductID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "product not found")
		return
	}
	if product.Stock < req.Quantity {
		response.Error(c, http.StatusBadRequest, "not enough stock")
		return
	}

	var item models.CartItem
	err = h.db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&item).Error
	if err == nil {
		newQty := item.Quantity + req.Quantity
		if product.Stock < newQty {
			response.Error(c, http.StatusBadRequest, "not enough stock")
			return
		}
		item.Quantity = newQty
		if err := h.db.Save(&item).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, "update cart item failed")
			return
		}
		h.db.Preload("Product").First(&item, item.ID)
		response.Success(c, "item updated", item)
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}

	item = models.CartItem{CartID: cart.ID, ProductID: req.ProductID, Quantity: req.Quantity}
	if err := h.db.Create(&item).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "add item failed")
		return
	}
	if err := h.db.Preload("Product").First(&item, item.ID).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}
	response.Created(c, "item added", item)
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	uidAny, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, _ := uidAny.(uint)

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req request.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	cart, err := h.getOrCreateCart(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}

	var item models.CartItem
	if err := h.db.Preload("Product").Where("id = ? AND cart_id = ?", uint(itemID), cart.ID).First(&item).Error; err != nil {
		response.Error(c, http.StatusNotFound, "item not found")
		return
	}

	if item.Product.Stock < req.Quantity {
		response.Error(c, http.StatusBadRequest, "not enough stock")
		return
	}

	item.Quantity = req.Quantity
	if err := h.db.Save(&item).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "update item failed")
		return
	}
	response.Success(c, "item updated", item)
}

func (h *CartHandler) DeleteItem(c *gin.Context) {
	uidAny, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userID, _ := uidAny.(uint)

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	cart, err := h.getOrCreateCart(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}

	if err := h.db.Where("id = ? AND cart_id = ?", uint(itemID), cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "delete item failed")
		return
	}
	response.Success(c, "item deleted", gin.H{})
}
