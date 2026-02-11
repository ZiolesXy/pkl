package seeder

import (
	"foca-store/database"
	"foca-store/models"
)

func SeedProducts() {
	database.DB.Create(&[]models.Product{
		{Name: "Laptop", Price: 15000000, Stock: 5},
		{Name: "Mouse", Price: 150000, Stock: 50},
	})
}
