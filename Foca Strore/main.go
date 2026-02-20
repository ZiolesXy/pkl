package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"voca-store/database"
	"voca-store/handlers"
	"voca-store/helper"
	"voca-store/middleware"
	"voca-store/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Cloudinary
	if err := helper.InitCloudinary(); err != nil {
		log.Println("Warning: Cloudinary not initialized:", err)
	} else {
		log.Println("Cloudinary initialized successfully")
	}

	// Initialize database
	db := database.InitDB()
	defer database.CloseDB()

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Checkout{},
		&models.CheckoutItem{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Setup Gin router
	r := gin.Default()

	//cors set
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://172.16.17.79:172", "http://172.16.17.79:3000", "http://localhost:172", "https://undeliberatingly-decemviral-petronila.ngrok-free.dev "},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	authHandler := handlers.NewAuthHandler(db)
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.RefreshToken)
	r.GET("/category", handlers.GetAllCategory(db))
	r.GET("/product/:slug", handlers.GetProductBySlug(db))
	r.GET("/products", handlers.GetAllProducts(db))

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth(db))
	{
		// User routes
		protected.GET("/profile", handlers.GetProfile(db))
		protected.PUT("/profile", handlers.UpdateProfile(db))
		protected.GET("/cart", handlers.ViewCart(db))
		protected.POST("/cart/items", handlers.AddToCart(db))
		protected.DELETE("/cart/items/:id", handlers.RemoveCartItem(db))
		protected.POST("/checkout", handlers.Checkout(db))
		protected.GET("/checkout/me", handlers.GetMyCheckout(db))

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminOnly())
		{
			admin.POST("/category", handlers.CreateCategory(db))
			admin.PUT("/category/:id", handlers.UpdateCategory(db))
			admin.DELETE("/category/:id", handlers.DeleteCategory(db))
			admin.POST("/products", handlers.CreateProduct(db))
			admin.PUT("/products/:id", handlers.UpdateProduct(db))
			admin.DELETE("/products/:id", handlers.DeleteProduct(db))
			admin.DELETE("/products", handlers.DeleteAllProducts(db))
			admin.GET("/checkout", handlers.GetCheckout(db))
			admin.PATCH("/checkout/:id/approve", handlers.ApproveCheckout(db))
			admin.PATCH("/checkout/:id/reject", handlers.RejectCheckout(db))
		}
	}

	// Seeder endpoint
	seed := r.Group("/seed")
	{
		seed.GET("/assets", handlers.SeedProductsFromAssetsHandler(db))
		seed.GET("/roles", handlers.SeedRoleHandler(db))
		seed.GET("/admin", handlers.SeedAdminHandler(db))
		seed.GET("/users", handlers.SeedUsersHandler(db))
		seed.GET("/products", handlers.SeedProductsHandler(db))
		seed.GET("/all", handlers.SeedAllHandler(db))
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
