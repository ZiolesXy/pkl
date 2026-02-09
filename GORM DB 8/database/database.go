package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnedtDB() {
	dsn := "host=172.16.17.123 user=postgres password=360589 dbname=gorm_auth_4 port=5432"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal Koneksi ke database")
	}

	DB = db
	fmt.Println("Datatbase Terkoneksi")
}