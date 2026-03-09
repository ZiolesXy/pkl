package handler

import (
	"net/http"
	"strconv"
	"voca-plane/internal/domain/dto"
	"voca-plane/internal/service"
	"voca-plane/pkg/response"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	tx, err := h.service.CreateTransaction(c.Request.Context(), userID.(uint), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "transaction created", tx)
}

func (h *TransactionHandler) Pay(c *gin.Context) {
	code := c.Param("code")

	err := h.service.PayTransaction(c.Request.Context(), code)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "payment successful", nil)
}

func (h *TransactionHandler) GetList(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	transaction, total, err := h.service.GetUserTransactions(c.Request.Context(), userID.(uint), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "transaction retrieved", transaction, meta)
}

func(h *TransactionHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	transaction, err := h.service.GetTransactionByCode(c.Request.Context(), code)
	if err != nil {
		response.Error(c, http.StatusNotFound, "transaction not found")
		return
	}

	response.Success(c, http.StatusOK, "transaction found", transaction)
}

func (h *TransactionHandler) Cancel(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	code := c.Param("code")

	err := h.service.CancelTransaction(c.Request.Context(), userID.(uint), code)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "transaction cancelled", nil)
}