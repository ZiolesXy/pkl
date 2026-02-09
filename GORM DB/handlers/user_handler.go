package handlers

import (
	"main/database"
	"main/helpers"
	"main/models"
	"main/request"
	"main/respons"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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

	entries := respons.NewEntries(userResponses)
	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Success", entries),
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

func GetProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Authorization header missing", nil),
		)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Invalid authorization format", nil),
		)
		return
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return helpers.ACCESS_SECRET, nil
	})

	if err != nil || !token.Valid {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Invalid or expired token", nil),
		)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Invalid token claims", nil),
		)
		return
	}

	if claims["type"] != "access" {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Invalid token type", nil),
		)
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Invalid user id", nil),
		)
		return
	}

	userID := uint(userIDFloat)

	var user models.User
	if err := database.DB.Preload("Role").Preload("Barangs").First(&user, userID).Error; err != nil {
		c.JSON(
			http.StatusNotFound,
			respons.NewJsonResponse("User not found", nil),
		)
		return
	}

	userResp := respons.Profile{
		ID: userID,
		Name: user.Name,
		Email: user.Email,
	}

	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Success", userResp),
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
