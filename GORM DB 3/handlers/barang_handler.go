package handlers

import (
	"main/database"
	"main/models"
	"main/request"
	"main/respons"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateBarang(c *gin.Context) {
	var req request.BarangPost

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
		request.NewJsonResponse("Invalid request", err.Error()))
		return
	}

	barang := models.Barang{Name: req.Name}

	if err := database.DB.Create(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
		request.NewJsonResponse("Error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, request.NewJsonResponse("Succes", respons.Barang{
		ID: barang.ID,
		Name: barang.Name,
	}))
}

func GetBarangs(c *gin.Context) {
	var barangs []models.Barang

	if err := database.DB.Find(&barangs).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			respons.NewJsonResponse("Failed get barangs", nil),
		)
		return
	}

	barangResponse := []respons.Barang{}

	for _, barang := range barangs {
		barangResponse = append(barangResponse, respons.Barang{
			ID: barang.ID,
			Name: barang.Name,
		})
	}

	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Succes", barangResponse),
	)
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