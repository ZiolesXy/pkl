package main

import (
	"foca-store/database"
	"foca-store/handlers"
	"foca-store/middleware"
	"foca-store/models"
	"foca-store/seeder"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	database.DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Checkout{},
		&models.CheckoutItem{},
	)

	r := gin.Default()
	r.POST("/seed", seeder.RunSeed)
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	auth := r.Group("/", middleware.JWTAuth())
	auth.GET("/products", handlers.GetProducts)
	auth.POST("/cart", handlers.AddToCart)
	auth.POST("/checkout", handlers.CreateCheckout)

	admin := auth.Group("/admin", middleware.AdminOnly())
	admin.POST("/products", handlers.CreateProduct)
	admin.POST("/checkout/:id/status", handlers.UpdateCheckoutStatus)

	r.Run(":3605")
}