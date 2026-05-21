# auth-go-gin

A RESTful authentication API built with Go, Gin, JWT, and PostgreSQL. Covers user registration, login, token refresh, and protected routes — designed as a hands-on learning project.

## Stack

- **Go 1.22+** — language
- **Gin** — HTTP framework
- **GORM + PostgreSQL** — ORM and database
- **golang-jwt/jwt v5** — JWT generation and validation
- **bcrypt** — password hashing
- **godotenv** — environment variable loading
- **Docker + Docker Compose** — containerization

## Running locally

```bash
cd backend
go run cmd/api/main.go
```

The server starts on port `8080`. Check it with:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

## Environment variables

Copy `.env.example` to `.env` and fill in the values:

```bash
cp .env.example .env
```

## Endpoints

| Method | Route | Auth | Description |
|--------|-------|------|-------------|
| GET | `/health` | No | Health check |
| POST | `/auth/register` | No | Create user |
| POST | `/auth/login` | No | Returns tokens |
| POST | `/auth/refresh` | No | Renew access token |
| GET | `/me` | Yes | Authenticated user data |

## Project structure

```
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
├── .env.example
└── .gitignore
```
