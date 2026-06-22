# Backend Architecture

Struktur final yang dipakai sebagai standar:

```text
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── cache/
│   │   └── redis.go
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   └── db.go
│   ├── handlers/
│   │   ├── address_handler.go
│   │   ├── admin_handler.go
│   │   ├── auth_handler.go
│   │   └── handler.go
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   └── jwt_auth.go
│   ├── models/
│   │   ├── models.go
│   │   └── product.go
│   ├── services/
│   │   ├── email_service.go
│   │   └── shipping_service.go
│   └── utils/
│       ├── jwt.go
│       ├── otp.go
│       ├── password.go
│       └── ...
├── pkg/
│   └── security/
│       └── jwt.go
└── ARCHITECTURE.md
```

## Prinsip

- `cmd/` hanya untuk entrypoint.
- `internal/app` hanya untuk bootstrap dependency dan router.
- `internal/handlers` hanya untuk HTTP request/response.
- `internal/services` hanya untuk business logic.
- `internal/db` hanya untuk koneksi database dan migrasi.
- `internal/models` hanya untuk entity/database model.
- `internal/cache` hanya untuk koneksi cache.
- `internal/middleware` hanya untuk auth dan request policy.
- `pkg/security` hanya untuk helper reusable yang aman dipakai lintas layer.

## Catatan perusahaan

- Hindari duplikasi folder dengan fungsi sama.
- Satu layer = satu tanggung jawab.
- Entry point dibuat tipis.
- Business logic tidak diletakkan di handler.
