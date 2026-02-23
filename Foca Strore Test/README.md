# Foca Store Backend API

Aplikasi backend e-commerce sederhana menggunakan Golang, Gin, GORM, PostgreSQL, dan JWT.

## Fitur Utama

- **Authentication**: Register, Login, Access Token (15 menit), Refresh Token (7 hari)
- **Role-based Authorization**: Admin dan User
- **Product Management**: CRUD produk (Admin only)
- **Category Management**: CRUD kategori produk (Admin only)
- **Address Management**: CRUD alamat pengiriman user
- **Cart System**: Add to cart, View cart, Remove item
- **Checkout System**: Checkout dengan status (pending, success, failed)
- **Coupon System**: Buat, klaim, dan gunakan kupon diskon
- **Image Upload**: Cloudinary integration untuk product images dan profile images
- **Seeder**: Endpoint untuk seeding data awal

---

## Instalasi

### Prasyarat

- Go 1.24+
- PostgreSQL 12+
- Git

### Langkah Instalasi

1. Clone repository ini:

```bash
git clone <repository-url>
cd "Foca Store"
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
DB_NAME=foca_store
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
CREATE DATABASE foca_store;
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

# Seed sample coupons
GET http://localhost:8080/seed/coupons

# Sync asset products
GET http://localhost:8080/seed/sync

# Seed all data (roles, admin, users)
GET http://localhost:8080/seed/all
```

Data yang di-seed:

- **Roles**: Admin, User
- **Admin**: 1 admin (email: admin@foca-store.com, password: admin123)
- **Users**: 3 sample users dengan password: password123
- **Products**: 5 sample products (dapat ditambah melalui endpoint /seed/products atau /seed/assets)
- **Coupons**: Sample coupons untuk testing

---

## API Endpoints

### Authentication (Public)

- `POST /register` - Register user baru
- `POST /login` - Login user
- `POST /refresh` - Refresh access token

### Categories (Public)

- `GET /category` - Lihat semua kategori
- `GET /category/:slug` - Lihat detail kategori berdasarkan slug

### Products (Public)

- `GET /products` - Lihat semua produk (public endpoint)
- `GET /product/:slug` - Lihat detail produk berdasarkan slug

### Coupons (Public)

- `GET /coupons` - Lihat semua kupon yang tersedia

### User Profile (Protected)

- `GET /api/profile` - Lihat profil user
- `PUT /api/profile` - Update profil user (support JSON & multipart/form-data untuk upload gambar)

### Address Management (Protected)

- `POST /api/addresses` - Tambah alamat baru
- `GET /api/addresses` - Lihat semua alamat user
- `GET /api/addresses/:uid` - Lihat detail alamat berdasarkan UID
- `PUT /api/addresses/:uid` - Update alamat
- `DELETE /api/addresses/:uid` - Hapus alamat

### Cart (Protected)

- `GET /api/cart` - Lihat keranjang
- `POST /api/cart/items` - Tambah item ke keranjang
- `DELETE /api/cart/items/:id` - Hapus item dari keranjang

### Checkout (Protected)

- `POST /api/checkout` - Checkout keranjang
- `GET /api/checkout/me` - Lihat history checkout user

### Coupon Management (Protected)

- `POST /api/coupons/claim` - Klaim kupon
- `GET /api/coupons/me` - Lihat kupon yang dimiliki user
- `DELETE /api/coupons/:id/remove` - Hapus kupon dari user

### Products Management (Admin Only)

- `POST /api/admin/products` - Buat produk baru
- `PUT /api/admin/products/:id` - Update produk
- `DELETE /api/admin/products/:id` - Hapus produk
- `DELETE /api/admin/products` - Hapus semua produk
- `DELETE /api/admin/products/assets` - Hapus semua gambar produk

### Categories Management (Admin Only)

- `POST /api/admin/category` - Buat kategori baru
- `PUT /api/admin/category/:id` - Update kategori
- `DELETE /api/admin/category/:id` - Hapus kategori

### Coupons Management (Admin Only)

- `POST /api/admin/coupons` - Buat kupon baru
- `PUT /api/admin/coupon/:id` - Update kupon
- `DELETE /api/admin/coupon/:id` - Hapus kupon

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

### Create Address

```bash
curl -X POST http://localhost:8080/api/addresses \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "label": "Rumah",
    "recipient_name": "John Doe",
    "phone": "08123456789",
    "address_line": "Jl. Contoh No. 123",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12345"
  }'
```

### Get My Addresses

```bash
curl -X GET http://localhost:8080/api/addresses \
  -H "Authorization: Bearer <access_token>"
```

### Update Profile with Image

```bash
curl -X PUT http://localhost:8080/api/profile \
  -H "Authorization: Bearer <access_token>" \
  -F "name=Updated Name" \
  -F "profile_image=@/path/to/image.jpg"
```

### Admin Endpoint

```bash
curl -X POST http://localhost:8080/api/admin/products \
  -H "Authorization: Bearer <admin_access_token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"New Product","price":100000,"stock":50,"categoryId":1}'
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
Foca Strore Test/
├── main.go
├── database/
│   └── database.go
├── models/
│   ├── role.go
│   ├── user.go
│   ├── category.go
│   ├── product.go
│   ├── cart.go
│   ├── cart_item.go
│   ├── checkout.go
│   ├── checkout_item.go
│   ├── address.go
│   ├── coupon.go
│   ├── user_coupon.go
│   └── refresh_token.go
├── request/
│   ├── auth.go
│   ├── category.go
│   ├── product.go
│   ├── cart.go
│   ├── checkout.go
│   ├── profile.go
│   └── address.go
├── response/
│   ├── response.go
│   ├── auth.go
│   ├── category.go
│   ├── product.go
│   ├── cart.go
│   ├── checkout.go
│   ├── profile.go
│   ├── address.go
│   ├── coupon.go
│   ├── role.go
│   └── user.go
├── handlers/
│   ├── auth.go
│   ├── product.go
│   ├── profile.go
│   ├── address.go
│   ├── cart.go
│   ├── checkout.go
│   ├── category.go
│   ├── coupon.go
│   ├── user_coupon.go
│   └── seed.go
├── middleware/
│   ├── jwt.go
│   └── admin.go
├── helper/
│   ├── jwt.go
│   ├── cloudinary.go
│   ├── slug.go
│   ├── password.go
│   └── context.go
├── seeders/
│   └── seed.go
├── trash/
│   └── trash.go
├── AssetPrivate/
└── README.md
```

---
