package handlers

import (
	"main/database"
	"main/models"
	"main/request"
	"main/respons"
	"main/seeders"
	"main/trash"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunSeeder(c *gin.Context) {
	seeders.RunSeed(c)
}

func ClearSeeder(c *gin.Context) {
	trash.ClearData(c)
	c.JSON(200, gin.H{"messege": "Data berhasil terhapus"})
}

func CreateUser(c *gin.Context) {
	var req request.UserPost

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			request.NewJsonResponse("Invalid request", err.Error()))
		return
	}

	var role models.Role
	if err := database.DB.First(&role, req.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			request.NewJsonResponse("Role not found", nil))
		return
	}

	user := models.User{
		Name:   req.Name,
		RoleID: req.RoleID,
		Role:   role,
	}

	if err := database.DB.Preload("Roles").Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			request.NewJsonResponse("Failed Create user", nil))
		return
	}

	c.JSON(http.StatusCreated,
		request.NewJsonResponse("User created", respons.User{
			ID:   user.ID,
			Name: user.Name,
			Role: respons.Role{
				ID:   user.Role.ID,
				Name: user.Role.Name,
			},
		}))
}

func GetUsers(c *gin.Context) {
	var users []models.User

	if err := database.DB.
		Preload("Role").
		Preload("Barangs").
		Order("id ASC").
		Find(&users).Error; err != nil {

		c.JSON(
			http.StatusNotFound,
			respons.NewJsonResponse("User not found", nil),
		)
		return
	}

	userResponses := []respons.User{}

	for _, user := range users {

		roleResp := respons.Role{
			ID:   user.Role.ID,
			Name: user.Role.Name,
		}

		barangResp := []respons.Barang{}
		for _, barang := range user.Barangs {
			barangResp = append(barangResp, respons.Barang{
				ID:   barang.ID,
				Name: barang.Name,
			})
		}

		userResponses = append(userResponses, respons.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  roleResp,
		})
	}

	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Success", userResponses),
	)
}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if err := database.DB.
		Preload("Role").
		Preload("Barangs").
		First(&user, id).Error; err != nil {
		c.JSON(
			404,
			respons.NewJsonResponse("User not found", nil),
		)
		return
	}

	roleResp := respons.Role{
		ID:   user.Role.ID,
		Name: user.Role.Name,
	}

	barangResp := []respons.Barang{}
	for _, barang := range user.Barangs {
		barangResp = append(barangResp, respons.Barang{
			ID:   barang.ID,
			Name: barang.Name,
		})
	}

	userResp := respons.User{
		ID:   user.ID,
		Name: user.Name,
		Email: user.Email,
		Role: roleResp,
	}

	c.JSON(
		200, respons.NewJsonResponse("Succes", userResp),
	)
}

func UpdateUsers(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := database.DB.Preload("Role").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, request.NewJsonResponse("User not found", nil))
		return
	}

	var req request.UserPut
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, request.NewJsonResponse("Invalid request", err.Error()))
		return
	}

	if req.Name != nil {
		if len(*req.Name) < 3 {
			c.JSON(http.StatusBadRequest, request.NewJsonResponse("Name too short", nil))
			return
		}
		user.Name = *req.Name
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.RoleID != nil {
		var role models.Role
		if err := database.DB.First(&role, *req.RoleID).Error; err != nil {
			c.JSON(http.StatusBadRequest, request.NewJsonResponse("Role not found", nil))
			return
		}
		user.RoleID = *req.RoleID
		user.Role = role
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, request.NewJsonResponse("Error", err.Error()))
		return
	}

	userResponse := respons.User{
		ID:   user.ID,
		Name: user.Name,
		Email: user.Email,
	}

	if user.Role.ID != 0 {
		userResponse.Role = respons.Role{
			ID:   user.Role.ID,
			Name: user.Role.Name,
		}
	}

	c.JSON(http.StatusOK, request.NewJsonResponse("User update", userResponse))
}

func DelUser(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound,
			request.NewJsonResponse("User not found", nil))
		return
	}

	c.JSON(http.StatusOK,
		request.NewJsonResponse("User deleted", nil))
}
