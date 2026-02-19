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

{Name:"MacBook Air M2",Price:18500000,Stock:10,ImageURL:defaultImg,CategoryID:laptop.ID},
{Name:"MacBook Pro M3",Price:32000000,Stock:10,ImageURL:defaultImg,CategoryID:laptop.ID},
{Name:"ASUS ROG Zephyrus G14",Price:25000000,Stock:10,ImageURL:defaultImg,CategoryID:laptop.ID},
{Name:"Lenovo Legion 5",Price:21000000,Stock:10,ImageURL:defaultImg,CategoryID:laptop.ID},
{Name:"Acer Aspire 5",Price:7500000,Stock:10,ImageURL:defaultImg,CategoryID:laptop.ID},
{Name:"HP Pavilion 14",Price:9000000,Stock:10,ImageURL:defaultImg,CategoryID:laptop.ID},
{Name:"Dell XPS 13",Price:24000000,Stock:10,ImageURL:defaultImg,CategoryID:laptop.ID},
{Name:"MSI Katana GF66",Price:17000000,Stock:10,ImageURL:defaultImg,CategoryID:laptop.ID},
{Name:"ASUS Vivobook 15",Price:8000000,Stock:10,ImageURL:defaultImg,CategoryID:laptop.ID},
{Name:"Lenovo IdeaPad 3",Price:6500000,Stock:10,ImageURL:defaultImg,CategoryID:laptop.ID},

// ================= SMARTPHONE =================

{Name:"iPhone 15 Pro",Price:21000000,Stock:10,ImageURL:defaultImg,CategoryID:smartphone.ID},
{Name:"Samsung S24 Ultra",Price:22000000,Stock:10,ImageURL:defaultImg,CategoryID:smartphone.ID},
{Name:"Xiaomi 14",Price:11000000,Stock:10,ImageURL:defaultImg,CategoryID:smartphone.ID},
{Name:"Oppo Find X6",Price:15000000,Stock:10,ImageURL:defaultImg,CategoryID:smartphone.ID},
{Name:"Vivo X100",Price:14000000,Stock:10,ImageURL:defaultImg,CategoryID:smartphone.ID},
{Name:"Realme GT5",Price:9000000,Stock:10,ImageURL:defaultImg,CategoryID:smartphone.ID},
{Name:"Samsung A54",Price:6000000,Stock:10,ImageURL:defaultImg,CategoryID:smartphone.ID},
{Name:"iPhone 13",Price:11000000,Stock:10,ImageURL:defaultImg,CategoryID:smartphone.ID},
{Name:"Redmi Note 13",Price:3500000,Stock:10,ImageURL:defaultImg,CategoryID:smartphone.ID},
{Name:"Infinix Zero Ultra",Price:5000000,Stock:10,ImageURL:defaultImg,CategoryID:smartphone.ID},

// ================= ACCESSORIES =================

{Name:"Logitech MX Master 3",Price:1500000,Stock:10,ImageURL:defaultImg,CategoryID:accessories.ID},
{Name:"Razer DeathAdder",Price:800000,Stock:10,ImageURL:defaultImg,CategoryID:accessories.ID},
{Name:"Keychron K2",Price:1400000,Stock:10,ImageURL:defaultImg,CategoryID:accessories.ID},
{Name:"Sony WH1000XM5",Price:5500000,Stock:10,ImageURL:defaultImg,CategoryID:accessories.ID},
{Name:"AirPods Pro 2",Price:3800000,Stock:10,ImageURL:defaultImg,CategoryID:accessories.ID},
{Name:"Samsung SSD T7",Price:1900000,Stock:10,ImageURL:defaultImg,CategoryID:accessories.ID},
{Name:"Sandisk Flashdisk 128GB",Price:150000,Stock:10,ImageURL:defaultImg,CategoryID:accessories.ID},
{Name:"Anker Powerbank",Price:600000,Stock:10,ImageURL:defaultImg,CategoryID:accessories.ID},
{Name:"UGREEN USB Hub",Price:300000,Stock:10,ImageURL:defaultImg,CategoryID:accessories.ID},
{Name:"HyperX Cloud II",Price:1200000,Stock:10,ImageURL:defaultImg,CategoryID:accessories.ID},

// ================= MAKANAN =================

{Name:"Indomie Goreng",Price:3500,Stock:100,ImageURL:defaultImg,CategoryID:makanan.ID},
{Name:"Mie Sedaap",Price:3200,Stock:100,ImageURL:defaultImg,CategoryID:makanan.ID},
{Name:"Beras Ramos 5kg",Price:75000,Stock:50,ImageURL:defaultImg,CategoryID:makanan.ID},
{Name:"Chitato",Price:10000,Stock:80,ImageURL:defaultImg,CategoryID:makanan.ID},
{Name:"SilverQueen",Price:15000,Stock:80,ImageURL:defaultImg,CategoryID:makanan.ID},
{Name:"Tango Wafer",Price:12000,Stock:80,ImageURL:defaultImg,CategoryID:makanan.ID},
{Name:"Roma Biscuit",Price:8000,Stock:80,ImageURL:defaultImg,CategoryID:makanan.ID},
{Name:"Sarden ABC",Price:12000,Stock:80,ImageURL:defaultImg,CategoryID:makanan.ID},
{Name:"Kornet Pronas",Price:25000,Stock:80,ImageURL:defaultImg,CategoryID:makanan.ID},
{Name:"Nugget Fiesta",Price:45000,Stock:80,ImageURL:defaultImg,CategoryID:makanan.ID},

// ================= MINUMAN =================

{Name:"Aqua 600ml",Price:4000,Stock:100,ImageURL:defaultImg,CategoryID:minuman.ID},
{Name:"Teh Botol Sosro",Price:5000,Stock:100,ImageURL:defaultImg,CategoryID:minuman.ID},
{Name:"Coca Cola",Price:6000,Stock:100,ImageURL:defaultImg,CategoryID:minuman.ID},
{Name:"Sprite",Price:6000,Stock:100,ImageURL:defaultImg,CategoryID:minuman.ID},
{Name:"Fanta",Price:6000,Stock:100,ImageURL:defaultImg,CategoryID:minuman.ID},
{Name:"Pocari Sweat",Price:7000,Stock:100,ImageURL:defaultImg,CategoryID:minuman.ID},
{Name:"Ultra Milk",Price:7000,Stock:100,ImageURL:defaultImg,CategoryID:minuman.ID},
{Name:"Good Day Coffee",Price:4000,Stock:100,ImageURL:defaultImg,CategoryID:minuman.ID},
{Name:"Yakult",Price:9000,Stock:100,ImageURL:defaultImg,CategoryID:minuman.ID},
{Name:"Floridina",Price:5000,Stock:100,ImageURL:defaultImg,CategoryID:minuman.ID},

// ================= PERTANIAN =================

{Name:"Bibit Padi",Price:50000,Stock:50,ImageURL:defaultImg,CategoryID:pertanian.ID},
{Name:"Bibit Jagung",Price:40000,Stock:50,ImageURL:defaultImg,CategoryID:pertanian.ID},
{Name:"Pupuk Urea",Price:120000,Stock:50,ImageURL:defaultImg,CategoryID:pertanian.ID},
{Name:"Pupuk Kompos",Price:60000,Stock:50,ImageURL:defaultImg,CategoryID:pertanian.ID},
{Name:"Cangkul",Price:80000,Stock:50,ImageURL:defaultImg,CategoryID:pertanian.ID},
{Name:"Sekop",Price:70000,Stock:50,ImageURL:defaultImg,CategoryID:pertanian.ID},
{Name:"Sprayer",Price:150000,Stock:50,ImageURL:defaultImg,CategoryID:pertanian.ID},
{Name:"Polybag",Price:20000,Stock:50,ImageURL:defaultImg,CategoryID:pertanian.ID},
{Name:"Bibit Cabai",Price:25000,Stock:50,ImageURL:defaultImg,CategoryID:pertanian.ID},
{Name:"Bibit Tomat",Price:25000,Stock:50,ImageURL:defaultImg,CategoryID:pertanian.ID},

// ================= MAINAN =================

{Name:"Lego Classic",Price:350000,Stock:30,ImageURL:defaultImg,CategoryID:mainan.ID},
{Name:"Hot Wheels",Price:50000,Stock:30,ImageURL:defaultImg,CategoryID:mainan.ID},
{Name:"Rubik 3x3",Price:40000,Stock:30,ImageURL:defaultImg,CategoryID:mainan.ID},
{Name:"Boneka Teddy",Price:80000,Stock:30,ImageURL:defaultImg,CategoryID:mainan.ID},
{Name:"RC Car",Price:250000,Stock:30,ImageURL:defaultImg,CategoryID:mainan.ID},
{Name:"Puzzle",Price:60000,Stock:30,ImageURL:defaultImg,CategoryID:mainan.ID},
{Name:"UNO Card",Price:30000,Stock:30,ImageURL:defaultImg,CategoryID:mainan.ID},
{Name:"Monopoly",Price:200000,Stock:30,ImageURL:defaultImg,CategoryID:mainan.ID},
{Name:"Action Figure Naruto",Price:150000,Stock:30,ImageURL:defaultImg,CategoryID:mainan.ID},
{Name:"Drone Mini",Price:500000,Stock:30,ImageURL:defaultImg,CategoryID:mainan.ID},

// ================= LAINNYA =================

{Name:"Kursi Plastik",Price:50000,Stock:30,ImageURL:defaultImg,CategoryID:lainnya.ID},
{Name:"Meja Lipat",Price:150000,Stock:30,ImageURL:defaultImg,CategoryID:lainnya.ID},
{Name:"Lampu LED",Price:40000,Stock:30,ImageURL:defaultImg,CategoryID:lainnya.ID},
{Name:"Kipas Angin",Price:200000,Stock:30,ImageURL:defaultImg,CategoryID:lainnya.ID},
{Name:"Jam Dinding",Price:70000,Stock:30,ImageURL:defaultImg,CategoryID:lainnya.ID},
{Name:"Botol Minum",Price:60000,Stock:30,ImageURL:defaultImg,CategoryID:lainnya.ID},
{Name:"Payung",Price:50000,Stock:30,ImageURL:defaultImg,CategoryID:lainnya.ID},
{Name:"Tas Ransel",Price:200000,Stock:30,ImageURL:defaultImg,CategoryID:lainnya.ID},
{Name:"Karpet",Price:150000,Stock:30,ImageURL:defaultImg,CategoryID:lainnya.ID},
{Name:"Bantal",Price:80000,Stock:30,ImageURL:defaultImg,CategoryID:lainnya.ID},

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