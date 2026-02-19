package seeders

import (
	// "errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"voca-store/helper"
	"voca-store/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedRoles(db *gorm.DB) error {
	roles := []models.Role{
		{Name: "Admin"},
		{Name: "User"},
	}

	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error
}

func SeedCategories(db *gorm.DB) error {
	categories := []models.Category{
		{Name: "Laptop"},
		{Name: "Smartphone"},
		{Name: "Accessories"},
		{Name: "Networking"},
	}

	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&categories).Error
}

func SeedAdmin(db *gorm.DB) error {
	var adminRole models.Role
	if err := db.Where("name = ?", "Admin").First(&adminRole).Error; err != nil {
		return err
	}

	var existingAdmin models.User
	if err := db.Where("email = ?", "pashaprabasakti@gmail.com").First(&existingAdmin).Error; err == nil {
		return nil // already exists
	}

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

	return db.Create(&admin).Error
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
		{"Amtir", "petir@gmail.com", "123456"},
		{"Jane Smith", "siti@gmail.com", "123456"},
		{"Bob Johnson", "bob@gmail.com", "123456"},
	}

	for _, u := range users {
		var existingUser models.User
		if err := db.Where("email = ?", u.email).First(&existingUser).Error; err == nil {
			continue
		}

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
	}

	return nil
}

func SeedProducts(db *gorm.DB) error {
	defaultImg := "https://picsum.photos/seed/picsum/200/300"

	// Ambil category dari DB
	var laptop, smartphone, accessories, networking models.Category

	if err := db.Where("name = ?", "Laptop").First(&laptop).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "Smartphone").First(&smartphone).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "Accessories").First(&accessories).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "Networking").First(&networking).Error; err != nil {
		return err
	}

	products := []models.Product{
		{
			Name:        "Apple MacBook Air M2",
			Description: "Laptop tipis dengan chip Apple M2",
			Price:       18500000,
			Stock:       10,
			ImageURL:    defaultImg,
			CategoryID:  laptop.ID,
		},
		{
			Name:        "ASUS ROG Zephyrus G14",
			Description: "Laptop gaming AMD Ryzen 9",
			Price:       25000000,
			Stock:       8,
			ImageURL:    defaultImg,
			CategoryID:  laptop.ID,
		},
		{
			Name:        "Samsung Galaxy S23",
			Description: "Flagship Snapdragon 8 Gen 2",
			Price:       12000000,
			Stock:       20,
			ImageURL:    defaultImg,
			CategoryID:  smartphone.ID,
		},
		{
			Name:        "Apple iPhone 14",
			Description: "iPhone chip A15 Bionic",
			Price:       13500000,
			Stock:       15,
			ImageURL:    defaultImg,
			CategoryID:  smartphone.ID,
		},
		{
			Name:        "Logitech G Pro X Superlight",
			Description: "Mouse gaming wireless",
			Price:       1800000,
			Stock:       30,
			ImageURL:    defaultImg,
			CategoryID:  accessories.ID,
		},
		{
			Name:        "Keychron K2 Mechanical Keyboard",
			Description: "Keyboard mechanical wireless",
			Price:       1400000,
			Stock:       25,
			ImageURL:    defaultImg,
			CategoryID:  accessories.ID,
		},
		{
			Name:        "LG UltraFine 27UL850 4K Monitor",
			Description: "Monitor 27 inch 4K",
			Price:       6500000,
			Stock:       12,
			ImageURL:    defaultImg,
			CategoryID:  accessories.ID,
		},
		{
			Name:        "Sony WH-1000XM5",
			Description: "Headphone noise cancelling",
			Price:       5500000,
			Stock:       18,
			ImageURL:    defaultImg,
			CategoryID:  accessories.ID,
		},
		{
			Name:        "Samsung T7 Portable SSD 1TB",
			Description: "SSD external cepat",
			Price:       1900000,
			Stock:       22,
			ImageURL:    defaultImg,
			CategoryID:  accessories.ID,
		},
		{
			Name:        "Xiaomi Mi Router AX3000",
			Description: "Router WiFi 6",
			Price:       900000,
			Stock:       28,
			ImageURL:    defaultImg,
			CategoryID:  networking.ID,
		},
	}

	for _, p := range products {
		var existing models.Product
		slug := helper.GenerateSlug(p.Name)

		if err := db.Where("slug = ?", slug).First(&existing).Error; err == nil {
			continue
		}

		p.Slug = slug
		p.ImagePublicID = "manual_link"

		if err := db.Create(&p).Error; err != nil {
			return err
		}
	}

	return nil
}

func SeedProductsFromAssets(db *gorm.DB) error {

	assetDir := "AssetPrivate"

	files, err := os.ReadDir(assetDir)
	if err != nil {
		return err
	}

	// default category = Accessories
	var defaultCategory models.Category
	if err := db.Where("name = ?", "Accessories").First(&defaultCategory).Error; err != nil {
		return err
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(assetDir, file.Name())
		name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		var count int64
		db.Model(&models.Product{}).
			Where("name = ?", name).
			Count(&count)

		if count > 0 {
			continue
		}

		uploadResult, err := helper.UploadImageFromFile(filePath, "products")
		if err != nil {
			fmt.Println("Upload failed:", err)
			continue
		}

		product := models.Product{
			Name:          name,
			Slug:          helper.GenerateSlug(name),
			Description:   name,
			ImageURL:      uploadResult.SecureURL,
			ImagePublicID: uploadResult.PublicID,
			Price:         1000000,
			Stock:         10,
			CategoryID:    defaultCategory.ID,
		}

		if err := db.Create(&product).Error; err != nil {
			fmt.Println("DB insert failed:", err)
			continue
		}

		fmt.Println("Seeded:", name)
	}

	return nil
}
