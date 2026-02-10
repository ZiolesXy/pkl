package main

import (
	"fmt"
	"log"
	"net/http"

	// Import internal packages sesuai struktur folder kita
	"main/internal/controller"
	"main/internal/database"
	"main/internal/middleware"
	"main/internal/repository"
	"main/internal/service"
	"main/pkg/config"
	"main/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Environment Variables (.env)
	config.LoadConfig()

	// 2. Inisialisasi Database & Migrasi Tabel
	database.ConnectDB()

	// 3. Setup Dependency Injection (DI)
	
	// --- Modul Auth ---
	authRepo := repository.NewAuthRepository(database.DB)
	authService := service.NewAuthService(authRepo)
	authController := controller.NewAuthController(authService)

	// --- Modul Product & Category ---
	productRepo := repository.NewProductRepository(database.DB)
	productService := service.NewProductService(productRepo)
	productController := controller.NewProductController(productService)

	// 4. Inisialisasi Framework Gin
	r := gin.Default()

	// 5. Middleware Global
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// 6. Routing API
	api := r.Group("/api")
	{
		// --- PUBLIC ROUTES (Tanpa Login) ---
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Katalog Produk (Bisa dilihat siapa saja)
		api.GET("/products", productController.GetAllProducts)
		api.GET("/categories", productController.GetAllCategories)

		// --- PROTECTED ROUTES (Harus Login) ---
		userRoutes := api.Group("/user")
		userRoutes.Use(middleware.AuthMiddleware())
		{
			userRoutes.GET("/profile", func(c *gin.Context) {
				userID, _ := c.Get("userID")
				role, _ := c.Get("role")
				utils.SuccessResponse(c, http.StatusOK, "Profile data", gin.H{
					"user_id": userID,
					"role":    role,
				})
			})
			
			// Nanti STEP 8 (Cart) akan masuk ke sini
		}

		// --- ADMIN ROUTES (Login + Role Admin) ---
		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
		{
			// CRUD Category
			adminRoutes.POST("/categories", productController.CreateCategory)
			
			// CRUD Product
			adminRoutes.POST("/products", productController.CreateProduct)
			// Kamu bisa menambahkan adminRoutes.PUT atau adminRoutes.DELETE di sini nanti
			
			adminRoutes.GET("/dashboard", func(c *gin.Context) {
				utils.SuccessResponse(c, http.StatusOK, "Welcome Admin Dashboard", nil)
			})
		}
	}

	// Health Check Endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "up",
			"message": "Foca Store API is ready",
		})
	})

	// 7. Start Server
	port := config.GetEnv("PORT", "8080")
	fmt.Printf("\n--- Foca Store API Running on Port %s ---\n", port)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}