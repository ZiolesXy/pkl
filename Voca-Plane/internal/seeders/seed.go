package seeders

import (
	"fmt"
	"log"
	"time"
	"voca-plane/internal/domain/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func InitSeeders(db *gorm.DB) {
	log.Println(">>> MEMERIKSA DATA SEEDER...")

	var count int64
	db.Model(&models.Airline{}).Count(&count)
	if count > 0 {
		log.Println(">>> SEEDER DIBATALKAN: Data sudah ada")
		return
	}

	// ======================
	// 1. Airlines
	// ======================
	airlines := []models.Airline{
		{Name: "Garuda Indonesia", Code: "GA"},
		{Name: "Lion Air", Code: "JT"},
		{Name: "AirAsia", Code: "QZ"},
		{Name: "Citilink", Code: "QG"},
		{Name: "Batik Air", Code: "ID"},
	}

	for i := range airlines {
		db.FirstOrCreate(&airlines[i], models.Airline{Code: airlines[i].Code})
	}

	// ======================
	// 2. Airports
	// ======================
	airports := []models.Airport{
		{Code: "CGK", Name: "Soekarno-Hatta", City: "Jakarta"},
		{Code: "DPS", Name: "Ngurah Rai", City: "Bali"},
		{Code: "SUB", Name: "Juanda", City: "Surabaya"},
		{Code: "KNO", Name: "Kualanamu", City: "Medan"},
		{Code: "UPG", Name: "Sultan Hasanuddin", City: "Makassar"},
	}

	for i := range airports {
		db.FirstOrCreate(&airports[i], models.Airport{Code: airports[i].Code})
	}

	// ======================
	// 3. Flight
	// ======================
	var flight models.Flight

	db.FirstOrCreate(&flight, models.Flight{
		FlightNumber: "GA400",
	})

	flight.AirlineID = airlines[0].ID
	flight.OriginID = airports[0].ID
	flight.DestinationID = airports[1].ID
	flight.DepartureTime = time.Now().Add(24 * time.Hour)
	flight.ArrivalTime = time.Now().Add(26 * time.Hour)

	db.Save(&flight)

	// ======================
	// 4. Flight Class
	// ======================
	var fClass models.FlightClass

	db.FirstOrCreate(&fClass, models.FlightClass{
		FlightID:  flight.ID,
		ClassType: "Economy",
	})

	fClass.Price = 1500000
	db.Save(&fClass)

	// ======================
	// 5. Seats
	// ======================
	for i := 0; i < 50; i++ {

		row := string(rune('A' + (i / 10)))
		num := i % 10

		seat := models.FlightSeat{
			FlightClassID: fClass.ID,
			SeatNumber:    fmt.Sprintf("%s%d", row, num),
			IsAvailable:   true,
		}

		db.FirstOrCreate(&seat, models.FlightSeat{
			FlightClassID: fClass.ID,
			SeatNumber:    seat.SeatNumber,
		})
	}

	// ======================
	// 6. Promo
	// ======================
	promos := []models.PromoCode{
		{Code: "HEMAT50", Discount: 50, IsActive: true},
		{Code: "NEWUSER", Discount: 20, IsActive: true},
		{Code: "LIBURAN", Discount: 10, IsActive: true},
	}

	for i := range promos {
		db.FirstOrCreate(&promos[i], models.PromoCode{
			Code: promos[i].Code,
		})
	}

	// ======================
	// 7. Super Admin
	// ======================
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	admin := models.User{
		Name:     "Super Admin",
		Email:    "admin@flightbooking.com",
		Password: string(hashedPassword),
		Role:     models.RoleAdmin,
	}

	db.FirstOrCreate(&admin, models.User{
		Email: admin.Email,
	})

	// ======================
	// 8. Regular User
	// ======================
	userPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)

	user := models.User{
		Name:     "John Doe",
		Email:    "user@flightbooking.com",
		Password: string(userPassword),
		Role:     models.RoleUser,
	}

	db.FirstOrCreate(&user, models.User{
		Email: user.Email,
	})

	log.Println(">>> SEEDER SELESAI")
}