package handlers

import (
	"foca-store/database"
	"foca-store/helper"
	"foca-store/models"
	"foca-store/request"
	"foca-store/response"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req request.AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	var exist models.User
	if err := database.DB.Where("email = ?", req.Email).First(&exist).Error; err == nil {
		response.Error(c, 400, "email already registered")
		return
	}

	var role models.Role
	if err := database.DB.Where("name = ?", "User").First(&role).Error; err != nil {
		response.Error(c, 500, "default role not found")
		return
	}

	hashed, _ := helper.HashPassword(req.Password)

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		RoleID:   role.ID,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, "register success", user)
}


func Login(c *gin.Context) {
	var req request.AuthRequest
	c.ShouldBindJSON(&req)

	var u models.User
	database.DB.Preload("Role").Where("email = ?", req.Email).First(&u)

	if !helper.CheckPassword(u.Password, req.Password) {
		response.Error(c, 401, "invalid credential")
		return
	}

	access, refresh, err := helper.GenerateToken(u)
	if err != nil {
		response.Error(c, 500, "failed generate token")
		return
	}

	response.Success(c, "login success", gin.H{
		"access_token": access,
		"refresh_token": refresh,
	})
}
