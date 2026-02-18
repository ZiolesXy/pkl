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
	// Gambar base64 yang Anda berikan
	imageBase64 := "data:image/jpeg;base64,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"

	products := []models.Product{
		{Name: "Laptop High-End", Description: "Performance tinggi dengan RAM 16GB dan SSD 512GB", Price: 15000000, Stock: 10, ImageURL: imageBase64},
		{Name: "Smartphone Pro", Description: "Penyimpanan 128GB dengan layar AMOLED 120Hz", Price: 5000000, Stock: 25, ImageURL: imageBase64},
		{Name: "Wireless Headphones", Description: "Noise-cancelling dengan daya tahan baterai 30 jam", Price: 1500000, Stock: 50, ImageURL: imageBase64},
		{Name: "Gaming Mouse RGB", Description: "Mouse ergonomis dengan sensor optik presisi tinggi", Price: 300000, Stock: 100, ImageURL: imageBase64},
		{Name: "Mechanical Keyboard", Description: "Keyboard tactile dengan backlight RGB kustom", Price: 800000, Stock: 75, ImageURL: imageBase64},
		{Name: "Monitor 4K", Description: "Layar tajam 27 inci untuk kebutuhan desain grafis", Price: 4500000, Stock: 15, ImageURL: imageBase64},
		{Name: "Webcam Full HD", Description: "Kamera jernih untuk meeting dan streaming", Price: 600000, Stock: 40, ImageURL: imageBase64},
		{Name: "External Hard Drive", Description: "Kapasitas 1TB untuk backup data aman Anda", Price: 900000, Stock: 30, ImageURL: imageBase64},
		{Name: "Power Bank 20k", Description: "Kapasitas besar 20.000mAh dengan fast charging", Price: 400000, Stock: 60, ImageURL: imageBase64},
		{Name: "USB-C Hub Multiport", Description: "Adaptor serbaguna untuk konektivitas maksimal", Price: 350000, Stock: 85, ImageURL: imageBase64},
	}

	for _, p := range products {
		var existingProduct models.Product
		// Cek berdasarkan Nama untuk menghindari duplikasi saat seeding ulang
		if err := db.Where("name = ?", p.Name).First(&existingProduct).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&p).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}