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

	// products := []models.Product{
	// 	{Name: "Laptop High-End", Description: "Performance tinggi dengan RAM 16GB dan SSD 512GB", Price: 15000000, Stock: 10, ImageURL: ""},
	// 	{Name: "Smartphone Pro", Description: "Penyimpanan 128GB dengan layar AMOLED 120Hz", Price: 5000000, Stock: 25, ImageURL: ""},
	// 	{Name: "Wireless Headphones", Description: "Noise-cancelling dengan daya tahan baterai 30 jam", Price: 1500000, Stock: 50, ImageURL: ""},
	// 	{Name: "Gaming Mouse RGB", Description: "Mouse ergonomis dengan sensor optik presisi tinggi", Price: 300000, Stock: 100, ImageURL: defaultImg},
	// 	{Name: "Mechanical Keyboard", Description: "Keyboard tactile dengan backlight RGB kustom", Price: 800000, Stock: 75, ImageURL: ""},
	// 	{Name: "Monitor 4K", Description: "Layar tajam 27 inci untuk kebutuhan desain grafis", Price: 4500000, Stock: 15, ImageURL: ""},
	// 	{Name: "Webcam Full HD", Description: "Kamera jernih untuk meeting dan streaming", Price: 600000, Stock: 40, ImageURL: ""},
	// 	{Name: "External Hard Drive", Description: "Kapasitas 1TB untuk backup data aman Anda", Price: 900000, Stock: 30, ImageURL: ""},
	// 	{Name: "Power Bank 20k", Description: "Kapasitas besar 20.000mAh dengan fast charging", Price: 400000, Stock: 60, ImageURL: ""},
	// 	{Name: "USB-C Hub Multiport", Description: "Adaptor serbaguna untuk konektivitas maksimal", Price: 350000, Stock: 85, ImageURL: ""},
	// 	{Name: "Bluetooth Speaker Mini", Description: "Speaker portable dengan suara jernih", Price: 700000, Stock: 45, ImageURL: ""},
	// 	{Name: "Tablet 10 Inch", Description: "Tablet ringan untuk hiburan dan produktivitas", Price: 3500000, Stock: 20, ImageURL: ""},
	// 	{Name: "Smartwatch Series X", Description: "Jam pintar dengan fitur kesehatan lengkap", Price: 2500000, Stock: 35, ImageURL: ""},
	// 	{Name: "Router WiFi 6", Description: "Kecepatan tinggi untuk kebutuhan internet rumah", Price: 1200000, Stock: 28, ImageURL: ""},
	// 	{Name: "Portable SSD 1TB", Description: "Penyimpanan eksternal super cepat", Price: 1800000, Stock: 22, ImageURL: ""},
	// }

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
