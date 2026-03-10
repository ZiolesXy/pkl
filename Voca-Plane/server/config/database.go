package config

import (
	"fmt"
	"log"
	"voca-plane/internal/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database; %v", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Airline{},
		&models.Airport{},
		&models.Flight{},
		&models.FlightClass{},
		&models.FlightSeat{},
		&models.Transaction{},
		&models.TransactionPassenger{},
		&models.PromoCode{},
	)

	return db
}