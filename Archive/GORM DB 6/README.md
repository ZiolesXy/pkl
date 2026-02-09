# Belajar GORM v2 + PostgreSQL + Gin (CRUD Lengkap)

> Dokumentasi ini dirangkum dari PDF + lanjutan CRUD & API JSON. **Siap copy–paste** sebagai `README.md`.

---

## STEP 1 — Setup Project

```bash
mkdir gorm-gin-basic
cd gorm-gin-basic
go mod init gorm-gin-basic

go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/gin-gonic/gin
```

### Struktur Folder

```
gorm-gin-basic/
├── main.go
├── database/
│   └── database.go
├── models/
│   ├── role.go
│   ├── user.go
│   └── barang.go
├── handlers/
│   ├── role_handler.go
│   ├── user_handler.go
│   └── barang_handler.go
└── seeders/
    └── seeder.go
```

---

## STEP 2 — Database Connection

```go
package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=postgres dbname=gorm_basic port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal koneksi database")
	}
	DB = db
	fmt.Println("Database terkoneksi")
}
```

---

## STEP 3 — Models

### Role

```go
type Role struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"type:varchar(50);not null"`
	Users []User `gorm:"foreignKey:RoleID"`
}
```

### User

```go
type User struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"type:varchar(100);not null"`
	RoleID  uint
	Role    Role
	Barangs []Barang `gorm:"many2many:user_barangs"`
}
```

### Barang

```go
type Barang struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"type:varchar(100);not null"`
	Users []User `gorm:"many2many:user_barangs"`
}
```

---

## STEP 4 — Migration & Router

```go
func main() {
	database.ConnectDatabase()
	database.DB.AutoMigrate(&models.Role{}, &models.User{}, &models.Barang{})
	seeders.RunSeeder()

	r := gin.Default()

	r.POST("/roles", handlers.CreateRole)
	r.GET("/roles", handlers.GetRoles)
	r.PUT("/roles/:id", handlers.UpdateRole)
	r.DELETE("/roles/:id", handlers.DeleteRole)

	r.POST("/users", handlers.CreateUser)
	r.GET("/users", handlers.GetUsers)
	r.GET("/users/:id", handlers.GetUserByID)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)

	r.POST("/barangs", handlers.CreateBarang)
	r.GET("/barangs", handlers.GetBarangs)
	r.PUT("/barangs/:id", handlers.UpdateBarang)
	r.DELETE("/barangs/:id", handlers.DeleteBarang)

	r.POST("/users/:id/barangs", handlers.AssignBarang)
	r.DELETE("/users/:id/barangs/:barang_id", handlers.RemoveBarang)

	r.Run(":8080")
}
```

---

## STEP 5 — CRUD API JSON

### CREATE ROLE

```http
POST /roles
```

```json
{ "name": "Admin" }
```

### CREATE USER

```http
POST /users
```

```json
{ "name": "Budi", "role_id": 1 }
```

### CREATE BARANG

```http
POST /barangs
```

```json
{ "name": "Laptop" }
```

### ASSIGN BARANG KE USER

```http
POST /users/1/barangs
```

```json
{ "id": 1 }
```

### GET USERS

```http
GET /users
```

### GET USER DETAIL

```http
GET /users/1
```

### UPDATE USER

```http
PUT /users/1
```

```json
{ "name": "Budi Update", "role_id": 2 }
```

### DELETE USER

```http
DELETE /users/1
```

### REMOVE BARANG DARI USER

```http
DELETE /users/1/barangs/1
```

---

## FINAL CHECKLIST

* [x] CRUD Role
* [x] CRUD User
* [x] CRUD Barang
* [x] Many to Many Assign & Remove
* [x] PostgreSQL FK Aman
* [x] Siap Production Dasar

---

🔥 **END — GORM + GIN CRUD COMPLETE**
