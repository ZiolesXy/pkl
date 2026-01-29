package seeders

import (
	"main/database"
	"main/models"

	"github.com/gin-gonic/gin"
)

func RunSeed(c *gin.Context) {

	//Role
	roles := []models.Role{
		{Name: "Admin"},
		{Name: "User"},
		{Name: "Sales"},
	}
	if err := database.DB.Create(&roles).Error; err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	//Users
	users := []models.User{
		{
			Name: "Pasha",
			RoleID: roles[0].ID,
		},
		{
			Name: "Amrl",
			RoleID: roles[1].ID,
		},
		{
			Name: "Siti",
			RoleID: roles[2].ID,
		},
	}
	if err := database.DB.Create(&users).Error; err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	//Barang
	barangs := []models.Barang{
		{Name: "Laptop"},
		{Name: "Mouse"},
		{Name: "Keyboard"},
		{Name: "Monitor"},
		{Name: "Printer"},
	}
	if err := database.DB.Create(&barangs).Error; err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	
	//Pivot
	if err := database.DB.Model(&users[0]).Association("Barangs").Append(&barangs[0], &barangs[1]); err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	if err := database.DB.Model(&users[1]).Association("Barangs").Append(&barangs[0], &barangs[1]); err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
}

func ClearData(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	for _, user := range users {
		if err := database.DB.Model(&user).Association("Barangs").Clear(); err != nil {
			c.JSON(400, gin.H{"Error1": err.Error()})
			return
		}
		if err := database.DB.Unscoped().Delete(&user).Error; err != nil {
			c.JSON(400, gin.H{"Error2": err.Error()})
			return
		}
		// if err := database.DB.Unscoped().Delete(&models.Role{}).Error; err != nil {
		// 	c.JSON(400, gin.H{"Error3": err.Error()})
		// 	return
		// }
	}

	// if err := database.DB.Unscoped().Delete(&models.User{}).Error; err != nil {
	// 	c.JSON(400, gin.H{"Error2": err.Error()})
	// 	return
	// }

	// if err := database.DB.Unscoped().Delete(&models.Role{}).Error; err != nil {
	// 	c.JSON(400, gin.H{"Error3": err.Error()})
	// 	return
	// }

	// c.JSON(200, gin.H{
	// 	"messege": "All table data cleared",
	// })
}