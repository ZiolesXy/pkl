package database

import (
	"fmt"
	"main/internal/model"
	"main/pkg/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	host := config.GetEnv("DB_HOST", "localhost")
	user := config.GetEnv("DB_USER", "postgres")
	password := config.GetEnv("DB_PASSWORD", "")
	dbName := config.GetEnv("DB_NAME", "go_ecommerce_db")
	port := config.GetEnv("DB_PORT", "5432")
	sslMode := config.GetEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host, user, password, dbName, port, sslMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("Database connection established")
	DB = db
	
	MigrateDB(db)
}

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{}, &model.Category{}, &model.Product{})
	if err != nil {
		log.Fatal("Migration Failed: ", err)
	}
}