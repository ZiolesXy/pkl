package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"focastore/database"
	"focastore/handlers"
	"focastore/middleware"
)

func main() {
	_ = godotenv.Load()

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("db migrate failed: %v", err)
	}

	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "ok", "data": gin.H{}})
		})

		authHandler := handlers.NewAuthHandler(db)
		productHandler := handlers.NewProductHandler(db)
		cartHandler := handlers.NewCartHandler(db)
		checkoutHandler := handlers.NewCheckoutHandler(db)
		seedHandler := handlers.NewSeedHandler(db)

		api.GET("/seed", seedHandler.Seed)

		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/refresh", authHandler.Refresh)

		protected := api.Group("")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			protected.GET("/me", authHandler.Me)

			protected.GET("/products", productHandler.List)
			protected.GET("/products/:id", productHandler.Get)

			protected.GET("/cart", cartHandler.GetCart)
			protected.POST("/cart/items", cartHandler.AddItem)
			protected.PUT("/cart/items/:id", cartHandler.UpdateItem)
			protected.DELETE("/cart/items/:id", cartHandler.DeleteItem)

			protected.POST("/checkout", checkoutHandler.Checkout)
			protected.GET("/checkouts", checkoutHandler.ListMyCheckouts)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuthMiddleware(), middleware.RequireRole("Admin"))
		{
			admin.POST("/products", productHandler.Create)
			admin.PUT("/products/:id", productHandler.Update)
			admin.DELETE("/products/:id", productHandler.Delete)

			admin.GET("/checkouts", checkoutHandler.ListAllCheckouts)
			admin.PUT("/checkouts/:id/status", checkoutHandler.UpdateStatus)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
