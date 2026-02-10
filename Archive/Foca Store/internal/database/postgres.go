package database

import (
	"fmt"
	"log"
	"main/config"
	"main/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectPostgres(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}

	sqlDB, err := db.DB()
	if err != nil{
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(
		&models.Role{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("failed migrate", err)
	}

	DB = db
	log.Println("PostgresSQL connected")
}