package handlers

import (
	"errors"
	"net/http"
	"voca-store/helper"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	var existingUser models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		response.ErrorResponse(c, http.StatusConflict, "Email already registered")
		return
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	var userRole models.Role
	if err := h.DB.Where("name = ?", "User").First(&userRole).Error; err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "User role not found")
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		RoleID:   userRole.ID,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	if err := h.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to load user data")
		return
	}

	cart := models.Cart{UserID: user.ID}
	if err := h.DB.Create(&cart).Error; err != nil {
		response.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"failed to create cart",
		)
		return
	}

	userResp := response.BuildUserResponse(user.ID, user.Name, user.Email, user.Role.Name)
	response.SuccessResponse(c, "user registered succesfully", userResp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invallid request body")
		return
	}

	var user models.User
	if err := h.DB.Preload("Role").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(c, http.StatusBadRequest, "invalid email or password")
		} else {
			response.ErrorResponse(c, http.StatusInternalServerError, "database error")
		}
		return
	}

	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		response.ErrorResponse(c, http. StatusBadRequest, "invalid email or password")
		return
	}

	accessToken, err := helper.GenerateAccessToken(user.ID, user.Role.Name)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get access token")
		return
	}

	refreshToken, err := helper.GenerateRefreshToken(user.ID, user.Role.Name)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get refresh token")
		return
	}

	userResp := response.BuildUserResponse(user.ID, user.Name, user.Email, user.Role.Name)
	authResp := response.BuildAuthResponse(userResp, accessToken, refreshToken)

	response.SuccessResponse(c, "login succesfull", authResp)
}

func(h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	claims, err := helper.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		response.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	var user models.User
	if err := h.DB.Preload("Role").First(&user, claims.UserID).Error; err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	accessToken, err := helper.GenerateAccessToken(user.ID, user.Role.Name)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to get access token")
		return
	}

	tokenResp := response.BuildToken(accessToken)
	response.SuccessResponse(c, "token refreshed succesfull", tokenResp)
}