package main

import (
	"main/database"
	"main/handlers"
	"main/middlewares"
	"main/models"
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

	//public
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	//auth
	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	admin := auth.Group("/")
	admin.Use(middlewares.OnlyAdmin())

	r.POST("/dummy", handlers.RunSeeder)
	r.DELETE("/reset", handlers.ClearSeeder)
	auth.GET("/", handlers.GetUserBarangs)

	admin.POST("/roles", handlers.CreateRole)
	auth.GET("/roles", handlers.GetRole)
	auth.GET("/role/:id", handlers.GetRoleByID)
	admin.PUT("/role/:id", handlers.UpdateRole)
	admin.DELETE("/role/:id", handlers.DelRole)
	
	auth.GET("/users", handlers.GetUsers)
	auth.GET("/user/:id", handlers.GetUserByID)
	admin.PUT("/user/:id", handlers.UpdateUsers)
	admin.DELETE("/user/:id", handlers.DelUser)

	admin.POST("/barangs", handlers.CreateBarang)
	auth.GET("/barangs", handlers.GetBarangs)
	auth.GET("/barang/:id", handlers.GetBarangByID)
	admin.PUT("/barang/:id", handlers.UpdateBarang)
	admin.DELETE("/barang/:id", handlers.DelBarang)

	auth.GET("/users/barangs", handlers.GetUserBarangs)
	auth.GET("/user/barang", handlers.GetUserBarangPivot)
	admin.POST("/user/:user_id/barang/:barang_id", handlers.AssignBarang)
	admin.DELETE("/user/:id/barang/:barang_id", handlers.RemoveBarang)
	r.Run(":3605")
}