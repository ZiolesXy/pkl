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

func UpdateUsers(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"Error": "User tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(200, user)
}

func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var role models.Role

	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(404, gin.H{"Error": "Role tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Save(&role).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(200, role)
}

func UpdateBarang(c *gin.Context) {
	id := c.Param("id")
	var barang models.Barang

	if err := database.DB.First(&barang, id).Error; err != nil {
		c.JSON(400, gin.H{"Error": "Barang tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&barang); err != nil {
		c.JSON(404, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Save(&barang).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(200, barang)
}

func DelUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"Error": "User tidak ditemukan"})
		return
	}

	if err := database.DB.Model(&user).Association("Barangs").Clear(); err != nil{
		c.JSON(400, gin.H{"Error": err})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"Messege" : "User berhasil dihapus"})
}

func DelRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role

	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(404, gin.H{"Error": "User tidak ditemukan"})
		return
	}

	
	if erro := database.DB.Delete(&role).Error; erro != nil {
		c.JSON(400, gin.H{"Error": erro.Error()})
		return
	}
	c.JSON(200, gin.H{"Messege" : "User berhasil dihapus"})
}

func DelBarang(c *gin.Context) {
	id := c.Param("id")
	var barang models.Barang

	if err := database.DB.First(&barang, id).Error; err != nil {
		c.JSON(404, gin.H{"Error": "Barang tidak ditemukan"})
		return
	}

	if err := database.DB.Delete(&barang).Error; err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"messege": "Barang berhasil dihapus"})
}