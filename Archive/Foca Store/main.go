package main

import (
	"log"
	"main/config"
	"main/internal/app"
	"main/internal/database"
)

func main() {
	cfg := config.LoadConfig()
	database.ConnectPostgres(cfg)
	router := app.SetupRouter()

	log.Printf("Server running on port %s\n", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}