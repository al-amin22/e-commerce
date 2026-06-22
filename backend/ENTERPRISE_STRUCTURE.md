# Enterprise Backend Structure

Struktur ini adalah versi yang lebih mendekati praktik di perusahaan besar.

## Target Struktur

```text
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── bootstrap/
│   │   └── container.go
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   ├── user.go
│   │   ├── address.go
│   │   ├── product.go
│   │   └── otp.go
│   ├── dto/
│   │   ├── auth.go
│   │   ├── address.go
│   │   └── product.go
│   ├── response/
│   │   └── response.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── address_repository.go
│   │   └── product_repository.go
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── address_service.go
│   │   └── shipping_service.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── address_handler.go
│   │   ├── admin_handler.go
│   │   └── health_handler.go
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   └── role_middleware.go
│   ├── platform/
│   │   ├── database/
│   │   ├── cache/
│   │   └── mailer/
│   └── util/
│       ├── jwt.go
│       ├── password.go
│       └── otp.go
├── pkg/
│   └── security/
│       └── jwt.go
└── docs/
    └── architecture.md
```

## Prinsip Perusahaan

### 1. Entry point tipis
`cmd/server/main.go` hanya memanggil bootstrap aplikasi.

### 2. Dependency inversion
Handler tidak langsung mengurus detail database bila tidak perlu.
Service menerima interface repository.
Repository menjadi satu-satunya layer yang tahu detail query.

### 3. DTO terpisah dari entity
Request/response API sebaiknya tidak sama dengan model database.
Ini mengurangi risiko field bocor dan membuat kontrak API lebih stabil.

### 4. Bootstrap container
Di perusahaan besar, wiring dependency biasanya ada di satu tempat:
- config
- database
- redis
- logger
- mailer
- handler
- router

### 5. Konvensi penamaan
Pilih satu:
- `handler`, bukan campuran `handler` dan `handlers`
- `service`, bukan campuran `service` dan `services`
- `repository`, bukan campuran `repository` dan akses DB acak di handler

## Mapping dari repo ini

- `internal/app` → bootstrap aplikasi
- `internal/config` → central config
- `internal/db` → database platform layer
- `internal/models` → domain/entity
- `internal/handlers` → HTTP layer
- `internal/services` → business logic
- `internal/dto` → request/response contract yang jelas
- `internal/response` → helper JSON response yang konsisten
- `internal/cache` → infra/cache layer

## Saran langkah berikutnya

Kalau ingin benar-benar enterprise-ready, langkah berikutnya adalah:
1. pindahkan logic DB dari handler ke service
2. pecah request/response ke DTO
3. tambahkan repository layer yang konsisten
4. tambah logger dan structured error response
5. satukan penamaan folder menjadi satu gaya saja
