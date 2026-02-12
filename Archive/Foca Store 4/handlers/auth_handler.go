package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"focastore/helper"
	"focastore/middleware"
	"focastore/models"
	"focastore/request"
	"focastore/response"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	var userRole models.Role
	if err := h.db.Where("name = ?", "User").First(&userRole).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "role not ready, seed first")
		return
	}

	var exists int64
	if err := h.db.Model(&models.User{}).Where("email = ?", req.Email).Count(&exists).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "database error")
		return
	}
	if exists > 0 {
		response.Error(c, http.StatusBadRequest, "email already registered")
		return
	}

	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "password error")
		return
	}

	user := models.User{Name: req.Name, Email: req.Email, PasswordHash: hash, RoleID: userRole.ID}
	if err := h.db.Create(&user).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "create user failed")
		return
	}

	data := gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": userRole.Name}
	response.Created(c, "registered", data)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	var user models.User
	if err := h.db.Preload("Role").Where("email = ?", req.Email).First(&user).Error; err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if !helper.CheckPassword(user.PasswordHash, req.Password) {
		response.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	access, err := helper.GenerateToken(user.ID, user.Role.Name, helper.TokenAccess, 15*time.Minute)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "token error")
		return
	}
	refresh, err := helper.GenerateToken(user.ID, user.Role.Name, helper.TokenRefresh, 7*24*time.Hour)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "token error")
		return
	}

	response.Success(c, "login success", gin.H{"access_token": access, "refresh_token": refresh})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req request.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	claims, err := helper.ParseToken(req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid token")
		return
	}
	if claims.Type != helper.TokenRefresh {
		response.Error(c, http.StatusUnauthorized, "invalid token type")
		return
	}

	access, err := helper.GenerateToken(claims.UserID, claims.Role, helper.TokenAccess, 15*time.Minute)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "token error")
		return
	}
	refresh, err := helper.GenerateToken(claims.UserID, claims.Role, helper.TokenRefresh, 7*24*time.Hour)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "token error")
		return
	}

	response.Success(c, "refreshed", gin.H{"access_token": access, "refresh_token": refresh})
}

func (h *AuthHandler) Me(c *gin.Context) {
	uidAny, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	uid, _ := uidAny.(uint)

	var user models.User
	if err := h.db.Preload("Role").First(&user, uid).Error; err != nil {
		response.Error(c, http.StatusNotFound, "user not found")
		return
	}

	response.Success(c, "me", gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role.Name})
}
