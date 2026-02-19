package seeders

import (
	"errors"
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
		{"Amtir", "petir@gmail.com", "123456"},
		{"Jane Smith", "siti@gmail.com", "123456"},
		{"Bob Johnson", "bob@gmail.com", "123456"},
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
	// Ganti base64 dengan URL gambar manual atau placeholder
	// Anda bisa menggunakan link berbeda untuk tiap produk di bawah
	defaultImg := "https://picsum.photos/seed/picsum/200/300"

	products := []models.Product{
		{
			Name:        "Apple MacBook Air M2",
			Description: "Laptop tipis dengan chip Apple M2, RAM 8GB, SSD 256GB, layar Retina 13.6 inci",
			Price:       18500000,
			Stock:       10,
			ImageURL:    defaultImg,
		},
		{
			Name:        "ASUS ROG Zephyrus G14",
			Description: "Laptop gaming dengan AMD Ryzen 9, RTX 4060, RAM 16GB, SSD 1TB",
			Price:       25000000,
			Stock:       8,
			ImageURL:    defaultImg,
		},
		{
			Name:        "Samsung Galaxy S23",
			Description: "Smartphone flagship dengan Snapdragon 8 Gen 2, RAM 8GB, storage 256GB",
			Price:       12000000,
			Stock:       20,
			ImageURL:    defaultImg,
		},
		{
			Name:        "Apple iPhone 14",
			Description: "iPhone dengan chip A15 Bionic, kamera 12MP, layar Super Retina XDR",
			Price:       13500000,
			Stock:       15,
			ImageURL:    defaultImg,
		},
		{
			Name:        "Logitech G Pro X Superlight",
			Description: "Mouse gaming wireless ultra ringan dengan sensor HERO 25K",
			Price:       1800000,
			Stock:       30,
			ImageURL:    defaultImg,
		},
		{
			Name:        "Keychron K2 Mechanical Keyboard",
			Description: "Keyboard mechanical wireless dengan RGB dan hot-swappable switch",
			Price:       1400000,
			Stock:       25,
			ImageURL:    defaultImg,
		},
		{
			Name:        "LG UltraFine 27UL850 4K Monitor",
			Description: "Monitor 27 inci resolusi 4K UHD dengan HDR10 dan USB-C",
			Price:       6500000,
			Stock:       12,
			ImageURL:    defaultImg,
		},
		{
			Name:        "Sony WH-1000XM5",
			Description: "Headphone wireless dengan noise cancelling terbaik di kelasnya",
			Price:       5500000,
			Stock:       18,
			ImageURL:    defaultImg,
		},
		{
			Name:        "Samsung T7 Portable SSD 1TB",
			Description: "SSD external cepat hingga 1050MB/s dengan koneksi USB-C",
			Price:       1900000,
			Stock:       22,
			ImageURL:    defaultImg,
		},
		{
			Name:        "Xiaomi Mi Router AX3000",
			Description: "Router WiFi 6 dengan kecepatan tinggi dan latency rendah",
			Price:       900000,
			Stock:       28,
			ImageURL:    defaultImg,
		},
	}

	for _, p := range products {
		var existingProduct models.Product
		slug := helper.GenerateSlug(p.Name)

		if err := db.Where("slug = ?", slug).First(&existingProduct).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				p.Slug = slug
				// Karena menggunakan link manual, ImagePublicID dikosongkan
				// agar tidak konflik jika Anda menggunakan Cloudinary logic di tempat lain
				p.ImagePublicID = "manual_link"

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

func SeedProductsFromAssets(db *gorm.DB) error {

	assetDir := "AssetPrivate"

	files, err := os.ReadDir(assetDir)
	if err != nil {
		return err
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(assetDir, file.Name())

		// ambil nama tanpa extension
		name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		// cek apakah sudah ada
		var count int64
		db.Model(&models.Product{}).
			Where("name = ?", name).
			Count(&count)

		if count > 0 {
			fmt.Println("Product already exists:", name)
			continue
		}

		// upload ke Cloudinary
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
		}

		if err := db.Create(&product).Error; err != nil {
			fmt.Println("DB insert failed:", err)
			continue
		}

		fmt.Println("Seeded:", name)
	}

	return nil
}