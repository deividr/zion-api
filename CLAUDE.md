# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Run locally with hot reload (Air)
docker-compose up

# Run directly
go run cmd/api/main.go

# Build
go build -o ./tmp/main ./cmd/api/main.go

# Run tests
go test ./...

# Run a single test
go test ./internal/application/use-cases/upload/ -run TestUpload

# Database migrations (requires .env with DATABASE_URL_MIGRATE)
make migrate_up              # Apply all migrations
make migrate_down_last       # Rollback last migration
make create_migration name=description_here  # Create new migration pair

# Data migration from legacy MySQL
make load_products
make load_customers
make load_addresses
make load_orders
```

## Architecture

This is a Go REST API following **Clean Architecture** with strict dependency inversion:

```
HTTP Request → Middleware (Auth/CORS) → Controller → Usecase → Domain Interface → Repository → PostgreSQL
```

**Layer rules:**
- `internal/domain/` — Entities, repository interfaces, and domain errors. Independent of frameworks.
- `internal/usecase/` — Business logic orchestration. Depends only on domain interfaces.
- `internal/controller/` — Gin HTTP handlers. Binds JSON, calls usecases, returns responses.
- `internal/infra/repository/postgres/` — PostgreSQL implementations of domain repository interfaces using pgx + Squirrel.
- `internal/infra/factory/` — Dependency injection (wires repositories → usecases → controllers → routes).
- `internal/middleware/` — JWT auth (Clerk RSA public key validation) and CORS.
- `internal/application/use-cases/` — Complex use cases involving multiple concerns (file uploads, order processing).

**Key conventions:**
- Usecases depend on domain interfaces, never concrete repository types
- All primary keys are UUIDs (`gen_random_uuid()`)
- Soft deletes via `is_deleted BOOLEAN DEFAULT FALSE` column
- All repository/usecase functions take `context.Context` as first parameter
- Squirrel query builder with `squirrel.Dollar` placeholder format for pgx compatibility
- Structured logging with zerolog throughout

**Naming pattern:**
- Domain: `Customer`, `CustomerRepository` (interface)
- Usecase: `CustomerUsecase`
- Controller: `CustomerController`
- Repository: `PgCustomerRepository`

## Tech Stack

- **Framework:** Gin (HTTP), pgx v5 (PostgreSQL driver), Squirrel (SQL builder)
- **Auth:** JWT validation with Clerk RSA public key (`CLERK_PEM_PUBLIC_KEY`)
- **Storage:** AWS SDK v2 targeting Tigris (S3-compatible) for file uploads with presigned URLs
- **Database:** PostgreSQL 16, migrations via golang-migrate
- **Logging:** zerolog with console output
- **Dev tooling:** Air for hot reload, Docker Compose for local env

## Database

- Migrations live in `internal/infra/database/migrations/` as `YYYYMMDDHHMMSS_description.{up,down}.sql` pairs
- Docker Compose runs PostgreSQL on port 5433 (host) → 5432 (container)
- Two DB URLs: `DATABASE_URL` (runtime, used by app) and `DATABASE_URL_MIGRATE` (migrations via Makefile, uses `pgx5://` scheme)

## API Endpoints

All routes are JWT-protected. Resources: products, customers, addresses (nested under customers), categories, orders, pre-signed-url (for uploads).

## Deployment

- Fly.io via GitHub Actions on push to `main`
- Region: `gru` (São Paulo), internal port 8000, auto-scales to zero
- Production build: multi-stage Docker (Go alpine builder → Debian bookworm runtime)

## Environment Variables

See `.env.example` for required variables: `DATABASE_URL`, `DATABASE_URL_MIGRATE`, `CLERK_PEM_PUBLIC_KEY`, `ALLOWED_ORIGINS`, `BUCKET_NAME`, AWS/Tigris credentials.
