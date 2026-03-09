package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type config struct {
	AppPort            string
	GinMode            string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	JWTSecret          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func GetEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func LoadConfig() *config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env using file found, using environment variables")
	}

	accessTokenExpiryStr := os.Getenv("ACCESS_TOKEN_EXPIRY")
	if accessTokenExpiryStr == "" {
		accessTokenExpiryStr = "24h"
	}
	accessTokenExpiry, _ := time.ParseDuration(accessTokenExpiryStr)

	refreshTokenExpiryStr := os.Getenv("REFRESH_TOKEN_EXPIRY")
	if refreshTokenExpiryStr == "" {
		refreshTokenExpiryStr = "168h"
	}
	refreshTokenExpiry, _ := time.ParseDuration(refreshTokenExpiryStr)

	return &config{
		AppPort:            GetEnv("APP_PORT", "8080"),
		GinMode:            GetEnv("GIN_MODE", "debug"),
		DBHost:             GetEnv("DB_HOST", "localhost"),
		DBPort:             GetEnv("DB_PORT", "5432"),
		DBUser:             GetEnv("DB_USER", "postgres"),
		DBPassword:         GetEnv("DB_PASSWORD", "postgres"),
		DBName:             GetEnv("DB_NAME", "flight_booking"),
		DBSSLMode:          GetEnv("DB_SSLMODE", "disable"),
		JWTSecret:          GetEnv("JWT_SECRET", "secret"),
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
	}
}
