package main

import (
	"main/database"
	"main/handlers"
	"main/models"
	"main/trash"
	"net/http"
	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "http://localhost:3000" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	database.ConnedtDB()

	database.DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Barang{},
	)

	// seeders.RunSeed()

	r := gin.Default()
	r.Use(corsMiddleware())
	r.POST("/dummy", handlers.RunSeeder)
	r.DELETE("/reset", handlers.ClearSeeder)
	
	r.GET("/usert", trash.GetUsers)
	r.GET("/usert/:id", trash.GetUserByID)

	r.POST("/roles", handlers.CreateRole)
	r.GET("/roles", handlers.GetRole)
	r.PUT("/role/:id", handlers.UpdateRole)
	r.DELETE("/role/:id", handlers.DelRole)

	r.POST("/users", handlers.CreateUser)
	r.GET("/users", handlers.GetUsers)
	r.GET("/user/:id", handlers.GetUserByID)
	r.PUT("/user/:id", handlers.UpdateUsers)
	r.DELETE("/user/:id", handlers.DelUser)

	r.POST("/barangs", handlers.CreateBarang)
	r.GET("/barangs", handlers.GetBarangs)
	r.GET("/barang/:id", handlers.GetBarangByID)
	r.PUT("/barang/:id", handlers.UpdateBarang)
	r.DELETE("/barang/:id", handlers.DelBarang)

	r.GET("/users/barangs", handlers.GetUserBarangs)
	r.GET("/user/barang", handlers.GetUserBarangPivot)
	r.POST("/user/:user_id/barang/:barang_id", handlers.AssignBarang)
	r.DELETE("/user/:id/barang/:barang_id", handlers.RemoveBarang)
	r.Run(":8080")
}