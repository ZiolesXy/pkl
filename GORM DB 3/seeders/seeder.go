package seeders

import (
	"main/database"
	"main/models"
	"github.com/gin-gonic/gin"
)

func RunSeed(c *gin.Context) {
    // 1. Seed Roles (5 Role)
    roles := []models.Role{
        {Name: "Admin"},
        {Name: "Manager"},
        {Name: "Staff IT"},
        {Name: "Sales"},
        {Name: "Inventory"},
    }
    if err := database.DB.Create(&roles).Error; err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 2. Seed Barang (20 Barang Asli)
    namaBarang := []string{
        "MacBook Pro M2", "Logitech MX Master 3", "Keychron K2 V2", "Dell UltraSharp 27", "Epson L3210 Printer",
        "Sony WH-1000XM5", "iPad Air 5", "Samsung Galaxy S23", "SteelSeries Mousepad", "ThinkPad X1 Carbon",
        "Wacom Intuos Pro", "Cisco Router ISR", "Ubiquiti UniFi AP", "GoPro Hero 11", "BenQ SW271 Monitor",
        "WD Black SN850 SSD", "Seagate IronWolf 4TB", "Asus ROG Strix VGA", "Corsair RM850 PSU", "Razer BlackWidow",
    }
    
    var barangs []models.Barang
    for _, name := range namaBarang {
        barangs = append(barangs, models.Barang{Name: name})
    }
    if err := database.DB.Create(&barangs).Error; err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 3. Seed Users (10 User)
    namaUser := []string{
        "Pasha", "Amrl", "Siti", "Budi", "Dewi",
        "Eko", "Farhan", "Gita", "Hadi", "Indah",
    }
    
    var users []models.User
    for i, name := range namaUser {
        users = append(users, models.User{
            Name:   name,
            RoleID: roles[i%len(roles)].ID, // Distribusi ke 5 role
        })
    }
    if err := database.DB.Create(&users).Error; err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 4. Seed Ownership (1 User memiliki 10 Barang)
    for i := range users {
        var barangUntukUser []models.Barang
        
        // Mengambil 10 barang untuk setiap user secara bergantian
        for j := 0; j < 10; j++ {
            idx := (i + j) % len(barangs)
            barangUntukUser = append(barangUntukUser, barangs[idx])
        }

        // Simpan relasi Many-to-Many ke tabel pivot
        if err := database.DB.Model(&users[i]).Association("Barangs").Append(barangUntukUser); err != nil {
            c.JSON(400, gin.H{"error": "Gagal buat ownership: " + err.Error()})
            return
        }
    }

    // Respon Tunggal
    c.JSON(200, gin.H{"message": "Seed data (5 Role, 10 User, 20 Barang Asli) berhasil!"})
}

func ClearData(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	for _, user := range users {
		if err := database.DB.Model(&user).Association("Barangs").Clear(); err != nil {
			c.JSON(400, gin.H{"Error1": err.Error()})
			return
		}
		if err := database.DB.Unscoped().Delete(&user).Error; err != nil {
			c.JSON(400, gin.H{"Error2": err.Error()})
			return
		}
		// if err := database.DB.Unscoped().Delete(&models.Role{}).Error; err != nil {
		// 	c.JSON(400, gin.H{"Error3": err.Error()})
		// 	return
		// }
	}

	// if err := database.DB.Unscoped().Delete(&models.User{}).Error; err != nil {
	// 	c.JSON(400, gin.H{"Error2": err.Error()})
	// 	return
	// }

	// if err := database.DB.Unscoped().Delete(&models.Role{}).Error; err != nil {
	// 	c.JSON(400, gin.H{"Error3": err.Error()})
	// 	return
	// }

	// c.JSON(200, gin.H{
	// 	"messege": "All table data cleared",
	// })
}