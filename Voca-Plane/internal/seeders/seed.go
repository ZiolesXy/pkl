package seeders

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"voca-plane/internal/domain/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) {
	log.Println(">>> STARTING FULL SEEDING...")
	SeedUsers(db)
	SeedAirlines(db)
	SeedAirports(db)
	// SeedFlights(db)
	SeedPromos(db)
	log.Println(">>> SEEDING COMPLETED")
}

func DropAll(db *gorm.DB) {
    log.Println(">>> DROPPING ALL TABLES...")

    // Daftar model yang ingin di-drop, urutkan dari yang memiliki Foreign Key terbanyak
    // untuk menghindari error constraint (atau gunakan DisableForeignKeyChecks)
    err := db.Migrator().DropTable(
        &models.FlightSeat{},
        &models.FlightClass{},
        &models.Flight{},
        &models.Airline{},
        &models.Airport{},
        &models.PromoCode{},
        &models.User{},
    )

    if err != nil {
        log.Fatalf(">>> FAILED TO DROP TABLES: %v", err)
    }
    
    log.Println(">>> ALL TABLES DROPPED SUCCESSFULLY")
}

func ResetDatabase(db *gorm.DB) {
    // 1. Drop semua tabel
    DropAll(db)

    // 2. Jalankan AutoMigrate lagi (sesuaikan dengan setup migrasi Anda)
    log.Println(">>> RE-MIGRATING TABLES...")
    db.AutoMigrate(
        &models.User{},
        &models.Airline{},
        &models.Airport{},
        &models.Flight{},
        &models.FlightClass{},
        &models.FlightSeat{},
        &models.PromoCode{},
    )

    // 3. Isi ulang data
    SeedAll(db)
}

func SeedUsers(db *gorm.DB) {
	log.Println(">>> Seeding Users...")
	
	// Super Admin
	superAdminPwd, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	superAdmin := models.User{
		Name:     "Super Admin",
		Email:    "superadmin@flightbooking.com",
		Password: string(superAdminPwd),
		Role:     models.RoleSuperAdmin,
	}
	db.FirstOrCreate(&superAdmin, models.User{Email: superAdmin.Email})

	// Admin
	adminPwd, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Name:     "Admin User",
		Email:    "admin@flightbooking.com",
		Password: string(adminPwd),
		Role:     models.RoleAdmin,
	}
	db.FirstOrCreate(&admin, models.User{Email: admin.Email})

	// Regular User
	userPwd, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	user := models.User{
		Name:     "John Doe",
		Email:    "user@flightbooking.com",
		Password: string(userPwd),
		Role:     models.RoleUser,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
}

func SeedAirlines(db *gorm.DB) {
	log.Println(">>> Seeding Airlines from JSON...")
	data, err := os.ReadFile("internal/seeders/airlines.json")
	if err != nil {
		log.Printf("Failed to read airlines.json: %v", err)
		return
	}

	var airlines []models.Airline
	if err := json.Unmarshal(data, &airlines); err != nil {
		log.Printf("Failed to unmarshal airlines.json: %v", err)
		return
	}

	for _, a := range airlines {
		db.FirstOrCreate(&a, models.Airline{Code: a.Code})
	}
}

func SeedAirports(db *gorm.DB) {
	log.Println(">>> Seeding Airports from JSON...")
	data, err := os.ReadFile("internal/seeders/airports.json")
	if err != nil {
		log.Printf("Failed to read airports.json: %v", err)
		return
	}

	var airportMap map[string]struct {
		IATA string `json:"iata"`
		Name string `json:"name"`
		City string `json:"city"`
	}
	if err := json.Unmarshal(data, &airportMap); err != nil {
		log.Printf("Failed to unmarshal airports.json: %v", err)
		return
	}

	// We only seed a subset if it's too many, or only those with IATA code
	count := 0
	for _, v := range airportMap {
		if v.IATA == "" {
			continue
		}
		airport := models.Airport{
			Code: v.IATA,
			Name: v.Name,
			City: v.City,
		}
		db.FirstOrCreate(&airport, models.Airport{Code: airport.Code})
		count++
		if count >= 100 { // Limit to 100 for now to avoid massive DB bloat if not needed
			break
		}
	}
}

func SeedFlights(db *gorm.DB) {
	log.Println(">>> Seeding 20 Sample Flights with new seat system...")
	var airlines []models.Airline
	var airports []models.Airport

	db.Find(&airlines)
	db.Find(&airports)

	if len(airlines) == 0 || len(airports) < 2 {
		log.Println("Insufficient airlines or airports for flight seeding")
		return
	}

	type classConfig struct {
		ClassType string
		Price     float64
		Percent   float64
	}

	for i := 1; i <= 20; i++ {
		airline := airlines[i%len(airlines)]
		origin := airports[i%len(airports)]
		dest := airports[(i+1)%len(airports)]

		// Vary class count: 1-7 = 1 class, 8-14 = 2 classes, 15-20 = 3 classes
		var classes []classConfig
		totalSeats := 120
		totalRows := 20
		totalColumns := 6

		switch {
		case i <= 7:
			// 1 class: Economy 100%
			classes = []classConfig{
				{"Economy", 1000000 + float64(i*50000), 1.0},
			}
		case i <= 14:
			// 2 classes: Business 30%, Economy 70%
			classes = []classConfig{
				{"Business", 3000000 + float64(i*100000), 0.30},
				{"Economy", 1000000 + float64(i*50000), 0.70},
			}
		default:
			// 3 classes: First 20%, Business 30%, Economy 50%
			classes = []classConfig{
				{"First", 5000000 + float64(i*150000), 0.20},
				{"Business", 3000000 + float64(i*100000), 0.30},
				{"Economy", 1000000 + float64(i*50000), 0.50},
			}
		}

		flight := models.Flight{
			AirlineID:     airline.ID,
			OriginID:      origin.ID,
			DestinationID: dest.ID,
			DepartureTime: time.Now().Add(time.Duration(i*6) * time.Hour),
			ArrivalTime:   time.Now().Add(time.Duration(i*6+2) * time.Hour),
			FlightNumber:  fmt.Sprintf("FL%03d", i),
			TotalSeats:    totalSeats,
			TotalRows:     totalRows,
			TotalColumns:  totalColumns,
		}
		db.FirstOrCreate(&flight, models.Flight{FlightNumber: flight.FlightNumber})

		// Generate classes and seats
		seatIndex := 0
		for _, c := range classes {
			seatCount := int(float64(totalSeats) * c.Percent)
			// Last class gets remaining seats to avoid rounding issues
			if c.ClassType == classes[len(classes)-1].ClassType {
				usedSeats := 0
				for _, prev := range classes[:len(classes)-1] {
					usedSeats += int(float64(totalSeats) * prev.Percent)
				}
				seatCount = totalSeats - usedSeats
			}

			fClass := models.FlightClass{
				FlightID:  flight.ID,
				ClassType: c.ClassType,
				Price:     c.Price,
			}
			db.FirstOrCreate(&fClass, models.FlightClass{FlightID: flight.ID, ClassType: fClass.ClassType})

			for j := 0; j < seatCount; j++ {
				row := seatIndex / totalColumns
				col := (seatIndex % totalColumns) + 1
				rowLetter := string(rune('A' + row))
				seatNumber := fmt.Sprintf("%s%d", rowLetter, col)

				seat := models.FlightSeat{
					FlightClassID: fClass.ID,
					SeatNumber:    seatNumber,
					IsAvailable:   true,
				}
				db.FirstOrCreate(&seat, models.FlightSeat{FlightClassID: fClass.ID, SeatNumber: seat.SeatNumber})
				seatIndex++
			}
		}
	}
}

func SeedPromos(db *gorm.DB) {
	log.Println(">>> Seeding Promos...")
	promos := []models.PromoCode{
		{Code: "HEMAT50", Discount: 50, IsActive: true},
		{Code: "NEWUSER", Discount: 20, IsActive: true},
	}
	for _, p := range promos {
		db.FirstOrCreate(&p, models.PromoCode{Code: p.Code})
	}
}

func InitSeeders(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		SeedAll(db)
	}
}
