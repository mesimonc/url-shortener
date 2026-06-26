# URL Shortener

A URL shortening service built with Go, Gin, PostgreSQL, and Redis.

## Tech Stack

- **Go 1.26** — Main language
- **Gin** — HTTP framework
- **PostgreSQL** — Primary database
- **Redis** — Caching layer
- **Docker Compose** — Local development environment
- **GORM** — ORM for PostgreSQL

## Features

- `POST /shorten` — Generate a short code for any URL
- `GET /:code` — Redirect to the original URL (301)
- `GET /api/stats/:code` — View click statistics for a short link

## Getting Started

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)

### Run Locally

1. Clone the repository

```bash
git clone https://github.com/mesimonc/url-shortener.git
cd url-shortener
```

2. Start PostgreSQL and Redis

```bash
docker compose up -d
```

3. Copy the environment variables

```bash
cp .env.example .env
```

4. Install dependencies and run

```bash
go mod tidy
go run main.go
```

The server runs on `http://localhost:8080`.

## API Usage

### Shorten a URL

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com"}'
```

Response:
```json
{
  "code": "mB6pCp"
}
```

### Redirect

```
GET http://localhost:8080/mB6pCp
→ 301 Redirect to https://github.com
```

### Get Stats

```bash
curl http://localhost:8080/api/stats/mB6pCp
```

Response:
```json
{
  "code": "mB6pCp",
  "original_url": "https://github.com",
  "clicks": 5,
  "created_at": "2026-06-26T20:10:02Z"
}
```

## Project Structure

```
url-shortener/
├── main.go                        # Entry point
├── config/
│   └── config.go                  # Environment variable loading
├── internal/
│   ├── handler/
│   │   └── url_handler.go         # HTTP handlers
│   ├── service/
│   │   └── url_service.go         # Business logic
│   └── repository/
│       ├── db.go                  # PostgreSQL connection and migration
│       ├── redis.go               # Redis cache
│       └── url_repository.go      # URL database operations
├── docker-compose.yml
└── .env.example
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `:8080` | Server port |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable` | PostgreSQL connection string |
| `REDIS_URL` | `redis://localhost:6379` | Redis connection string |
