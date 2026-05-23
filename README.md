# auth-go-gin

A full-stack authentication app built with Go, Gin, JWT, PostgreSQL, Redis, React, and Caddy. It covers user registration, login, cookie-based refresh token rotation, logout with token invalidation, protected routes, and a production-style SPA reverse proxy.

## Stack

- **Go 1.26.3** — language
- **Gin** — HTTP framework
- **GORM + PostgreSQL** — ORM and database
- **golang-jwt/jwt v5** — JWT generation and validation
- **bcrypt** — password hashing
- **go-redis/v9** — token blacklist store
- **godotenv** — environment variable loading
- **React + Vite** — frontend SPA
- **Caddy** — production static file server and `/api` reverse proxy
- **Docker + Docker Compose** — containerization

## Running with Docker (recommended)

```bash
cp .env.example .env   # fill in values if needed
docker compose up --build
```

The SPA, API, PostgreSQL, and Redis start together. Open `http://localhost`; Caddy serves the frontend and proxies `/api/*` to the API on the internal Docker network.

Verify it's up:

```bash
curl http://localhost/api/health
# {"status":"ok"}
```

## Running locally for development

Requires running PostgreSQL and Redis instances. Set backend connection variables in a `.env` file, then start the API:

```bash
cd backend
go run cmd/api/main.go
```

Verify it's up:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

In another shell, start the Vite frontend:

```bash
cd frontend
npm run dev
```

Open `http://localhost:5173`. Vite proxies `/api/*` to `http://localhost:8080`.

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
| `CORS_ALLOWED_ORIGINS` | `http://localhost:5173,http://127.0.0.1:5173,http://localhost,http://127.0.0.1` | Comma-separated browser origins allowed by CORS and the origin-check middleware on `/auth/*`. |
| `WEB_PORT` | `80` | Host port mapped to the Caddy `web` container's port 80. |

## Production deployment

Use the production overlay to serve the SPA through Caddy and keep the API private to the Docker network. The overlay requires an explicit production origin because `COOKIE_SECURE=true` is enabled.

```bash
CORS_ALLOWED_ORIGINS=https://yourdomain.com \
docker compose -f docker-compose.yml -f docker-compose.prod.yml up --build
```

PowerShell:

```powershell
$env:CORS_ALLOWED_ORIGINS = "https://yourdomain.com"
docker compose -f docker-compose.yml -f docker-compose.prod.yml up --build
```

This publishes only `web` on `WEB_PORT` and keeps `api`, `db`, and `redis` internal. Browser requests go to `/api/...`; Caddy strips `/api` and forwards to `api:8080`. Use the base `docker compose up --build` command for plain HTTP local browser testing.

## Auth model

- **Access token** — short-lived JWT returned in the JSON response. The client keeps it in memory and sends it as `Authorization: Bearer <token>`.
- **Refresh token** — long-lived JWT delivered as an **httpOnly, `SameSite=Strict`, `Path=/auth` cookie**. JavaScript can't read it, so XSS can't exfiltrate it; the browser sends it automatically only to `/auth/*` on the same site.
- **Rotation** — every `POST /auth/refresh` blacklists the incoming refresh token and issues a fresh pair. A stolen refresh token has a short useful life.
- **CSRF** — `SameSite=Strict` is the primary defense; an `Origin` allow-list middleware on `/auth/*` is belt-and-suspenders for browser callers.

## Endpoints

The public Docker entrypoint is Caddy, so API calls use `/api/...`. Caddy removes `/api` before forwarding to the Go API.

| Method | Public route | Backend route | Auth | Description |
|--------|--------------|---------------|------|-------------|
| GET | `/api/health` | `/health` | No | Health check |
| POST | `/api/auth/register` | `/auth/register` | No | Create user |
| POST | `/api/auth/login` | `/auth/login` | No | Returns access token in body, sets refresh-token cookie |
| POST | `/api/auth/refresh` | `/auth/refresh` | Cookie | Reads refresh cookie, rotates pair, sets new cookie |
| POST | `/api/auth/logout` | `/auth/logout` | Yes | Blacklists both tokens, clears refresh cookie |
| GET | `/api/me` | `/me` | Yes | Authenticated user data |

## Usage examples

### Health check

```bash
curl http://localhost/api/health
```

```json
{"status": "ok"}
```

### Register

```bash
curl -X POST http://localhost/api/auth/register \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost" \
  -d '{"email": "user@example.com", "password": "secret123"}'
```

```json
{"id": 1, "email": "user@example.com", "created_at": "...", "updated_at": "..."}
```

### Login

```bash
curl -X POST http://localhost/api/auth/login \
  -c cookies.txt \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost" \
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
curl http://localhost/api/me \
  -H "Authorization: Bearer <access_token>"
```

```json
{"id": 1, "email": "user@example.com", "created_at": "...", "updated_at": "..."}
```

### Refresh access token

No body — the refresh token rides on the cookie set by `/auth/login`. Each call rotates the refresh token (old one is blacklisted, new one is written back to the cookie).

```bash
curl -X POST http://localhost/api/auth/refresh \
  -b cookies.txt -c cookies.txt \
  -H "Origin: http://localhost"
```

```json
{"access_token": "<new_jwt>", "expires_in": 900}
```

### Logout

Blacklists both tokens (reading the refresh from the cookie) and clears the refresh cookie. No body required.

```bash
curl -X POST http://localhost/api/auth/logout \
  -b cookies.txt -c cookies.txt \
  -H "Authorization: Bearer <access_token>" \
  -H "Origin: http://localhost"
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
├── frontend/
│   ├── src/                 # React SPA
│   ├── Caddyfile            # SPA + /api reverse proxy
│   ├── Dockerfile           # production frontend image
│   └── package.json
├── Dockerfile               # API image
├── docker-compose.yml
├── docker-compose.prod.yml
├── .env.example
└── .gitignore
```
