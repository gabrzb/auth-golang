# auth-go-gin

A RESTful authentication API built with Go, Gin, JWT, PostgreSQL, and Redis. Covers user registration, login, token refresh, logout with token invalidation, and protected routes — designed as a hands-on learning project.

## Stack

- **Go 1.26.3** — language
- **Gin** — HTTP framework
- **GORM + PostgreSQL** — ORM and database
- **golang-jwt/jwt v5** — JWT generation and validation
- **bcrypt** — password hashing
- **go-redis/v9** — token blacklist store
- **godotenv** — environment variable loading
- **Docker + Docker Compose** — containerization

## Running with Docker (recommended)

```bash
cp .env.example .env   # fill in values if needed
docker compose up
```

API, PostgreSQL, and Redis start together. The API waits for both database and Redis healthchecks before booting.

## Running locally (Go only)

Requires running PostgreSQL and Redis instances. Set connection variables in a `.env` file, then:

```bash
cd backend
go run cmd/api/main.go
```

Verify it's up:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## Environment variables

Copy `.env.example` to `.env` and adjust as needed:

```bash
cp .env.example .env
```

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DB_HOST` | `localhost` | Postgres host |
| `DB_PORT` | `5432` | Postgres port |
| `DB_USER` | `postgres` | Postgres user |
| `DB_PASSWORD` | `postgres` | Postgres password |
| `DB_NAME` | `auth_db` | Database name |
| `JWT_SECRET` | — | Secret for signing tokens (required) |
| `JWT_ACCESS_EXPIRATION` | `15m` | Access token TTL |
| `JWT_REFRESH_EXPIRATION` | `168h` | Refresh token TTL |
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `COOKIE_SECURE` | `false` | When `true`, the refresh-token cookie is marked `Secure` (HTTPS-only). Keep `false` for `http://localhost`. |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:5173` | Comma-separated browser origins allowed by CORS and the origin-check middleware on `/auth/*`. |

## Auth model

- **Access token** — short-lived JWT returned in the JSON response. The client keeps it in memory and sends it as `Authorization: Bearer <token>`.
- **Refresh token** — long-lived JWT delivered as an **httpOnly, `SameSite=Strict`, `Path=/auth` cookie**. JavaScript can't read it, so XSS can't exfiltrate it; the browser sends it automatically only to `/auth/*` on the same site.
- **Rotation** — every `POST /auth/refresh` blacklists the incoming refresh token and issues a fresh pair. A stolen refresh token has a short useful life.
- **CSRF** — `SameSite=Strict` is the primary defense; an `Origin` allow-list middleware on `/auth/*` is belt-and-suspenders for browser callers.

## Endpoints

| Method | Route | Auth | Description |
|--------|-------|------|-------------|
| GET | `/health` | No | Health check |
| POST | `/auth/register` | No | Create user |
| POST | `/auth/login` | No | Returns access token in body, sets refresh-token cookie |
| POST | `/auth/refresh` | Cookie | Reads refresh cookie, rotates pair, sets new cookie |
| POST | `/auth/logout` | Yes | Blacklists both tokens, clears refresh cookie |
| GET | `/me` | Yes | Authenticated user data |

## Usage examples

### Health check

```bash
curl http://localhost:8080/health
```

```json
{"status": "ok"}
```

### Register

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret123"}'
```

```json
{"id": 1, "email": "user@example.com", "created_at": "...", "updated_at": "..."}
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -c cookies.txt \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:5173" \
  -d '{"email": "user@example.com", "password": "secret123"}'
```

Response body — only the access token:

```json
{"access_token": "<jwt>", "expires_in": 900}
```

Response also carries:

```
Set-Cookie: refresh_token=<jwt>; Path=/auth; Max-Age=604800; HttpOnly; SameSite=Strict
```

The refresh token never appears in JSON.

### Access protected route

```bash
curl http://localhost:8080/me \
  -H "Authorization: Bearer <access_token>"
```

```json
{"id": 1, "email": "user@example.com", "created_at": "...", "updated_at": "..."}
```

### Refresh access token

No body — the refresh token rides on the cookie set by `/auth/login`. Each call rotates the refresh token (old one is blacklisted, new one is written back to the cookie).

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -b cookies.txt -c cookies.txt \
  -H "Origin: http://localhost:5173"
```

```json
{"access_token": "<new_jwt>", "expires_in": 900}
```

### Logout

Blacklists both tokens (reading the refresh from the cookie) and clears the refresh cookie. No body required.

```bash
curl -X POST http://localhost:8080/auth/logout \
  -b cookies.txt -c cookies.txt \
  -H "Authorization: Bearer <access_token>" \
  -H "Origin: http://localhost:5173"
```

```json
{"message": "logged out"}
```

## Error responses

All errors follow the same shape:

```json
{"error": "description"}
```

Common HTTP status codes:
- `400` — missing or invalid request body
- `401` — missing, invalid, expired, or revoked token / wrong credentials
- `409` — email already in use
- `500` — unexpected server error

## Project structure

```text
├── backend/
│   ├── cmd/api/main.go      # entrypoint and dependency wiring
│   ├── internal/
│   │   ├── config/          # env loading
│   │   ├── database/        # GORM connection and migration
│   │   ├── handlers/        # HTTP handlers (auth, user)
│   │   ├── middleware/      # JWT auth middleware
│   │   ├── models/          # data models
│   │   ├── routes/          # route registration
│   │   ├── services/        # business logic and JWT service
│   │   └── store/           # Redis token blacklist
│   ├── go.mod
│   └── go.sum
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── .gitignore
```