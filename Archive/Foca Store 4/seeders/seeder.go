package seeders

import (
	"errors"
	"os"

	"gorm.io/gorm"

	"focastore/helper"
	"focastore/models"
)

func Seed(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	roles := []models.Role{{Name: "Admin"}, {Name: "User"}}
	for _, r := range roles {
		var existing models.Role
		err := db.Where("name = ?", r.Name).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&r).Error; err != nil {
					return err
				}
				continue
			}
			return err
		}
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	adminName := os.Getenv("ADMIN_NAME")
	if adminEmail == "" {
		adminEmail = "admin@local.dev"
	}
	if adminPass == "" {
		adminPass = "Admin123!"
	}
	if adminName == "" {
		adminName = "Admin"
	}

	var adminRole models.Role
	if err := db.Where("name = ?", "Admin").First(&adminRole).Error; err != nil {
		return err
	}

	var user models.User
	err := db.Where("email = ?", adminEmail).First(&user).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := helper.HashPassword(adminPass)
	if err != nil {
		return err
	}

	admin := models.User{Name: adminName, Email: adminEmail, PasswordHash: hash, RoleID: adminRole.ID}
	return db.Create(&admin).Error
}
