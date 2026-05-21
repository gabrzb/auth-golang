# auth-go-gin

A RESTful authentication API built with Go, Gin, JWT, and PostgreSQL. Covers user registration, login, token refresh, and protected routes — designed as a hands-on learning project.

## Stack

- **Go 1.26.3** — language
- **Gin** — HTTP framework
- **GORM + PostgreSQL** — ORM and database
- **golang-jwt/jwt v5** — JWT generation and validation
- **bcrypt** — password hashing
- **godotenv** — environment variable loading
- **Docker + Docker Compose** — containerization

## Running with Docker (recommended)

```bash
cp .env.example .env   # fill in values if needed
docker compose up
```

Both the API and PostgreSQL start together. The API waits for the database healthcheck before booting.

## Running locally (Go only)

Requires a running PostgreSQL instance. Set the connection variables in a `.env` file, then:

```bash
cd backend
go run cmd/api/main.go
```

Check it with:

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
| `JWT_SECRET` | — | Secret for signing tokens |
| `JWT_ACCESS_EXPIRATION` | `15m` | Access token TTL |
| `JWT_REFRESH_EXPIRATION` | `168h` | Refresh token TTL |

## Endpoints

| Method | Route | Auth | Description |
|--------|-------|------|-------------|
| GET | `/health` | No | Health check |
| POST | `/auth/register` | No | Create user |
| POST | `/auth/login` | No | Returns tokens |
| POST | `/auth/refresh` | No | Renew access token |
| GET | `/me` | Yes | Authenticated user data |

## Project structure

```text
├── backend/
│   ├── cmd/api/main.go      # entrypoint
│   ├── internal/
│   │   ├── config/          # env loading
│   │   ├── database/        # GORM connection
│   │   ├── handlers/        # HTTP handlers
│   │   ├── middleware/      # JWT middleware
│   │   ├── models/          # data models
│   │   ├── routes/          # route registration
│   │   └── services/        # business logic + JWT service
│   ├── go.mod
│   └── go.sum
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── .gitignore
```
