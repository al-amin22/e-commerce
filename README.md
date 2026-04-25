# E-Commerce Golang + Next.js (Module 1: Identity & Security)

## Stack
- Backend: Gin + GORM + SQLite + Redis
- Frontend: Next.js + Axios Interceptor

## Implemented
- Register + OTP email verification (goroutine send email)
- RBAC (`buyer`, `admin`)
- Double token auth:
  - Access Token: 15 minutes
  - Refresh Token: 7 days
- Refresh Token Rotation with Redis
- Revoke previous refresh token on same device re-login
- HttpOnly cookies for `access_token` and `refresh_token`
- Silent refresh in frontend via Axios interceptor on `401`

## Backend Setup
1. Copy `backend/.env.example` to `backend/.env`
2. Run Redis locally (default `localhost:6379`)
3. Install Go dependencies then run server:
   - `go mod tidy`
   - `go run cmd/server/main.go`

## Frontend Setup
1. Copy `frontend/.env.example` to `frontend/.env.local`
2. Install dependencies and run:
   - `npm install`
   - `npm run dev`

## Important Endpoints
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/verify-email`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me` (protected)

## Notes
- For production set `COOKIE_SECURE=true`
- Set `FRONTEND_ORIGIN` to your frontend domain
