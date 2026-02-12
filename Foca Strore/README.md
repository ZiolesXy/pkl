# E-commerce Backend API

Aplikasi backend e-commerce sederhana menggunakan Golang, Gin, GORM, PostgreSQL, dan JWT.

## Fitur Utama

- **Authentication**: Register, Login, Access Token (15 menit), Refresh Token (7 hari)
- **Role-based Authorization**: Admin dan User
- **Product Management**: CRUD produk (Admin only)
- **Cart System**: Add to cart, View cart, Remove item
- **Checkout System**: Checkout dengan status (pending, success, failed)
- **Seeder**: Endpoint untuk seeding data awal

---

## Instalasi

### Prasyarat

- Go 1.19+
- PostgreSQL 12+
- Git

### Langkah Instalasi

1. Clone repository ini:

```bash
git clone <repository-url>
cd ecommerce-app
```

2. Install dependencies:

```bash
go mod download
```

3. Buat file `.env` dari contoh:

```bash
cp .env.example .env
```

4. Edit file `.env` sesuai konfigurasi database Anda:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ecommerce_db
JWT_SECRET=your_super_secret_jwt_key_here_min_32_chars
```

---

## Setup Database

1. Buat database PostgreSQL:

```sql
CREATE DATABASE ecommerce_db;
```

2. Pastikan user PostgreSQL memiliki akses ke database tersebut.

---

## Migrasi Database

Aplikasi menggunakan auto-migration GORM.  
Tabel akan dibuat otomatis saat aplikasi pertama kali dijalankan.

---

## Menjalankan Aplikasi

1. Jalankan aplikasi:

```bash
go run main.go
```

2. Aplikasi akan berjalan di:

```
http://localhost:8080
```

---

## Seeder

Untuk mengisi data awal (roles, admin, users, products), akses endpoint:

```bash
GET http://localhost:8080/seed
```

Data yang di-seed:

- Roles: Admin, User
- 1 Admin (email: admin@ecommerce.com, password: admin123)
- 3 Users (john@example.com, jane@example.com, bob@example.com, password: password123)
- 5 Products

---

## API Endpoints

### Authentication

- `POST /register` - Register user baru
- `POST /login` - Login user
- `POST /refresh` - Refresh access token

### Products (Public)

- `GET /products` - Lihat semua produk

### Cart (User Only)

- `GET /api/cart` - Lihat keranjang
- `POST /api/cart/items` - Tambah item ke keranjang
- `DELETE /api/cart/items/:id` - Hapus item dari keranjang

### Checkout (User Only)

- `POST /api/checkout` - Checkout keranjang

### Products (Admin Only)

- `POST /api/admin/products` - Buat produk baru
- `PUT /api/admin/products/:id` - Update produk
- `DELETE /api/admin/products/:id` - Hapus produk
- `GET /api/admin/products` - Lihat semua produk

### Checkout (Admin Only)

- `PUT /api/admin/checkout/:id/status` - Update status checkout

---

## Format Response

### Success (single)

```json
{
  "status": "success",
  "message": "Success message",
  "data": {}
}
```

### Success (list)

```json
{
  "status": "success",
  "message": "Success message",
  "data": {
    "entries": []
  }
}
```

### Error

```json
{
  "status": "error",
  "message": "Error message"
}
```

---

## Testing dengan cURL

### Register

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Access Protected Endpoint

```bash
curl -X GET http://localhost:8080/api/cart \
  -H "Authorization: Bearer <access_token>"
```

### Admin Endpoint

```bash
curl -X POST http://localhost:8080/api/admin/products \
  -H "Authorization: Bearer <admin_access_token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"New Product","price":100000,"stock":50}'
```

---

## Keamanan

- Password di-hash menggunakan bcrypt
- JWT untuk autentikasi
- Role-based authorization
- Validasi input pada semua endpoint
- Error handling yang komprehensif

---

## Struktur Folder

```
ecommerce-app/
├── main.go
├── database/
│   └── database.go
├── models/
│   ├── role.go
│   ├── user.go
│   ├── product.go
│   ├── cart.go
│   ├── cart_item.go
│   └── checkout.go
├── request/
│   ├── auth.go
│   ├── product.go
│   ├── cart.go
│   └── checkout.go
├── response/
│   └── response.go
├── handlers/
│   ├── auth.go
│   ├── product.go
│   ├── cart.go
│   ├── checkout.go
│   └── seed.go
├── middleware/
│   ├── jwt.go
│   └── admin.go
├── helper/
│   └── jwt.go
├── seeders/
│   └── seed.go
└── README.md
```

---
