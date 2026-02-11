package seeder

import (
	"fmt"
	// "net/http"

	"foca-store/database"
	"foca-store/helper"
	"foca-store/models"
	"foca-store/response"

	"github.com/gin-gonic/gin"
)

func RunSeed(c *gin.Context) {

	/* =========================
	   1. SEED ROLES (IDEMPOTENT)
	========================= */

	roleNames := []string{
		"ADMIN",
		"USER",
	}

	var roles []models.Role
	for _, name := range roleNames {
		role := models.Role{Name: name}
		if err := database.DB.
			Where("name = ?", name).
			FirstOrCreate(&role).Error; err != nil {
			response.Error(c, 400, err.Error())
			return
		}
		roles = append(roles, role)
	}

	/* =========================
	   2. SEED PRODUCTS
	========================= */

	productSeeds := []models.Product{
		{Name: "MacBook Pro M2", Price: 25000000, Stock: 5},
		{Name: "Logitech MX Master 3", Price: 1500000, Stock: 20},
		{Name: "Keychron K2 V2", Price: 1200000, Stock: 15},
		{Name: "Dell UltraSharp 27", Price: 7000000, Stock: 10},
		{Name: "Sony WH-1000XM5", Price: 5000000, Stock: 8},
		{Name: "WD Black SN850 SSD", Price: 2000000, Stock: 25},
	}

	var products []models.Product
	for _, p := range productSeeds {
		product := p
		if err := database.DB.
			Where("name = ?", p.Name).
			FirstOrCreate(&product).Error; err != nil {
			response.Error(c, 400, err.Error())
			return
		}
		products = append(products, product)
	}

	/* =========================
	   3. SEED USERS
	========================= */

	type UserSeed struct {
		Name  string
		Email string
		Role  string
	}

	userSeeds := []UserSeed{
		{"Pasha", "admin@test.com", "ADMIN"},
		{"Acheron", "user1@test.com", "USER"},
		{"Siti", "user2@test.com", "USER"},
		{"Budi", "user3@test.com", "USER"},
	}

	var users []models.User

	for _, u := range userSeeds {

		var role models.Role
		if err := database.DB.
			Where("name = ?", u.Role).
			First(&role).Error; err != nil {
			response.Error(c, 400, "role not found")
			return
		}

		hashed, _ := helper.HashPassword("123456")

		user := models.User{
			Name:     u.Name,
			Email:    u.Email,
			Password: hashed,
			RoleID:   role.ID,
		}

		if err := database.DB.
			Where("email = ?", u.Email).
			FirstOrCreate(&user).Error; err != nil {
			response.Error(c, 400, err.Error())
			return
		}

		users = append(users, user)
	}

	/* =========================
	   4. SEED CART DUMMY
	========================= */

	for i := range users {

		var cart models.Cart
		database.DB.FirstOrCreate(&cart, models.Cart{
			UserID: users[i].ID,
		})

		// setiap user punya 2 product random
		for j := 0; j < 2; j++ {
			idx := (i + j) % len(products)

			item := models.CartItem{
				CartID:   cart.ID,
				ProductID: products[idx].ID,
				Quantity: 1,
			}

			database.DB.
				Where("cart_id = ? AND product_id = ?", cart.ID, products[idx].ID).
				FirstOrCreate(&item)
		}
	}

	/* =========================
	   5. RESPONSE
	========================= */

	message := fmt.Sprintf(
		"Seeder success: %d Role, %d User, %d Product",
		len(roles),
		len(users),
		len(products),
	)

	response.Success(c, message, nil)
}
