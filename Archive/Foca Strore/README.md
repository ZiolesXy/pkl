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
cd voca-store
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
DB_NAME=voca_store
JWT_SECRET=your_super_secret_jwt_key_here_min_32_chars
CLOUDINARY_CLOUD_NAME=your_cloudinary_cloud_name
CLOUDINARY_API_KEY=your_cloudinary_api_key
CLOUDINARY_API_SECRET=your_cloudinary_api_secret
PORT=8080
```

---

## Setup Database

1. Buat database PostgreSQL:

```sql
CREATE DATABASE voca_store;
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

Untuk mengisi data awal, terdapat beberapa endpoint seeding:

```bash
# Seed roles (Admin, User)
GET http://localhost:8080/seed/roles

# Seed admin user (email: admin@voca-store.com, password: admin123)
GET http://localhost:8080/seed/admin

# Seed sample users
GET http://localhost:8080/seed/users

# Seed sample products
GET http://localhost:8080/seed/products

# Seed products from assets folder
GET http://localhost:8080/seed/assets

# Seed all data (roles, admin, users)
GET http://localhost:8080/seed/all
```

Data yang di-seed:

- **Roles**: Admin, User
- **Admin**: 1 admin (email: admin@voca-store.com, password: admin123)
- **Users**: 3 sample users dengan password: password123
- **Products**: 5 sample products (dapat ditambah melalui endpoint /seed/products atau /seed/assets)

---

## API Endpoints

### Authentication

- `POST /register` - Register user baru
- `POST /login` - Login user
- `POST /refresh` - Refresh access token

### Products (Public)

- `GET /products` - Lihat semua produk (public endpoint)

### User Profile (Protected)

- `GET /api/profile` - Lihat profil user
- `PUT /api/profile` - Update profil user

### Products (Protected)

- `GET /api/products` - Lihat semua produk (dengan detail lengkap)
- `GET /api/product/:slug` - Lihat detail produk berdasarkan slug

### Cart (Protected)

- `GET /api/cart` - Lihat keranjang
- `POST /api/cart/items` - Tambah item ke keranjang
- `DELETE /api/cart/items/:id` - Hapus item dari keranjang

### Checkout (Protected)

- `POST /api/checkout` - Checkout keranjang

### Products Management (Admin Only)

- `POST /api/admin/products` - Buat produk baru
- `PUT /api/admin/products/:id` - Update produk
- `DELETE /api/admin/products/:id` - Hapus produk

### Checkout Management (Admin Only)

- `GET /api/admin/checkout` - Lihat semua checkout
- `PATCH /api/admin/checkout/:id/approve` - Approve checkout
- `PATCH /api/admin/checkout/:id/reject` - Reject checkout

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
voca-store/
├── main.go
├── database/
│   └── database.go
├── models/
│   ├── role.go
│   ├── user.go
│   ├── product.go
│   ├── cart.go
│   ├── cart_item.go
│   ├── checkout.go
│   └── checkout_item.go
├── request/
│   ├── auth.go
│   ├── product.go
│   ├── cart.go
│   ├── checkout.go
│   └── profile.go
├── response/
│   ├── response.go
│   ├── auth.go
│   ├── product.go
│   ├── cart.go
│   ├── checkout.go
│   └── profile.go
├── handlers/
│   ├── auth.go
│   ├── product.go
│   ├── profile.go
│   ├── cart.go
│   ├── checkout.go
│   └── seed.go
├── middleware/
│   ├── jwt.go
│   └── admin.go
├── helper/
│   ├── jwt.go
│   ├── cloudinary.go
│   └── slug.go
├── seeders/
│   └── seed.go
├── trash/
│   └── trash.go
├── AssetPrivate/
└── README.md
```

---
