package trash

import (
	"main/database"
	"main/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var users []models.User

	if err := database.DB.Preload("Role").Preload("Barangs").Find(&users).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := database.DB.Preload("Role").First(&user, "id = ?", id).Error; err != nil {
		c.JSON(400, gin.H{"Error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetBarangs(c *gin.Context) {
	var barangs []models.Barang
	
	if err := database.DB.Find(&barangs).Error; err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(200, barangs)
}

func GetUserBarangs(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := database.DB.Preload("Barangs").First(&user, userID).Error;err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetRole(c *gin.Context) {
	var role []models.Role
	if err := database.DB.Find(&role).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &user)
}

func CreateBarang(c *gin.Context) {
	var barang models.Barang

	if err := c.ShouldBindJSON(&barang); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Create(&barang).Error; err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, barang)
}

func CreateRole(c *gin.Context) {
	var role models.Role

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Create(&role).Error; err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}