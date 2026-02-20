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
		{Name: "Makanan"},
		{Name: "Minuman"},
		{Name: "Pertanian"},
		{Name: "Mainan"},
		{Name: "Lainnya"},
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

	defaultImg := "https://picsum.photos/seed/picsum/400/400"

	var laptop, smartphone, accessories, makanan, minuman, pertanian, mainan, lainnya models.Category

	db.Where("name = ?", "Laptop").First(&laptop)
	db.Where("name = ?", "Smartphone").First(&smartphone)
	db.Where("name = ?", "Accessories").First(&accessories)
	db.Where("name = ?", "Makanan").First(&makanan)
	db.Where("name = ?", "Minuman").First(&minuman)
	db.Where("name = ?", "Pertanian").First(&pertanian)
	db.Where("name = ?", "Mainan").First(&mainan)
	db.Where("name = ?", "Lainnya").First(&lainnya)

	products := []models.Product{

		// ================= LAPTOP =================

		{Name: "MacBook Air M2", Price: 18500000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "MacBook Pro M3", Price: 32000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "ASUS ROG Zephyrus G14", Price: 25000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Lenovo Legion 5", Price: 21000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Acer Aspire 5", Price: 7500000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "HP Pavilion 14", Price: 9000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Dell XPS 13", Price: 24000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "MSI Katana GF66", Price: 17000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "ASUS Vivobook 15", Price: 8000000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},
		{Name: "Lenovo IdeaPad 3", Price: 6500000, Stock: 10, ImageURL: defaultImg, CategoryID: laptop.ID},

		// ================= SMARTPHONE =================

		{Name: "iPhone 15 Pro", Price: 21000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Samsung S24 Ultra", Price: 22000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Xiaomi 14", Price: 11000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Oppo Find X6", Price: 15000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Vivo X100", Price: 14000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Realme GT5", Price: 9000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Samsung A54", Price: 6000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "iPhone 13", Price: 11000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Redmi Note 13", Price: 3500000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},
		{Name: "Infinix Zero Ultra", Price: 5000000, Stock: 10, ImageURL: defaultImg, CategoryID: smartphone.ID},

		// ================= ACCESSORIES =================

		{Name: "Logitech MX Master 3", Price: 1500000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Razer DeathAdder", Price: 800000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Keychron K2", Price: 1400000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Sony WH1000XM5", Price: 5500000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "AirPods Pro 2", Price: 3800000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Samsung SSD T7", Price: 1900000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Sandisk Flashdisk 128GB", Price: 150000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "Anker Powerbank", Price: 600000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "UGREEN USB Hub", Price: 300000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},
		{Name: "HyperX Cloud II", Price: 1200000, Stock: 10, ImageURL: defaultImg, CategoryID: accessories.ID},

		// ================= MAKANAN =================

		{Name: "Indomie Goreng", Price: 3500, Stock: 100, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Mie Sedaap", Price: 3200, Stock: 100, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Beras Ramos 5kg", Price: 75000, Stock: 50, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Chitato", Price: 10000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "SilverQueen", Price: 15000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Tango Wafer", Price: 12000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Roma Biscuit", Price: 8000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Sarden ABC", Price: 12000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Kornet Pronas", Price: 25000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},
		{Name: "Nugget Fiesta", Price: 45000, Stock: 80, ImageURL: defaultImg, CategoryID: makanan.ID},

		// ================= MINUMAN =================

		{Name: "Aqua 600ml", Price: 4000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Teh Botol Sosro", Price: 5000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Coca Cola", Price: 6000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Sprite", Price: 6000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Fanta", Price: 6000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Pocari Sweat", Price: 7000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Ultra Milk", Price: 7000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Good Day Coffee", Price: 4000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Yakult", Price: 9000, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},
		{Name: "Floridina", Price: 3500, Stock: 100, ImageURL: defaultImg, CategoryID: minuman.ID},

		// ================= PERTANIAN =================

		{Name: "Bibit Padi", Price: 50000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Bibit Jagung", Price: 40000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Pupuk Urea", Price: 120000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Pupuk Kompos", Price: 60000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Cangkul", Price: 80000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Sekop", Price: 70000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Sprayer", Price: 150000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Polybag", Price: 20000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Bibit Cabai", Price: 25000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},
		{Name: "Bibit Tomat", Price: 25000, Stock: 50, ImageURL: defaultImg, CategoryID: pertanian.ID},

		// ================= MAINAN =================

		{Name: "Lego Classic", Price: 350000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Hot Wheels", Price: 50000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Rubik 3x3", Price: 40000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Boneka Teddy", Price: 80000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "RC Car", Price: 250000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Puzzle", Price: 60000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "UNO Card", Price: 30000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Monopoly", Price: 200000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Action Figure Naruto", Price: 150000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},
		{Name: "Drone Mini", Price: 500000, Stock: 30, ImageURL: defaultImg, CategoryID: mainan.ID},

		// ================= LAINNYA =================

		{Name: "Kursi Plastik", Price: 50000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Meja Lipat", Price: 150000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Lampu LED", Price: 40000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Kipas Angin", Price: 200000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Jam Dinding", Price: 70000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Botol Minum", Price: 60000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Payung", Price: 50000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Tas Ransel", Price: 200000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Karpet", Price: 150000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
		{Name: "Bantal", Price: 80000, Stock: 30, ImageURL: defaultImg, CategoryID: lainnya.ID},
	}

	for _, p := range products {

		var existing models.Product

		slug := helper.GenerateSlug(p.Name)

		if err := db.Where("slug = ?", slug).First(&existing).Error; err == nil {
			continue
		}

		p.Slug = slug
		p.ImagePublicID = "seed"

		db.Create(&p)

	}

	return nil
}

func SeedProductsFromAssets(db *gorm.DB) error {

	assetDir := "AssetPrivate"

	// ambil semua category dari database
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		return err
	}

	// buat map: lowercase(name) -> category
	categoryMap := make(map[string]models.Category)
	for _, c := range categories {
		categoryMap[strings.ToLower(c.Name)] = c
	}

	// baca folder utama
	dirs, err := os.ReadDir(assetDir)
	if err != nil {
		return err
	}

	for _, dir := range dirs {

		if !dir.IsDir() {
			continue
		}

		folderName := strings.ToLower(dir.Name())

		category, ok := categoryMap[folderName]
		if !ok {
			fmt.Println("Category tidak ditemukan di DB:", folderName)
			continue
		}

		categoryPath := filepath.Join(assetDir, dir.Name())

		files, err := os.ReadDir(categoryPath)
		if err != nil {
			fmt.Println("Gagal baca folder:", err)
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
				fmt.Println("Upload gagal:", err)
				continue
			}

			product := models.Product{
				Name:          name,
				Slug:          helper.GenerateSlug(name),
				Description:   name,
				ImageURL:      uploadResult.SecureURL,
				ImagePublicID: uploadResult.PublicID,
				Price:         1000000,
				Stock:         100,
				CategoryID:    category.ID,
			}

			if err := db.Create(&product).Error; err != nil {
				fmt.Println("Insert gagal:", err)
				continue
			}

			fmt.Println("Seeded:", name, "| Category:", category.Name)
		}
	}

	return nil
}

func SyncAssetProductsWithDefaultSeed(db *gorm.DB) error {

	// default data dari SeedProducts (tanpa ImageURL & CategoryID)
	defaultProducts := []models.Product{

		// ================= LAPTOP =================
		{Name: "MacBook Air M2", Price: 18500000, Stock: 10, Description: "MacBook Air M2"},
		{Name: "MacBook Pro M3", Price: 32000000, Stock: 10, Description: "MacBook Pro M3"},
		{Name: "ASUS ROG Zephyrus G14", Price: 25000000, Stock: 10, Description: "ASUS ROG Zephyrus G14"},
		{Name: "Lenovo Legion 5", Price: 21000000, Stock: 10, Description: "Lenovo Legion 5"},
		{Name: "Acer Aspire 5", Price: 7500000, Stock: 10, Description: "Acer Aspire 5"},
		{Name: "HP Pavilion 14", Price: 9000000, Stock: 10, Description: "HP Pavilion 14"},
		{Name: "Dell XPS 13", Price: 24000000, Stock: 10, Description: "Dell XPS 13"},
		{Name: "MSI Katana GF66", Price: 17000000, Stock: 10, Description: "MSI Katana GF66"},
		{Name: "ASUS Vivobook 15", Price: 8000000, Stock: 10, Description: "ASUS Vivobook 15"},
		{Name: "Lenovo IdeaPad 3", Price: 6500000, Stock: 10, Description: "Lenovo IdeaPad 3"},

		// ================= SMARTPHONE =================
		{Name: "iPhone 15 Pro", Price: 21000000, Stock: 10, Description: "iPhone 15 Pro"},
		{Name: "Samsung S24 Ultra", Price: 22000000, Stock: 10, Description: "Samsung S24 Ultra"},
		{Name: "Xiaomi 14", Price: 11000000, Stock: 10, Description: "Xiaomi 14"},
		{Name: "Oppo Find X6", Price: 15000000, Stock: 10, Description: "Oppo Find X6"},
		{Name: "Vivo X100", Price: 14000000, Stock: 10, Description: "Vivo X100"},
		{Name: "Realme GT5", Price: 9000000, Stock: 10, Description: "Realme GT5"},
		{Name: "Samsung A54", Price: 6000000, Stock: 10, Description: "Samsung A54"},
		{Name: "iPhone 13", Price: 11000000, Stock: 10, Description: "iPhone 13"},
		{Name: "Redmi Note 13", Price: 3500000, Stock: 10, Description: "Redmi Note 13"},
		{Name: "Infinix Zero Ultra", Price: 5000000, Stock: 10, Description: "Infinix Zero Ultra"},

		// ================= ACCESSORIES =================
		{Name: "Logitech MX Master 3", Price: 1500000, Stock: 10, Description: "Logitech MX Master 3"},
		{Name: "Razer DeathAdder", Price: 800000, Stock: 10, Description: "Razer DeathAdder"},
		{Name: "Keychron K2", Price: 1400000, Stock: 10, Description: "Keychron K2"},
		{Name: "Sony WH1000XM5", Price: 5500000, Stock: 10, Description: "Sony WH1000XM5"},
		{Name: "AirPods Pro 2", Price: 3800000, Stock: 10, Description: "AirPods Pro 2"},
		{Name: "Samsung SSD T7", Price: 1900000, Stock: 10, Description: "Samsung SSD T7"},
		{Name: "Sandisk Flashdisk 128GB", Price: 150000, Stock: 10, Description: "Sandisk Flashdisk 128GB"},
		{Name: "Anker Powerbank", Price: 600000, Stock: 10, Description: "Anker Powerbank"},
		{Name: "UGREEN USB Hub", Price: 300000, Stock: 10, Description: "UGREEN USB Hub"},
		{Name: "HyperX Cloud II", Price: 1200000, Stock: 10, Description: "HyperX Cloud II"},

		// ================= MAKANAN =================
		{Name: "Indomie Goreng", Price: 3500, Stock: 100, Description: "Indomie Goreng"},
		{Name: "Mie Sedaap", Price: 3200, Stock: 100, Description: "Mie Sedaap"},
		{Name: "Beras Ramos 5kg", Price: 75000, Stock: 50, Description: "Beras Ramos 5kg"},
		{Name: "Chitato", Price: 10000, Stock: 80, Description: "Chitato"},
		{Name: "SilverQueen", Price: 15000, Stock: 80, Description: "SilverQueen"},
		{Name: "Tango Wafer", Price: 12000, Stock: 80, Description: "Tango Wafer"},
		{Name: "Roma Biscuit", Price: 8000, Stock: 80, Description: "Roma Biscuit"},
		{Name: "Sarden ABC", Price: 12000, Stock: 80, Description: "Sarden ABC"},
		{Name: "Kornet Pronas", Price: 25000, Stock: 80, Description: "Kornet Pronas"},
		{Name: "Nugget Fiesta", Price: 45000, Stock: 80, Description: "Nugget Fiesta"},

		// ================= MINUMAN =================
		{Name: "Aqua 600ml", Price: 4000, Stock: 100, Description: "Aqua 600ml"},
		{Name: "Teh Botol Sosro", Price: 5000, Stock: 100, Description: "Teh Botol Sosro"},
		{Name: "Coca Cola", Price: 6000, Stock: 100, Description: "Coca Cola"},
		{Name: "Sprite", Price: 6000, Stock: 100, Description: "Sprite"},
		{Name: "Fanta", Price: 6000, Stock: 100, Description: "Fanta"},
		{Name: "Pocari Sweat", Price: 7000, Stock: 100, Description: "Pocari Sweat"},
		{Name: "Ultra Milk", Price: 7000, Stock: 100, Description: "Ultra Milk"},
		{Name: "Good Day Coffee", Price: 4000, Stock: 100, Description: "Good Day Coffee"},
		{Name: "Yakult", Price: 9000, Stock: 100, Description: "Yakult"},
		{Name: "Floridina", Price: 5000, Stock: 100, Description: "Floridina"},

		// ================= PERTANIAN =================
		{Name: "Bibit Padi", Price: 50000, Stock: 50, Description: "Bibit Padi"},
		{Name: "Bibit Jagung", Price: 40000, Stock: 50, Description: "Bibit Jagung"},
		{Name: "Pupuk Urea", Price: 120000, Stock: 50, Description: "Pupuk Urea"},
		{Name: "Pupuk Kompos", Price: 60000, Stock: 50, Description: "Pupuk Kompos"},
		{Name: "Cangkul", Price: 80000, Stock: 50, Description: "Cangkul"},
		{Name: "Sekop", Price: 70000, Stock: 50, Description: "Sekop"},
		{Name: "Sprayer", Price: 150000, Stock: 50, Description: "Sprayer"},
		{Name: "Polybag", Price: 20000, Stock: 50, Description: "Polybag"},
		{Name: "Bibit Cabai", Price: 25000, Stock: 50, Description: "Bibit Cabai"},
		{Name: "Bibit Tomat", Price: 25000, Stock: 50, Description: "Bibit Tomat"},

		// ================= MAINAN =================
		{Name: "Lego Classic", Price: 350000, Stock: 30, Description: "Lego Classic"},
		{Name: "Hot Wheels", Price: 50000, Stock: 30, Description: "Hot Wheels"},
		{Name: "Rubik 3x3", Price: 40000, Stock: 30, Description: "Rubik 3x3"},
		{Name: "Boneka Teddy", Price: 80000, Stock: 30, Description: "Boneka Teddy"},
		{Name: "RC Car", Price: 250000, Stock: 30, Description: "RC Car"},
		{Name: "Puzzle", Price: 60000, Stock: 30, Description: "Puzzle"},
		{Name: "UNO Card", Price: 30000, Stock: 30, Description: "UNO Card"},
		{Name: "Monopoly", Price: 200000, Stock: 30, Description: "Monopoly"},
		{Name: "Action Figure Naruto", Price: 150000, Stock: 30, Description: "Action Figure Naruto"},
		{Name: "Drone Mini", Price: 500000, Stock: 30, Description: "Drone Mini"},

		// ================= LAINNYA =================
		{Name: "Kursi Plastik", Price: 50000, Stock: 30, Description: "Kursi Plastik"},
		{Name: "Meja Lipat", Price: 150000, Stock: 30, Description: "Meja Lipat"},
		{Name: "Lampu LED", Price: 40000, Stock: 30, Description: "Lampu LED"},
		{Name: "Kipas Angin", Price: 200000, Stock: 30, Description: "Kipas Angin"},
		{Name: "Jam Dinding", Price: 70000, Stock: 30, Description: "Jam Dinding"},
		{Name: "Botol Minum", Price: 60000, Stock: 30, Description: "Botol Minum"},
		{Name: "Payung", Price: 50000, Stock: 30, Description: "Payung"},
		{Name: "Tas Ransel", Price: 200000, Stock: 30, Description: "Tas Ransel"},
		{Name: "Karpet", Price: 150000, Stock: 30, Description: "Karpet"},
		{Name: "Bantal", Price: 80000, Stock: 30, Description: "Bantal"},
	}

	// buat map name -> default product
	defaultMap := make(map[string]models.Product)
	for _, p := range defaultProducts {
		defaultMap[p.Name] = p
	}

	// ambil semua product di DB
	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}

	for _, product := range products {

		defaultData, exists := defaultMap[product.Name]
		if !exists {
			continue
		}

		updates := map[string]interface{}{
			"description": defaultData.Description,
			"price":       defaultData.Price,
			"stock":       defaultData.Stock,
		}

		if err := db.Model(&product).Updates(updates).Error; err != nil {
			fmt.Println("Gagal update:", product.Name)
			continue
		}

		fmt.Println("Updated:", product.Name)
	}

	return nil
}

func SeedCoupons(db *gorm.DB) error {

	coupons := []models.Coupon{
		// percentage
		{
			Code:      "DISC10",
			Type:      "percentage",
			Value:     10,
			Quota:     100,
			UsedCount: 0,
		},

		{
			Code:      "DISC20",
			Type:      "percentage",
			Value:     20,
			Quota:     50,
			UsedCount: 0,
		},

		{
			Code:      "DISC30",
			Type:      "percentage",
			Value:     30,
			Quota:     30,
			UsedCount: 0,
		},

		{
			Code:      "WELCOME",
			Type:      "percentage",
			Value:     15,
			Quota:     200,
			UsedCount: 0,
		},

		{
			Code:      "FLASH50",
			Type:      "percentage",
			Value:     50,
			Quota:     10,
			UsedCount: 0,
		},

		// fixed
		{
			Code:      "FIXED10K",
			Type:      "fixed",
			Value:     10000,
			Quota:     100,
			UsedCount: 0,
		},

		{
			Code:      "FIXED25K",
			Type:      "fixed",
			Value:     25000,
			Quota:     100,
			UsedCount: 0,
		},

		{
			Code:      "FIXED50K",
			Type:      "fixed",
			Value:     50000,
			Quota:     50,
			UsedCount: 0,
		},

		{
			Code:      "FIXED100K",
			Type:      "fixed",
			Value:     100000,
			Quota:     20,
			UsedCount: 0,
		},

		{
			Code:      "BIGSALE",
			Type:      "fixed",
			Value:     150000,
			Quota:     10,
			UsedCount: 0,
		},

		{
			Code:      "LIMITED",
			Type:      "fixed",
			Value:     200000,
			Quota:     5,
			UsedCount: 0,
		},

		{
			Code:      "ONGKIR",
			Type:      "fixed",
			Value:     15000,
			Quota:     300,
			UsedCount: 0,
		},
	}

	return db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoNothing: true,
		}).
		Create(&coupons).Error
}