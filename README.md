# E-Commerce Stage 2: Authentication & Security (JWT)

Tujuan tahap ini: sistem login aman dengan hash password, JWT access/refresh token, middleware proteksi route, dan silent refresh di frontend.

## Stack

- Backend: Go, Gin, GORM, PostgreSQL, Redis, Viper, Bcrypt, JWT
- Frontend: React + Vite + TypeScript + Tailwind + React Hook Form + Axios Interceptor

## Struktur Backend (Clean Architecture)

- `backend/cmd/api/main.go` dan `backend/cmd/server/main.go` → thin entrypoint yang memanggil shared bootstrap
- `backend/cmd/server/main.go` → entrypoint resmi
- `backend/cmd/api/main.go` → entrypoint alternatif untuk kompatibilitas
- `backend/internal/app` → bootstrap aplikasi: load config, DB, Redis, service, route
- `backend/internal/models` → entity/database model utama (`User`, `Address`, `EmailOTP`, `Product`)
- `backend/internal/dto` → request/response contract API
- `backend/internal/handlers` → layer HTTP (request/response) untuk fitur utama
- `backend/internal/services` → business logic terpisah dari handler
- `backend/internal/response` → helper response JSON yang konsisten
- `backend/internal/db` → koneksi PostgreSQL + AutoMigrate
- `backend/internal/cache` → koneksi Redis
- `backend/internal/middleware` → proteksi JWT untuk route private
- `backend/internal/config` → baca `.env` dan buat DSN
- `backend/pkg/security` → helper JWT (generate/parse)

Catatan standar perusahaan:
- yang aktif dipakai hanya jalur di atas

Detail struktur final ada di [backend/ARCHITECTURE.md](backend/ARCHITECTURE.md).
Versi yang lebih enterprise-ready ada di [backend/ENTERPRISE_STRUCTURE.md](backend/ENTERPRISE_STRUCTURE.md).

## Auto-Migrate

Saat server start, sistem otomatis membuat/memperbarui tabel:
- `users`
- `email_otps`
- `addresses`
- `categories`
- `products`

Model ada di:
- [backend/internal/models/product.go](backend/internal/models/product.go)
- [backend/internal/models/models.go](backend/internal/models/models.go)

## Konfigurasi Env (Backend)

Copy file [backend/.env.example](backend/.env.example) menjadi `.env`, lalu isi:

- `APP_PORT=8080`
- `DB_HOST=localhost`
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=postgres`
- `DB_NAME=ecommerce_auth`
- `DB_SSLMODE=disable`

## Menjalankan Backend

Jalankan dari folder `backend`:

- `go mod tidy`
- `go run cmd/server/main.go`

Catatan: `cmd/api/main.go` masih tersedia, tetapi yang resmi dipakai adalah `cmd/server/main.go`.

Jika sukses, server aktif di `http://localhost:8080`.

## Endpoint Aktif Saat Ini

- `GET /health`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/verify-email`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me` (protected)
- `POST /api/v1/bootstrap-admin`
- `GET /api/v1/addresses` (protected)
- `POST /api/v1/addresses` (protected)
- `GET /api/v1/shipping/cost` (protected)
- `GET /api/v1/admin/users` (protected, admin only)

## Frontend Setup (React + Vite + TS + Tailwind)

File penting:
- [frontend/src/api/axios.ts](frontend/src/api/axios.ts) → Axios instance (`baseURL` backend)
- [frontend/src/main.tsx](frontend/src/main.tsx) → entrypoint React
- [frontend/src/App.tsx](frontend/src/App.tsx) → router login/register/dashboard
- [frontend/src/context/AuthContext.tsx](frontend/src/context/AuthContext.tsx) → state auth
- [frontend/src/pages/LoginPage.tsx](frontend/src/pages/LoginPage.tsx) + [frontend/src/pages/RegisterPage.tsx](frontend/src/pages/RegisterPage.tsx) → React Hook Form
- [frontend/tailwind.config.ts](frontend/tailwind.config.ts) + [frontend/postcss.config.js](frontend/postcss.config.js)

Env frontend:
- copy [frontend/.env.example](frontend/.env.example) ke `.env`
- isi `VITE_API_BASE_URL=http://localhost:8080/api/v1`

Run frontend dari folder `frontend`:
- `npm install`
- `npm run dev`

## Cara berpikirnya (untuk basic Laravel)

Konsepnya mirip Laravel, beda bahasa dan style:

- Laravel `Controller` ≈ Go `handler`
- Laravel `Service` ≈ Go `service`
- Laravel `Model/Eloquent` ≈ Go `models + db/repository (GORM)`
- Laravel `Auth Guard/Middleware` ≈ Go `middleware + JWT`
- Laravel `.env` ≈ Go `.env` (dibaca Viper)

Perbedaan utama:
- Laravel/PHP: dinamis
- Go: statically typed + compile, jadi error banyak ketahuan lebih awal.

Alur request juga sama:

$$
Client \to Handler \to Service \to Repository \to Database
$$

Alur auth yang baru:

$$
Login \to AccessToken(15m) + RefreshToken(7d) \to API Protected \to 401 \to Silent Refresh \to Retry Request
$$

## Yang Anda pelajari di Tahap 2

- `Hashing` dengan Bcrypt: password asli tidak pernah disimpan.
- `JWT`: access token sebagai tiket akses singkat.
- `Refresh token` di Redis: bisa revoke/rotasi token lama.
- `Middleware`: API private aman, hanya user login yang bisa akses.
- `Axios Interceptor`: token expired ditangani otomatis tanpa ganggu user.

Itu sebabnya kode lebih rapi, tidak campur query dan business logic di satu file.

# e-commerce
