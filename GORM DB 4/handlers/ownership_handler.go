package handlers

import (
	"main/database"
	"main/models"
	"main/respons"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserBarangs(c *gin.Context) {
	var users []models.User

	if err := database.DB.Preload("Role").Preload("Barangs").Order("id ASC").Find(&users).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			respons.NewJsonResponse("Failed get users barangs", nil),
		)
		return
	}

	userResponses := []respons.UserWithBarang{}

	for _, user := range users {
		barangResp := []respons.Barang{}
		for _, barang := range user.Barangs{
			barangResp = append(barangResp, respons.Barang{
				ID: barang.ID,
				Name: barang.Name,
			})
		}

		roleResp := respons.Role{
			ID: user.Role.ID,
			Name: user.Role.Name,
		}

		userResponses = append(userResponses, respons.UserWithBarang{
			ID: user.ID,
			Name: user.Name,
			Email: user.Email,
			Role: roleResp,
			Barangs: barangResp,
		})
	}

	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Succes", userResponses),
	)
}

func GetUserBarangPivot (c *gin.Context) {
	var results []map[string]interface{}

	if err := database.DB.Table("user_barangs").Select("user_id, barang_id").Order("user_id ASC").Find(&results).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			respons.NewJsonResponse("Failed to fetch pivot data", err.Error()),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Succes", results),
	)
}

func AssignBarang(c *gin.Context) {
	userID := c.Param("user_id")
	barangID := c.Param("barang_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound,
			respons.NewJsonResponse("User not found", nil))
		return
	}

	var barang models.Barang
	if err := database.DB.First(&barang, barangID).Error; err != nil {
		c.JSON(http.StatusNotFound,
			respons.NewJsonResponse("Barang not found", nil))
		return
	}

	if err := database.DB.Model(&user).Association("Barangs").Append(&barang); err != nil {
		c.JSON(http.StatusInternalServerError,
			respons.NewJsonResponse("Failed add barang", nil))
		return
	}

	barangResp := respons.Barang{
		ID:   barang.ID,
		Name: barang.Name,
	}

	c.JSON(http.StatusOK,
		respons.NewJsonResponse("Barang added to user", respons.OwnerPost{
			ID:      user.ID,
			Name:    user.Name,
			Barangs: barangResp,
		}))
}

func RemoveBarang(c *gin.Context) {
	userID := c.Param("id")
	barangID := c.Param("barang_id")

	var user models.User
	var barang models.Barang

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{"Error": "User tidak ditemukan"})
		return
	}

	if err := database.DB.First(&barang, barangID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	if err := database.DB.Model(&user).Association("Barangs").Delete(&barang); err != nil {
		c.JSON(404, gin.H{"error": err.Error() })
		return
	}
	c.JSON(200, gin.H{"Messege": "Barang berhasil dilepas dari user"})
}