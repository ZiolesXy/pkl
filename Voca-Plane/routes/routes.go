package routes

import (
	"net/http"
	"voca-plane/internal/handler"
	"voca-plane/middleware"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine,
	authHandler *handler.AuthHandler,
	flightHandler *handler.FlightHandler,
	transactionHandler *handler.TransactionHandler,
	userHandler *handler.UserHandler,
	adminHandler *handler.AdminHandler,
	systemHandler *handler.SystemHandler,
	jwtSecret string,
	allowedOrigins string,
	appPassword string) {
	r.Use(middleware.CORS(allowedOrigins))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// System Routes
		sys := v1.Group("/system")
		sys.Use(middleware.JWTAuth(jwtSecret))
		sys.Use(middleware.RequireSuperAdmin())
		sys.Use(middleware.AppPassword(appPassword))
		{
			sys.POST("/seed", systemHandler.Seed)
		}

		// Public Auth
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/auth/refresh", authHandler.RefreshToken)
		v1.POST("/transactions/midtrans/callback", transactionHandler.MidtransCallback)

		// Public Flight Search
		v1.GET("/flights", flightHandler.GetAll)
		v1.GET("/flights/search", flightHandler.Search)
		v1.GET("/flights/:id", flightHandler.GetByID)

		// Protected User Routes
		userProtected := v1.Group("")
		userProtected.Use(middleware.JWTAuth(jwtSecret))
		{
			userProtected.GET("/user/profile", userHandler.GetProfile)
			userProtected.PATCH("/user/profile", userHandler.UpdateProfile)
			v1.GET("/user/device-info", userHandler.GetDeviceInfo)

			userProtected.GET("/transactions", transactionHandler.GetList)
			userProtected.GET("/transactions/:code", transactionHandler.GetByCode)
			userProtected.POST("/transactions", transactionHandler.Create)
			userProtected.PATCH("/transactions/:code/pay", transactionHandler.Pay)
			userProtected.DELETE("/transactions/:code", transactionHandler.Cancel)
		}

		// Protected Admin Routes
		adminProtected := v1.Group("/admin")
		adminProtected.Use(middleware.JWTAuth(jwtSecret))
		adminProtected.Use(middleware.RequireAdmin())
		{
			adminProtected.GET("/dashboard", adminHandler.GetDashboard)

			adminProtected.GET("/users", adminHandler.GetUsers)
			adminProtected.PATCH("/users/:id/role", adminHandler.UpdateUserRole)

			adminProtected.GET("/transactions", adminHandler.GetTransactions)

			adminProtected.GET("/flights", adminHandler.GetFlights)
			adminProtected.POST("/flights", adminHandler.CreateFlight)
			adminProtected.PUT("/flights/:id", adminHandler.UpdateFlight)
			adminProtected.DELETE("/flights/:id", adminHandler.DeleteFlight)

			adminProtected.GET("/airlines", adminHandler.GetAirlines)
			adminProtected.POST("/airlines", adminHandler.CreateAirline)
			adminProtected.PUT("/airlines/:id", adminHandler.UpdateAirline)
			adminProtected.DELETE("/airlines/:id", adminHandler.DeleteAirline)

			adminProtected.GET("/airports", adminHandler.GetAirports)
			adminProtected.POST("/airports", adminHandler.CreateAirport)
			adminProtected.PUT("/airports/:id", adminHandler.UpdateAirport)
			adminProtected.DELETE("/airports/:id", adminHandler.DeleteAirport)

			adminProtected.GET("/promos", adminHandler.GetPromos)
			adminProtected.POST("/promos", adminHandler.CreatePromo)
			adminProtected.PUT("/promos/:id", adminHandler.UpdatePromo)
			adminProtected.DELETE("/promos/:id", adminHandler.DeletePromo)
		}
	}
}
