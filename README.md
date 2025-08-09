# trade_company

Scaffold for a BizBuySell-like project.

## Stack
- Go 1.22, Gin, GORM (MySQL), Redis, JWT, Zap
- REST API (v1), GraphQL (todo), Wire (todo)
- Templates with Tailwind CDN
- Docker Compose (mysql, redis, app, adminer)

## Getting started
1. Copy `.env.example` to `.env` if needed.
2. Start MySQL and Redis via Docker:
   - `docker compose up -d`
3. Run locally:
   - `go mod tidy`
   - `go run ./cmd/server`
4. Build:
   - `make build`

## REST endpoints
- POST `/api/v1/auth/register` { email, password }
- POST `/api/v1/auth/login` { email, password }
- GET `/api/v1/listings`
- GET `/api/v1/listings/:id`
- POST `/api/v1/listings` (Bearer token)

## Todo
- GraphQL schema and resolvers (gqlgen)
- Wire DI providers
- Cache listing search in Redis with TTL and invalidation
- More pages: integrate server-side templates
- Robust migrations and seeds
