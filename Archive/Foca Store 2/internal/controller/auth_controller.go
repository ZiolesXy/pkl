package controller

import (
	"main/internal/model"
	"main/internal/service"
	"main/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) *AuthController {
	return &AuthController{service}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var input model.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	if err := ctrl.service.Register(input); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Registration failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", nil)
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	token, err := ctrl.service.Login(input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{"token": token})
}