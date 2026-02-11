package handlers

import (
	"time"
	"foca-store/database"
	"foca-store/helper"
	"foca-store/models"
	"foca-store/request"
	"foca-store/response"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req request.AuthRequest
	c.ShouldBindJSON(&req)

	hash, _ := helper.HashPassword(req.Password)
	user := models.User{Email: req.Email, Password: hash, Role: "USER"}
	database.DB.Create(&user)

	response.Success(c, "registered", user)
}

func Login(c *gin.Context) {
	var req request.AuthRequest
	c.ShouldBindJSON(&req)

	var u models.User
	database.DB.Where("email = ?", req.Email).First(&u)

	if !helper.CheckPassword(u.Password, req.Password) {
		response.Error(c, 401, "invalid credential")
		return
	}

	access, _ := helper.GenerateToken(u.ID, 15*time.Minute)
	refresh, _ := helper.GenerateToken(u.ID, 7*24*time.Hour)

	response.Success(c, "login success", gin.H{
		"access_token": access,
		"refresh_token": refresh,
	})
}
