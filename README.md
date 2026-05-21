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

## Endpoints

| Method | Route | Auth | Description |
|--------|-------|------|-------------|
| GET | `/health` | No | Health check |
| POST | `/auth/register` | No | Create user |
| POST | `/auth/login` | No | Returns access + refresh tokens |
| POST | `/auth/refresh` | No | Renew access token |
| POST | `/auth/logout` | Yes | Invalidate both tokens |
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
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret123"}'
```

```json
{
  "access_token": "<jwt>",
  "refresh_token": "<jwt>",
  "expires_in": 900
}
```

### Access protected route

```bash
curl http://localhost:8080/me \
  -H "Authorization: Bearer <access_token>"
```

```json
{"id": 1, "email": "user@example.com", "created_at": "...", "updated_at": "..."}
```

### Refresh access token

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "<refresh_token>"}'
```

```json
{"access_token": "<new_jwt>", "expires_in": 900}
```

### Logout

Invalidates both the access token and refresh token. After this, neither token will be accepted.

```bash
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "<refresh_token>"}'
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