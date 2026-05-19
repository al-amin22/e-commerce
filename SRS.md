# Software Requirements Specification (SRS)
# E-Commerce Platform — Go + React

**Versi:** 1.0.0  
**Tanggal:** 2026-05-19  
**Status:** Active  
**Tim:** Full-Stack (Go Backend + React Frontend)

---

## Daftar Isi

1. [Pendahuluan](#1-pendahuluan)
2. [Gambaran Sistem](#2-gambaran-sistem)
3. [Arsitektur Teknis](#3-arsitektur-teknis)
4. [Modul 1: Auth & Identity](#modul-1-auth--identity)
5. [Modul 2: Product Catalog](#modul-2-product-catalog)
6. [Modul 3: Shopping Cart](#modul-3-shopping-cart)
7. [Modul 4: Order Management](#modul-4-order-management)
8. [Modul 5: Payment](#modul-5-payment)
9. [Modul 6: Shipping & Delivery](#modul-6-shipping--delivery)
10. [Modul 7: Admin Panel](#modul-7-admin-panel)
11. [Modul 8: User Profile & Address](#modul-8-user-profile--address)
12. [Non-Functional Requirements](#12-non-functional-requirements)
13. [Database Schema](#13-database-schema)
14. [API Contract](#14-api-contract)

---

## 1. Pendahuluan

### 1.1 Tujuan
Dokumen ini mendefinisikan spesifikasi lengkap untuk platform e-commerce berbasis Go (Gin + GORM + PostgreSQL + Redis) di backend dan React 18 (TypeScript + Vite + Tailwind CSS) di frontend. Setiap modul didefinisikan secara terpisah agar dapat diimplementasikan dan diuji secara independen.

### 1.2 Ruang Lingkup
Platform ini memungkinkan:
- **Buyer**: Mendaftar, browse produk, mengelola keranjang belanja, melakukan checkout, melacak pesanan.
- **Admin**: Mengelola produk, kategori, pesanan, dan pengguna melalui panel admin.

### 1.3 Stack Teknologi
| Layer | Teknologi |
|---|---|
| Backend | Go 1.22, Gin v1.10, GORM v1.25, PostgreSQL 15 |
| Cache / Session | Redis 7 |
| Auth | JWT (access + refresh token) |
| Email | SMTP (OTP verification) |
| Shipping | RajaOngkir API |
| Frontend | React 18, TypeScript, Vite, Tailwind CSS 3 |
| State Management | React Context API + localStorage |
| HTTP Client | Axios 1.7 |
| Form Handling | React Hook Form 7 |

---

## 2. Gambaran Sistem

```
┌─────────────────────────────────────────────────────┐
│                    FRONTEND (React)                  │
│  LoginPage │ RegisterPage │ ProductPage │ CartPage  │
│  OrderPage │ DashboardPage │ AdminPanel │ ProfilePage│
└────────────────────────┬────────────────────────────┘
                         │ HTTP REST API
┌────────────────────────▼────────────────────────────┐
│                    BACKEND (Go/Gin)                  │
│  Auth Handler │ Product Handler │ Cart Handler      │
│  Order Handler │ Payment Handler │ Shipping Handler │
│  Admin Handler │ Profile Handler                    │
├─────────────────────────────────────────────────────┤
│   Service Layer (Business Logic)                    │
├─────────────────────────────────────────────────────┤
│   Repository Layer (GORM / PostgreSQL)              │
├──────────────────┬──────────────────────────────────┤
│   PostgreSQL     │   Redis (Sessions + Cache)        │
└──────────────────┴──────────────────────────────────┘
```

### 2.1 Role & Permission
| Role | Permission |
|---|---|
| `guest` | Lihat produk, search, filter |
| `buyer` | Semua guest + cart, order, profile, address |
| `admin` | Semua buyer + kelola produk, kategori, order status, user |

---

## 3. Arsitektur Teknis

### 3.1 Struktur Direktori Backend
```
backend/
├── cmd/server/main.go              # Entry point
├── internal/
│   ├── config/config.go            # Environment config
│   ├── db/db.go                    # DB connection & migration
│   ├── models/                     # GORM models (semua tabel)
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── category.go
│   │   ├── cart.go
│   │   ├── order.go
│   │   ├── payment.go
│   │   ├── shipping.go
│   │   └── address.go
│   ├── repository/                 # Data access layer
│   │   ├── user_repository.go
│   │   ├── product_repository.go
│   │   ├── cart_repository.go
│   │   ├── order_repository.go
│   │   └── address_repository.go
│   ├── services/                   # Business logic
│   │   ├── auth_service.go
│   │   ├── product_service.go
│   │   ├── cart_service.go
│   │   ├── order_service.go
│   │   ├── payment_service.go
│   │   ├── shipping_service.go
│   │   └── email_service.go
│   ├── handlers/                   # HTTP handlers
│   │   ├── handler.go
│   │   ├── auth_handler.go
│   │   ├── product_handler.go
│   │   ├── cart_handler.go
│   │   ├── order_handler.go
│   │   ├── payment_handler.go
│   │   ├── shipping_handler.go
│   │   ├── address_handler.go
│   │   └── admin_handler.go
│   ├── middleware/
│   │   ├── auth_middleware.go      # JWT verification
│   │   └── role_middleware.go      # Role-based access
│   ├── cache/redis.go
│   └── utils/
│       ├── jwt.go
│       ├── password.go
│       ├── otp.go
│       └── response.go             # Standar JSON response
```

### 3.2 Struktur Direktori Frontend
```
frontend/src/
├── api/
│   ├── axios.ts                    # Axios instance + interceptors
│   ├── auth.ts                     # Auth API calls
│   ├── products.ts                 # Product API calls
│   ├── cart.ts                     # Cart API calls
│   ├── orders.ts                   # Order API calls
│   └── address.ts                  # Address API calls
├── context/
│   ├── AuthContext.tsx             # Auth state
│   └── CartContext.tsx             # Cart state
├── components/
│   ├── layout/
│   │   ├── Navbar.tsx
│   │   ├── Footer.tsx
│   │   └── Sidebar.tsx
│   ├── product/
│   │   ├── ProductCard.tsx
│   │   └── ProductGrid.tsx
│   ├── cart/
│   │   └── CartItem.tsx
│   └── ui/
│       ├── Button.tsx
│       ├── Input.tsx
│       └── Badge.tsx
└── pages/
    ├── LoginPage.tsx
    ├── RegisterPage.tsx
    ├── ProductListPage.tsx
    ├── ProductDetailPage.tsx
    ├── CartPage.tsx
    ├── CheckoutPage.tsx
    ├── OrderListPage.tsx
    ├── OrderDetailPage.tsx
    ├── DashboardPage.tsx
    ├── ProfilePage.tsx
    └── admin/
        ├── AdminDashboard.tsx
        ├── AdminProducts.tsx
        ├── AdminOrders.tsx
        └── AdminUsers.tsx
```

### 3.3 Standard JSON Response
```json
{
  "success": true,
  "message": "...",
  "data": {...},
  "error": null
}
```

Error response:
```json
{
  "success": false,
  "message": "Validation error",
  "data": null,
  "error": "email is required"
}
```

---

## Modul 1: Auth & Identity

### Deskripsi
Mengelola registrasi pengguna, verifikasi email, login/logout, dan manajemen token JWT.

### Requirements

#### FR-AUTH-001: Registrasi Pengguna
- User mengisi name, email, password
- Password di-hash menggunakan bcrypt (cost 12)
- OTP 6-digit dikirim ke email (expire 10 menit)
- Akun belum aktif hingga email diverifikasi
- Validasi: email unik, password min 8 karakter, name min 3 karakter

#### FR-AUTH-002: Verifikasi Email
- User submit email + OTP code
- OTP valid 10 menit, hanya boleh digunakan 1x
- Setelah verifikasi, akun menjadi aktif (is_verified = true)

#### FR-AUTH-003: Login
- User submit email + password
- Cek email sudah diverifikasi
- Issue access token (15 menit) + refresh token (7 hari)
- Refresh token disimpan di Redis + HTTP-only cookie
- Access token dikembalikan di response body

#### FR-AUTH-004: Refresh Token
- Client kirim refresh token (cookie atau body)
- Validasi token di Redis (JTI-based)
- Token lama direvoke, token baru di-issue
- Implementasi token rotation (mencegah replay attack)

#### FR-AUTH-005: Logout
- Revoke refresh token dari Redis
- Clear auth cookies

#### FR-AUTH-006: Get Current User (Me)
- Endpoint protected yang mengembalikan data user saat ini

### API Endpoints (Modul 1)
```
POST /api/v1/auth/register
POST /api/v1/auth/verify-email
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
GET  /api/v1/auth/me          [Protected]
```

### Frontend Pages (Modul 1)
- `LoginPage.tsx` — Form login dengan validasi
- `RegisterPage.tsx` — Form registrasi dengan validasi
- `AuthContext.tsx` — Global state: user, token, login(), logout(), register()

### Bug yang Harus Diperbaiki
1. `RegisterPage.tsx` memanggil `registerUser()` yang tidak ada di `AuthContext`
2. Axios response interceptor kosong — tidak ada auto-refresh saat 401
3. Tidak ada route `/verify-email` di frontend

---

## Modul 2: Product Catalog

### Deskripsi
Mengelola produk dan kategori. Guest bisa browse & search. Admin bisa CRUD produk.

### Requirements

#### FR-PROD-001: Daftar Produk (Publik)
- Endpoint publik mengembalikan daftar produk dengan paginasi
- Filter: category_id, min_price, max_price, search (by name)
- Sort: price_asc, price_desc, newest, oldest
- Default: 20 item per halaman

#### FR-PROD-002: Detail Produk (Publik)
- Mengembalikan detail lengkap produk berdasarkan ID/slug
- Termasuk: images, category, stock, description

#### FR-PROD-003: Kategori (Publik)
- Daftar semua kategori produk
- Kategori bisa bersarang (parent-child)

#### FR-PROD-004: CRUD Produk (Admin)
- Create, Update, Delete produk
- Upload gambar produk (base64 atau URL)
- Set stok, harga, kategori, status (active/inactive)

#### FR-PROD-005: Manajemen Stok
- Stok berkurang otomatis saat order dibuat
- Stok kembali jika order dibatalkan
- Validasi: tidak bisa order jika stok = 0

### Database Models (Modul 2)
```go
// Category
type Category struct {
    ID       uuid.UUID  `gorm:"primarykey"`
    Name     string     `gorm:"not null;unique"`
    Slug     string     `gorm:"unique"`
    ParentID *uuid.UUID // nullable untuk root category
    Parent   *Category
    Children []Category `gorm:"foreignKey:ParentID"`
}

// Product
type Product struct {
    ID          uuid.UUID `gorm:"primarykey"`
    Name        string    `gorm:"not null"`
    Slug        string    `gorm:"unique"`
    Description string
    Price       float64   `gorm:"not null"`
    Stock       int       `gorm:"not null;default:0"`
    ImageURL    string
    CategoryID  uuid.UUID
    Category    Category
    Weight      float64   // gram, untuk kalkulasi ongkir
    IsActive    bool      `gorm:"default:true"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### API Endpoints (Modul 2)
```
GET    /api/v1/products                    [Public] — list + filter + paginate
GET    /api/v1/products/:id                [Public] — detail produk
GET    /api/v1/categories                  [Public] — semua kategori
POST   /api/v1/admin/products              [Admin] — create produk
PUT    /api/v1/admin/products/:id          [Admin] — update produk
DELETE /api/v1/admin/products/:id          [Admin] — hapus produk
POST   /api/v1/admin/categories            [Admin] — create kategori
```

### Frontend Pages (Modul 2)
- `ProductListPage.tsx` — Grid produk dengan filter & pagination
- `ProductDetailPage.tsx` — Detail produk + tombol Add to Cart
- `ProductCard.tsx` — Komponen card produk reusable

---

## Modul 3: Shopping Cart

### Deskripsi
Keranjang belanja per-user. Disimpan di database (persistent) bukan localStorage.

### Requirements

#### FR-CART-001: Lihat Keranjang
- Mengembalikan semua item dalam keranjang user
- Termasuk subtotal per item dan total keseluruhan

#### FR-CART-002: Tambah Item ke Keranjang
- Validasi produk exists dan stok cukup
- Jika produk sudah ada di cart, tambahkan quantity
- Batas maksimum: quantity tidak boleh melebihi stok

#### FR-CART-003: Update Quantity
- Update jumlah item di keranjang
- Quantity 0 = hapus item dari keranjang

#### FR-CART-004: Hapus Item dari Keranjang
- Hapus item spesifik dari keranjang

#### FR-CART-005: Kosongkan Keranjang
- Hapus semua item dari keranjang (dipanggil setelah checkout sukses)

### Database Models (Modul 3)
```go
type Cart struct {
    ID        uuid.UUID  `gorm:"primarykey"`
    UserID    uuid.UUID  `gorm:"not null;uniqueIndex"`
    User      User
    Items     []CartItem
    CreatedAt time.Time
    UpdatedAt time.Time
}

type CartItem struct {
    ID        uuid.UUID `gorm:"primarykey"`
    CartID    uuid.UUID `gorm:"not null"`
    Cart      Cart
    ProductID uuid.UUID `gorm:"not null"`
    Product   Product
    Quantity  int       `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### API Endpoints (Modul 3)
```
GET    /api/v1/cart                   [Buyer] — lihat keranjang
POST   /api/v1/cart/items             [Buyer] — tambah item
PUT    /api/v1/cart/items/:id         [Buyer] — update quantity
DELETE /api/v1/cart/items/:id         [Buyer] — hapus item
DELETE /api/v1/cart                   [Buyer] — kosongkan keranjang
```

### Frontend Pages (Modul 3)
- `CartPage.tsx` — Daftar item, quantity control, subtotal, total, tombol Checkout
- `CartContext.tsx` — Global cart state (count, items)
- `CartItem.tsx` — Komponen item di cart
- Navbar badge: jumlah item di cart

---

## Modul 4: Order Management

### Deskripsi
Proses checkout dari cart ke order, manajemen status order.

### Requirements

#### FR-ORDER-001: Buat Order (Checkout)
- User memilih alamat pengiriman
- User memilih metode pengiriman (dari kalkulasi ongkir)
- Sistem membuat order dari item di cart
- Stok produk dikurangi
- Cart dikosongkan
- Order status: `pending_payment`

#### FR-ORDER-002: Daftar Order (Buyer)
- User melihat daftar semua ordernya
- Filter by status

#### FR-ORDER-003: Detail Order
- Detail lengkap order termasuk items, shipping info, status history

#### FR-ORDER-004: Batalkan Order
- Buyer bisa membatalkan order yang masih `pending_payment`
- Stok dikembalikan saat order dibatalkan

#### FR-ORDER-005: Update Status Order (Admin)
- Admin bisa update status: confirmed → shipped → delivered
- Setiap perubahan status dicatat di order_status_history

### Status Flow Order
```
pending_payment → paid → confirmed → shipped → delivered
                ↓
              cancelled
```

### Database Models (Modul 4)
```go
type Order struct {
    ID              uuid.UUID   `gorm:"primarykey"`
    UserID          uuid.UUID   `gorm:"not null"`
    User            User
    AddressID       uuid.UUID   `gorm:"not null"`
    Address         Address
    Status          OrderStatus `gorm:"default:'pending_payment'"`
    TotalAmount     float64
    ShippingCost    float64
    ShippingCourier string
    ShippingService string
    Notes           string
    Items           []OrderItem
    Payment         *Payment
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type OrderItem struct {
    ID        uuid.UUID `gorm:"primarykey"`
    OrderID   uuid.UUID `gorm:"not null"`
    ProductID uuid.UUID `gorm:"not null"`
    Product   Product
    Quantity  int
    UnitPrice float64   // harga saat checkout (snapshot)
    Subtotal  float64
}

type OrderStatus string
const (
    OrderStatusPendingPayment OrderStatus = "pending_payment"
    OrderStatusPaid           OrderStatus = "paid"
    OrderStatusConfirmed      OrderStatus = "confirmed"
    OrderStatusShipped        OrderStatus = "shipped"
    OrderStatusDelivered      OrderStatus = "delivered"
    OrderStatusCancelled      OrderStatus = "cancelled"
)
```

### API Endpoints (Modul 4)
```
POST   /api/v1/orders                    [Buyer] — checkout & buat order
GET    /api/v1/orders                    [Buyer] — daftar order saya
GET    /api/v1/orders/:id                [Buyer] — detail order
PUT    /api/v1/orders/:id/cancel         [Buyer] — batalkan order
GET    /api/v1/admin/orders              [Admin] — semua order
PUT    /api/v1/admin/orders/:id/status   [Admin] — update status
```

### Frontend Pages (Modul 4)
- `CheckoutPage.tsx` — Pilih alamat, metode kirim, konfirmasi total
- `OrderListPage.tsx` — Daftar order dengan filter status
- `OrderDetailPage.tsx` — Detail order + tracking status

---

## Modul 5: Payment

### Deskripsi
Pengelolaan pembayaran. Mendukung manual transfer (upload bukti) dan simulasi payment gateway.

### Requirements

#### FR-PAY-001: Upload Bukti Pembayaran
- Buyer upload bukti transfer (base64 image atau URL)
- Order status berubah menjadi `paid` (menunggu konfirmasi admin)

#### FR-PAY-002: Konfirmasi Pembayaran (Admin)
- Admin mengkonfirmasi atau menolak bukti pembayaran
- Jika dikonfirmasi: status → `confirmed`
- Jika ditolak: status tetap `pending_payment`, notifikasi dikirim

#### FR-PAY-003: Riwayat Pembayaran
- User melihat status pembayaran per order

### Database Models (Modul 5)
```go
type Payment struct {
    ID            uuid.UUID     `gorm:"primarykey"`
    OrderID       uuid.UUID     `gorm:"not null;uniqueIndex"`
    Order         Order
    Amount        float64
    Method        PaymentMethod `gorm:"default:'manual_transfer'"`
    Status        PaymentStatus `gorm:"default:'pending'"`
    ProofImageURL string        // bukti pembayaran
    ConfirmedAt   *time.Time
    ConfirmedBy   *uuid.UUID    // admin user ID
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type PaymentMethod string
const (
    PaymentMethodManualTransfer PaymentMethod = "manual_transfer"
)

type PaymentStatus string
const (
    PaymentStatusPending   PaymentStatus = "pending"
    PaymentStatusPaid      PaymentStatus = "paid"
    PaymentStatusRejected  PaymentStatus = "rejected"
)
```

### API Endpoints (Modul 5)
```
POST   /api/v1/orders/:id/payment           [Buyer] — upload bukti bayar
GET    /api/v1/orders/:id/payment           [Buyer] — status pembayaran
PUT    /api/v1/admin/payments/:id/confirm   [Admin] — konfirmasi pembayaran
PUT    /api/v1/admin/payments/:id/reject    [Admin] — tolak pembayaran
```

---

## Modul 6: Shipping & Delivery

### Deskripsi
Kalkulasi ongkos kirim menggunakan RajaOngkir API, dan pengelolaan data pengiriman.

### Requirements

#### FR-SHIP-001: Kalkulasi Ongkir
- Input: origin (kota toko), destination (kota user), weight, courier
- Output: daftar service + estimasi harga
- Fallback ke mock data jika API key tidak tersedia

#### FR-SHIP-002: Update Nomor Resi
- Admin input nomor resi saat order di-ship
- Buyer bisa melihat nomor resi di detail order

### API Endpoints (Modul 6)
```
GET  /api/v1/shipping/cost                    [Buyer] — kalkulasi ongkir
PUT  /api/v1/admin/orders/:id/tracking        [Admin] — input nomor resi
```

---

## Modul 7: Admin Panel

### Deskripsi
Dashboard admin untuk mengelola seluruh platform.

### Requirements

#### FR-ADMIN-001: Dashboard Stats
- Total revenue, total orders, total users, total products
- Orders by status (chart data)

#### FR-ADMIN-002: Manajemen Produk
- CRUD produk & kategori (sudah didefinisikan di Modul 2)

#### FR-ADMIN-003: Manajemen Order
- Lihat semua order, filter by status, update status
- Konfirmasi/tolak pembayaran

#### FR-ADMIN-004: Manajemen User
- Lihat daftar user, detail user
- Suspend/aktifkan akun user

### API Endpoints (Modul 7)
```
GET  /api/v1/admin/stats           [Admin] — dashboard statistics
GET  /api/v1/admin/users           [Admin] — daftar semua user
GET  /api/v1/admin/users/:id       [Admin] — detail user
PUT  /api/v1/admin/users/:id/suspend  [Admin] — suspend user
```

### Frontend Pages (Modul 7)
- `AdminDashboard.tsx` — Stats overview
- `AdminProducts.tsx` — Tabel produk + CRUD form
- `AdminOrders.tsx` — Tabel order + filter + update status
- `AdminUsers.tsx` — Tabel user + suspend action

---

## Modul 8: User Profile & Address

### Deskripsi
Pengelolaan profil pengguna dan alamat pengiriman.

### Requirements

#### FR-PROF-001: Lihat & Edit Profil
- User bisa update name dan password
- Update password: wajib submit password lama

#### FR-PROF-002: Manajemen Alamat
- User bisa CRUD alamat pengiriman
- Satu alamat bisa di-set sebagai default
- Maksimal 5 alamat per user

### API Endpoints (Modul 8)
```
GET    /api/v1/profile              [Buyer] — lihat profil
PUT    /api/v1/profile              [Buyer] — update profil
PUT    /api/v1/profile/password     [Buyer] — ganti password
GET    /api/v1/addresses            [Buyer] — daftar alamat
POST   /api/v1/addresses            [Buyer] — tambah alamat
PUT    /api/v1/addresses/:id        [Buyer] — update alamat
DELETE /api/v1/addresses/:id        [Buyer] — hapus alamat
PUT    /api/v1/addresses/:id/default [Buyer] — set default
```

### Frontend Pages (Modul 8)
- `ProfilePage.tsx` — Form edit profil + ganti password
- `AddressPage.tsx` — Daftar & form alamat

---

## 12. Non-Functional Requirements

| ID | Requirement |
|---|---|
| NFR-01 | Semua response API harus < 500ms untuk operasi normal |
| NFR-02 | Password di-hash dengan bcrypt cost 12 |
| NFR-03 | JWT access token expire 15 menit |
| NFR-04 | JWT refresh token expire 7 hari dengan token rotation |
| NFR-05 | Semua input user wajib divalidasi di backend |
| NFR-06 | Semua endpoint admin wajib cek role = "admin" |
| NFR-07 | CORS dikonfigurasi hanya untuk origin frontend |
| NFR-08 | Rate limiting: max 10 request/menit untuk auth endpoints |
| NFR-09 | Database menggunakan UUID v4 untuk semua primary key |
| NFR-10 | Semua waktu disimpan dalam UTC, ditampilkan dalam WIB |

---

## 13. Database Schema

### Tabel Utama
```sql
-- users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'buyer',
    is_verified BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- email_otps
CREATE TABLE email_otps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    code VARCHAR(10) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- categories
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    parent_id UUID REFERENCES categories(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- products
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    price DECIMAL(15,2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    image_url TEXT,
    category_id UUID REFERENCES categories(id),
    weight DECIMAL(10,2) DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- addresses
CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label VARCHAR(100),
    recipient VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address_line TEXT NOT NULL,
    city VARCHAR(255) NOT NULL,
    province VARCHAR(255) NOT NULL,
    postal_code VARCHAR(10) NOT NULL,
    courier_code VARCHAR(20),
    destination_id INT,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- carts
CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- cart_items
CREATE TABLE cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(cart_id, product_id)
);

-- orders
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    address_id UUID NOT NULL REFERENCES addresses(id),
    status VARCHAR(30) NOT NULL DEFAULT 'pending_payment',
    total_amount DECIMAL(15,2) NOT NULL,
    shipping_cost DECIMAL(15,2) NOT NULL DEFAULT 0,
    shipping_courier VARCHAR(50),
    shipping_service VARCHAR(100),
    tracking_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- order_items
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INT NOT NULL,
    unit_price DECIMAL(15,2) NOT NULL,
    subtotal DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- payments
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL UNIQUE REFERENCES orders(id),
    amount DECIMAL(15,2) NOT NULL,
    method VARCHAR(50) NOT NULL DEFAULT 'manual_transfer',
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    proof_image_url TEXT,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    confirmed_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

## 14. API Contract

### Standard Response Format
```json
{
  "success": true | false,
  "message": "string",
  "data": object | array | null,
  "error": "string | null",
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

### Pagination Query Params
```
?page=1&per_page=20&sort=created_at&order=desc
```

### Auth Headers
```
Authorization: Bearer <access_token>
```

---

*Dokumen ini adalah panduan utama implementasi. Setiap modul diimplementasikan satu per satu sesuai urutan di atas.*
