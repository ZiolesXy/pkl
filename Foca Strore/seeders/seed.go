package seeders

import (
	"errors"
	"voca-store/helper"
	"voca-store/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedRoles(db *gorm.DB) error {
	roles := []string{"Admin", "User"}
	for _, roleName := range roles {
		var existingRole models.Role
		if err := db.Where("name = ?", roleName).First(&existingRole).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				role := models.Role{Name: roleName}
				if err := db.Create(&role).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func SeedBasicRole(db *gorm.DB) error {
    roles := []models.Role{
        {Name: "Admin"},
        {Name: "User"},
    }

    // Menggunakan Clause OnConflict agar jika nama sudah ada, dia tidak error/duplikat
    // Ini jauh lebih efisien daripada melakukan loop + query manual
    return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error
}

func SeedAdmin(db *gorm.DB) error {
	var adminRole models.Role
	if err := db.Where("name = ?", "Admin").First(&adminRole).Error; err != nil {
		return err
	}

	var existingAdmin models.User
	if err := db.Where("email = ?", "admin@ecommerce.com").First(&existingAdmin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hashedPassword, err := helper.HashPassword("360589")
			if err != nil {
				return err
			}
			
			admin := models.User{
				Name:     "Pasha",
				Email:    "pashaprabasakti@gmail.com",
				Password: hashedPassword,
				RoleID:   adminRole.ID,
			}
			if err := db.Create(&admin).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func SeedUsers(db *gorm.DB) error {
	var userRole models.Role
	if err := db.Where("name = ?", "User").First(&userRole).Error; err != nil {
		return err
	}

	users := []struct {
		name  string
		email string
		pass  string
	}{
		{"John Doe", "john@example.com", "password123"},
		{"Jane Smith", "jane@example.com", "password123"},
		{"Bob Johnson", "bob@example.com", "password123"},
	}

	for _, u := range users {
		var existingUser models.User
		if err := db.Where("email = ?", u.email).First(&existingUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				hashedPassword, err := helper.HashPassword(u.pass)
				if err != nil {
					return err
				}
				
				user := models.User{
					Name:     u.name,
					Email:    u.email,
					Password: hashedPassword,
					RoleID:   userRole.ID,
				}
				if err := db.Create(&user).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func SeedProducts(db *gorm.DB) error {
	products := []struct {
		name        string
		description string
		imageURL    string
		price       float64
		stock       int
	}{
		{
			name:        "Laptop",
			description: "High performance laptop with 16GB RAM and 512GB SSD",
			imageURL:    "",
			price:       15000000,
			stock:       10,
		},
		{
			name:        "Smartphone",
			description: "Latest smartphone with 128GB storage and 6GB RAM",
			imageURL:    "",
			price:       5000000,
			stock:       25,
		},
		{
			name:        "Headphones",
			description: "Wireless noise-cancelling headphones with 30hr battery life",
			imageURL:    "",
			price:       1500000,
			stock:       50,
		},
		{
			name:        "Mouse",
			description: "Ergonomic wireless mouse with RGB lighting",
			imageURL:    "",
			price:       300000,
			stock:       100,
		},
		{
			name:        "Keyboard",
			description: "Mechanical keyboard with customizable RGB backlight",
			imageURL:    "",
			price:       800000,
			stock:       75,
		},
	}

	for _, p := range products {
		var existingProduct models.Product
		if err := db.Where("name = ?", p.name).First(&existingProduct).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				product := models.Product{
					Name:        p.name,
					Description: p.description,
					ImageURL:    p.imageURL,
					Price:       p.price,
					Stock:       p.stock,
				}
				if err := db.Create(&product).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}