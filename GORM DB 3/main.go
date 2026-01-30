package main

import (
	"main/database"
	"main/handlers"
	"main/models"
	"main/trash"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnedtDB()

	database.DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Barang{},
	)

	// seeders.RunSeed()

	r := gin.Default()
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
	r.PUT("/barang/:id", handlers.UpdateBarang)
	r.DELETE("/barang/:id", handlers.DelBarang)

	r.GET("/users/barangs", handlers.GetUserBarangs)
	r.POST("/user/:user_id/barang/:barang_id", handlers.AssignBarang)
	r.DELETE("/user/:id/barang/:barang_id", handlers.RemoveBarang)
	r.Run(":8080")
}