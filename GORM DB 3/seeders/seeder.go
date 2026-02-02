package seeders

import (
	"fmt"
	"main/database"
	"main/models"
	"main/respons"
	"net/http"

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
        "Async", "grace", "Furina", "Hutao", "Acheron",
        "Zephyro", "Jing yuan", "Feixiao", "Cyrene", "Vylan",
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
        
        // Mengambil 5 barang untuk setiap user secara bergantian
        for j := 0; j < 5; j++ {
            idx := (i + j) % len(barangs)
            barangUntukUser = append(barangUntukUser, barangs[idx])
        }

        // Simpan relasi Many-to-Many ke tabel pivot
        if err := database.DB.Model(&users[i]).Association("Barangs").Append(barangUntukUser); err != nil {
            c.JSON(400, gin.H{"error": "Gagal buat ownership: " + err.Error()})
            return
        }
    }

    totalUser := len(namaUser)
    totalBarang := len(namaBarang)
    totalRole := len(roles)

    messege := fmt.Sprintf("Seed data (%d User, %d Barang, %d Role berhasil)", totalUser, totalBarang, totalRole)

    c.JSON(
        http.StatusOK,
        respons.NewJsonResponse("Succes", messege),
    )
}