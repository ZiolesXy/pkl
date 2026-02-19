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

		// ================= LAPTOP (12) =================
		{Name: "Apple MacBook Air M2", Description: "Chip Apple M2, 13-inch", Price: 18500000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Apple MacBook Pro M3 14-inch", Description: "Chip Apple M3 terbaru", Price: 32000000, Stock: 6, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "ASUS ROG Zephyrus G14", Description: "Ryzen 9 RTX Series", Price: 25000000, Stock: 8, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "ASUS TUF Gaming F15", Description: "Intel i7 RTX 4060", Price: 19000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Lenovo Legion 5 Pro", Description: "Ryzen 7 RTX 4070", Price: 27000000, Stock: 5, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Lenovo IdeaPad Slim 5", Description: "Laptop tipis Ryzen 5", Price: 9500000, Stock: 15, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "HP Spectre x360", Description: "Laptop premium convertible", Price: 22000000, Stock: 7, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "HP Pavilion Gaming 15", Description: "Gaming entry RTX", Price: 15000000, Stock: 11, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Dell XPS 13", Description: "Ultrabook premium", Price: 24000000, Stock: 9, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Dell Inspiron 14", Description: "Laptop kerja harian", Price: 8500000, Stock: 14, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Acer Predator Helios 300", Description: "Gaming RTX Series", Price: 21000000, Stock: 6, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Acer Aspire 5", Description: "Laptop kerja ringan", Price: 7500000, Stock: 18, ImageURL: defaultImg, CategoryID: laptop.ID},

		// ================= SMARTPHONE (12) =================
		{Name: "Apple iPhone 15 Pro", Description: "Chip A17 Pro", Price: 21000000, Stock: 12, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Apple iPhone 14", Description: "Chip A15 Bionic", Price: 13500000, Stock: 15, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Samsung Galaxy S24 Ultra", Description: "Snapdragon 8 Gen 3", Price: 22000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Samsung Galaxy S23 FE", Description: "Flagship killer", Price: 9000000, Stock: 18, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Xiaomi 14", Description: "Leica camera flagship", Price: 11000000, Stock: 20, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Xiaomi Redmi Note 13 Pro", Description: "Midrange AMOLED", Price: 4500000, Stock: 25, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "OPPO Find X6 Pro", Description: "Kamera flagship", Price: 15000000, Stock: 8, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "OPPO Reno11", Description: "Stylish midrange", Price: 5500000, Stock: 17, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Vivo X100 Pro", Description: "Zeiss camera system", Price: 14000000, Stock: 9, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Realme GT 5 Pro", Description: "Snapdragon flagship", Price: 10000000, Stock: 13, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Google Pixel 8", Description: "Tensor G3 chip", Price: 12000000, Stock: 7, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Nothing Phone 2", Description: "Glyph interface", Price: 9500000, Stock: 11, ImageURL: defaultImg, CategoryID: smartphone.ID},

		// ================= ACCESSORIES (14) =================
		{Name: "Logitech MX Master 3S", Description: "Mouse productivity", Price: 1500000, Stock: 30, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Logitech G Pro X Superlight", Description: "Mouse gaming wireless", Price: 1800000, Stock: 30, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Razer BlackWidow V4", Description: "Mechanical gaming keyboard", Price: 2500000, Stock: 20, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Keychron K2", Description: "Mechanical wireless keyboard", Price: 1400000, Stock: 25, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Sony WH-1000XM5", Description: "Noise cancelling headphone", Price: 5500000, Stock: 18, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Apple AirPods Pro 2", Description: "ANC true wireless", Price: 3800000, Stock: 22, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Samsung T7 SSD 1TB", Description: "Portable SSD", Price: 1900000, Stock: 22, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Seagate Expansion 2TB", Description: "External HDD", Price: 1200000, Stock: 26, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "LG UltraFine 27UL850", Description: "4K IPS Monitor", Price: 6500000, Stock: 12, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Dell P2723QE", Description: "4K Office Monitor", Price: 7200000, Stock: 9, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Anker PowerCore 20000", Description: "Powerbank fast charging", Price: 800000, Stock: 35, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "UGREEN USB-C Hub 7in1", Description: "USB-C multiport hub", Price: 600000, Stock: 28, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "SteelSeries Arctis 7", Description: "Wireless gaming headset", Price: 2200000, Stock: 14, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "HyperX Cloud II", Description: "Gaming headset", Price: 1500000, Stock: 19, ImageURL: defaultImg, CategoryID: accessories.ID},

		// ================= NETWORKING (12) =================
		{Name: "TP-Link Archer AX73", Description: "WiFi 6 Router", Price: 1800000, Stock: 20, ImageURL: defaultImg, CategoryID: networking.ID},
		{Name: "TP-Link Deco X20", Description: "Mesh WiFi 6", Price: 2500000, Stock: 15, ImageURL: defaultImg, CategoryID: networking.ID},
		{Name: "ASUS RT-AX88U", Description: "Gaming Router WiFi 6", Price: 4200000, Stock: 10, ImageURL: defaultImg, CategoryID: networking.ID},
		{Name: "Xiaomi Router AX3000", Description: "WiFi 6 Router", Price: 900000, Stock: 28, ImageURL: defaultImg, CategoryID: networking.ID},
		{Name: "MikroTik hAP ac3", Description: "Advanced router", Price: 1700000, Stock: 18, ImageURL: defaultImg, CategoryID: networking.ID},
		{Name: "Netgear Nighthawk AX12", Description: "High performance router", Price: 5500000, Stock: 6, ImageURL: defaultImg, CategoryID: networking.ID},
		{Name: "Ubiquiti UniFi AP AC Lite", Description: "Access point enterprise", Price: 1600000, Stock: 14, ImageURL: defaultImg, CategoryID: networking.ID},
		{Name: "D-Link DIR-X5460", Description: "WiFi 6 Router", Price: 2300000, Stock: 12, ImageURL: defaultImg, CategoryID: networking.ID},
		{Name: "Tenda AC23", Description: "Dual band router", Price: 650000, Stock: 30, ImageURL: defaultImg, CategoryID: networking.ID},
		{Name: "Huawei AX3 Pro", Description: "WiFi 6 router", Price: 1200000, Stock: 17, ImageURL: defaultImg, CategoryID: networking.ID},
		{Name: "Cisco RV340", Description: "Business VPN router", Price: 4800000, Stock: 5, ImageURL: defaultImg, CategoryID: networking.ID},
		{Name: "Linksys MR9600", Description: "Mesh WiFi 6 router", Price: 3900000, Stock: 8, ImageURL: defaultImg, CategoryID: networking.ID},
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

	// category map
	validCategories := map[string]string{
		"laptop":      "Laptop",
		"smartphone":  "Smartphone",
		"accessories": "Accessories",
		"networking":  "Networking",
	}

	// loop folder category
	dirs, err := os.ReadDir(assetDir)
	if err != nil {
		return err
	}

	for _, dir := range dirs {

		if !dir.IsDir() {
			continue
		}

		folderName := strings.ToLower(dir.Name())

		categoryName, ok := validCategories[folderName]
		if !ok {
			fmt.Println("Category not allowed:", folderName)
			continue
		}

		// find category in DB
		var category models.Category
		if err := db.Where("name = ?", categoryName).First(&category).Error; err != nil {
			fmt.Println("Category not found:", categoryName)
			continue
		}

		categoryPath := filepath.Join(assetDir, folderName)

		files, err := os.ReadDir(categoryPath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, file := range files {

			if file.IsDir() {
				continue
			}

			filePath := filepath.Join(categoryPath, file.Name())

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
				Stock:         3200,
				CategoryID:    category.ID,
			}

			if err := db.Create(&product).Error; err != nil {

				fmt.Println("Insert failed:", err)

				continue

			}

			fmt.Println("Seeded:", name, "Category:", categoryName)

		}

	}

	return nil
}