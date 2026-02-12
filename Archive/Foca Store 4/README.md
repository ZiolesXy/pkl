# Foca Store API

## Tech

- Gin
- GORM
- PostgreSQL
- JWT (access 15m, refresh 7d)
- bcrypt

## Setup

1. Copy `.env.example` to `.env` and edit values.
2. Create database `DB_NAME` in PostgreSQL.
3. Install deps:

```bash
go mod tidy
```

4. Run:

```bash
go run .
```

## Endpoints

Base URL: `/api`

- `GET /health`
- `GET /seed`

### Auth

- `POST /register`
- `POST /login`
- `POST /refresh`
- `GET /me` (auth)

### Products

- `GET /products` (auth)
- `GET /products/:id` (auth)
- `POST /admin/products` (admin)
- `PUT /admin/products/:id` (admin)
- `DELETE /admin/products/:id` (admin)

### Cart

- `GET /cart` (auth)
- `POST /cart/items` (auth)
- `PUT /cart/items/:id` (auth)
- `DELETE /cart/items/:id` (auth)

### Checkout

- `POST /checkout` (auth)
- `GET /checkouts` (auth)
- `GET /admin/checkouts` (admin)
- `PUT /admin/checkouts/:id/status` (admin)

## JSON Response Format

```json
{
  "status": "success",
  "message": "string",
  "data": {}
}
```
