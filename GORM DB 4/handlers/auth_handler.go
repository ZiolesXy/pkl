package handlers

import (
	"main/database"
	"main/helpers"
	"main/models"
	"main/request"
	"main/respons"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req request.RegisterRequest
	c.ShouldBindJSON(&req)

	hashed, err := helpers.HashPassword(req.Password)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			respons.NewJsonResponse("Hash error", nil),
		)
		return
	}

	user := models.User{
		Name: req.Name,
		Email: req.Email,
		Password: hashed,
		RoleID: 2,
	}

	if err := database.DB.Preload("roles").Create(&user).Error; err != nil {
		c.JSON(
			http.StatusBadRequest,
			request.NewJsonResponse("Error", err.Error()),
		)
		return
	}

	if err := database.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			respons.NewJsonResponse("Error", err.Error()),
		)
	}

	c.JSON(http.StatusOK, request.NewJsonResponse("Succes", respons.User{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		Role: respons.Role{
			ID: user.Role.ID,
			Name: user.Role.Name,
		},
	}))
}

func Login(c *gin.Context) {
	var req request.LoginRequest
	c.ShouldBindJSON(&req)

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, request.NewJsonResponse("Error", err.Error()))
	}

	if user.ID == 0 {
		c.JSON(401, respons.NewJsonResponse("User not found", nil))
	}

	if err := helpers.CheckPassword(user.Password, req.Password); err != nil {
		c.JSON(http.StatusForbidden, respons.NewJsonResponse("Wrong Password", nil))
		return
	}

	token, _ := helpers.GenerateJWT(user)

	c.JSON(
		200,
		respons.NewJsonResponse("token", token),
	)
}