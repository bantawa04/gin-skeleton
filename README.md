# Gin Skeleton - Go API Boilerplate

A Go API starter kit built with Gin, GORM, PostgreSQL, and Uber Fx. It provides a domain-oriented project structure, authentication, middleware, migrations, Swagger documentation, linting, API test examples, and placeholder adapters for common external integrations.

> See [ARCHITECTURE.md](./ARCHITECTURE.md) for the architectural patterns and package responsibilities used by the project.

## Features

- **Domain-oriented architecture** with Handler → Service → Repository separation
- **Dependency injection** with Uber Fx
- **JWT authentication** with access and refresh tokens
- **PostgreSQL + GORM** for persistence
- **Goose migrations** with a dedicated migration command
- **HTTP middleware** for CORS, request IDs, logging, sanitization, case conversion, rate limiting, transactions, and centralized error handling
- **Swagger/OpenAPI** documentation
- **Structured logging** with logrus and log rotation
- **Request validation** with go-playground/validator
- **API tests** using Go's `testing` and `net/http/httptest`
- **golangci-lint** configuration and Makefile commands
- **External integration placeholders** for Stripe, AWS S3, and Resend
- **Domain scaffolding** through Makefile templates
- **GitHub Actions PR checks** that run tests and linting for pull requests targeting `main`

## Project Structure

```text
gin-skeleton/
├── .github/
│   └── workflows/
│       └── pr-checks.yml                # Test + lint checks for PRs to main
├── cmd/
│   ├── api/
│   │   └── main.go                    # HTTP API entrypoint
│   └── migrate/
│       └── main.go                    # Goose migration command
│
├── database/
│   └── migrations/                    # SQL migrations
│
├── docker/
│   ├── entrypoint.sh
│   └── web.Dockerfile
│
├── docs/                              # Generated Swagger artifacts
│
├── internal/
│   ├── domain/                        # Business capabilities
│   │   ├── auth/
│   │   ├── health/
│   │   ├── refresh_token/
│   │   └── user/
│   │       ├── handler/
│   │       ├── repository/
│   │       ├── service/
│   │       ├── dto.go
│   │       ├── model.go
│   │       └── request.go
│   │
│   ├── infra/                         # Framework and external adapters
│   │   ├── bootstrap/                 # Fx application/module composition
│   │   ├── config/                    # Environment/runtime configuration
│   │   ├── integration/               # Third-party service adapters
│   │   │   ├── resend/
│   │   │   │   └── client.go
│   │   │   ├── s3/
│   │   │   │   └── client.go
│   │   │   └── stripe/
│   │   │       └── client.go
│   │   ├── logger/
│   │   ├── middleware/
│   │   └── router/
│   │       ├── router.go
│   │       ├── web.go
│   │       └── web_test.go            # Router-level API tests
│   │
│   └── shared/                        # Cross-domain primitives/helpers
│       ├── constant/
│       ├── exception/
│       ├── response/
│       ├── utils/
│       └── validator/
│
├── templates/                         # Domain scaffolding templates
├── .golangci.yml                      # golangci-lint configuration
├── docker-compose.yml
├── env.example
├── go.mod
├── Makefile
└── README.md
```

## Architecture

Business code lives under `internal/domain`, while framework concerns and external systems live under `internal/infra`.

```text
HTTP Request
    │
    ▼
 Handler
    │
    ▼
 Service
    │
    ▼
Repository
    │
    ▼
PostgreSQL
```

External systems follow the same dependency direction:

```text
Domain interface / application need
            │
            ▼
internal/infra/integration/<provider>
            │
            ▼
      External service
```

This keeps vendor-specific SDKs and implementation details outside the business domains.

## External Integrations

Placeholder adapters are included for:

```text
internal/infra/integration/
├── stripe/
├── s3/
└── resend/
```

They intentionally do **not** install Stripe, AWS, or Resend SDKs. Each placeholder contains a small configuration/client shape and returns `ErrNotImplemented` until an application wires in the provider it actually needs.

Typical usage when adapting the skeleton:

1. Define the capability required by the domain, preferably as an interface.
2. Implement that capability under `internal/infra/integration/<provider>`.
3. Add the provider SDK only when the integration is needed.
4. Register the concrete adapter through the Fx bootstrap modules.

This makes it easier to replace providers later, for example Stripe with another payment gateway or S3 with another object-storage implementation.

## Prerequisites

- Go 1.25.0 or higher
- PostgreSQL 12 or higher
- Make (optional but recommended)
- Docker / Docker Compose (optional)

## Installation

```bash
git clone <repository-url>
cd gin-skeleton

go mod download
cp env.example .env
```

Update `.env` with your local configuration, then run migrations:

```bash
make migrate-up
```

Start the API:

```bash
go run ./cmd/api
```

The default API address is:

```text
http://localhost:8000
```

## Common Commands

```bash
make test            # Run all Go tests
make test-cover      # Run tests and generate coverage.out

make lint            # Run golangci-lint
make lint-fix        # Run golangci-lint with automatic fixes
make lint-install    # Install the pinned golangci-lint version

make migrate-up      # Apply pending migrations
make migrate-down    # Roll back the latest migration
make migrate-status  # Show migration status
make migrate-create NAME=add_products
make migrate-baseline
make migrate-fresh   # Drop public schema and re-run migrations

make swagger         # Regenerate Swagger artifacts
make scaffold name=book
```

Run `make help` to see the available Makefile commands.

## Testing

The skeleton includes router-level API tests in:

```text
internal/infra/router/web_test.go
```

The tests use Go's standard `testing` and `net/http/httptest` packages and exercise the real route/middleware/handler path while replacing service/database boundaries with test doubles.

Covered example flows include:

- health check
- signup
- login
- refresh-token rotation
- logout
- protected endpoint rejection without a JWT
- user list
- user lookup
- user update
- user deletion

Run the suite with:

```bash
make test
```

Generate coverage output with:

```bash
make test-cover
```

`coverage.out` is ignored by Git.

## Linting

The project uses `golangci-lint` and includes a repository-level `.golangci.yml`.

Install the pinned version:

```bash
make lint-install
```

Run linting:

```bash
make lint
```

Apply supported automatic fixes:

```bash
make lint-fix
```

## Pull Request Checks

GitHub Actions runs automated validation for every pull request targeting `main` from any source branch.

The workflow lives at:

```text
.github/workflows/pr-checks.yml
```

It runs two independent jobs in parallel:

- `Test` — runs `make test`
- `Lint` — runs `golangci-lint` using the repository configuration

The workflow intentionally handles validation only; no deployment or release workflow is bundled.

## Scaffolding a New Domain

Generate repository, service, and Fx module boilerplate:

```bash
make scaffold name=book
```

This creates:

```text
internal/domain/book/repository/book_repository.go
internal/domain/book/repository/book_repository_interface.go
internal/domain/book/service/book_service.go
internal/domain/book/service/book_service_interface.go
internal/infra/bootstrap/modules/book_module.go
```

You can then add the model, handler, request, and DTO files required by the domain.

## Configuration

Copy `env.example` to `.env` and configure the application.

### Database

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSL_MODE=disable
```

### Server

```env
SERVER_PORT=8000
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
```

### JWT

```env
JWT_SECRET_KEY=your-super-secret-jwt-key-change-in-production
JWT_ACCESS_EXPIRY=168h
JWT_REFRESH_EXPIRY=720h
```

## API Endpoints

### Utility / documentation

```text
GET  /                         Skeleton status information
GET  /health                   Database-aware health check
GET  /api/health               Database-aware API health check
GET  /swagger/*any             Swagger UI (basic auth)
GET  /docs/swagger.yaml        Swagger specification (basic auth)
```

### Authentication

```text
POST /api/auth/signup          Create a user account
POST /api/auth/login           Login and receive access/refresh tokens
POST /api/auth/refresh         Rotate refresh token and issue new tokens
POST /api/auth/logout          Logout; requires an access token
```

### Users

```text
GET    /api/users              List users with pagination
GET    /api/users/:id          Get a user by ID
PUT    /api/users/:id          Update a user; requires JWT
DELETE /api/users/:id          Delete a user; requires JWT
```

## Authentication

Login using:

```text
POST /api/auth/login
```

Send the access token on protected routes:

```http
Authorization: Bearer <access_token>
```

Refresh tokens are rotated through:

```text
POST /api/auth/refresh
```

## Database Migrations

Migrations use [Goose](https://github.com/pressly/goose) and are stored in `database/migrations`.

Migration execution is intentionally separate from normal API startup.

Create a migration:

```bash
make migrate-create NAME=create_products_table
```

Apply migrations:

```bash
make migrate-up
```

Rollback the most recent migration:

```bash
make migrate-down
```

Check migration state:

```bash
make migrate-status
```

## Swagger

Regenerate Swagger artifacts with:

```bash
make swagger
```

Then open:

```text
http://localhost:8000/swagger/index.html
```

## Middleware

The application includes middleware for:

1. CORS
2. Request IDs
3. Structured request logging
4. Input sanitization
5. Request/response case conversion
6. Centralized application error handling
7. Authentication rate limiting
8. JWT authorization
9. Database transactions for write endpoints

## Responses and Error Handling

Successful API responses use a common envelope:

```json
{
  "success": true,
  "message": "operation successful",
  "data": {}
}
```

Validation errors use a Laravel-style field map:

```json
{
  "message": "The given data was invalid.",
  "errors": {
    "email": ["The email field must be a valid email address."]
  }
}
```

Application errors are centralized through the exception middleware and mapped to the appropriate HTTP status code.

## Docker

Build the image:

```bash
docker build -f docker/web.Dockerfile -t gin-skeleton .
```

Run with Docker Compose:

```bash
docker compose up
```

Docker startup does not run migrations by default. Run migrations separately, or opt in where appropriate for your deployment setup.

## Logging

Application logs are written under `logs/` and rotated automatically:

```text
logs/app.log
logs/errors.log
```

The `logs/` directory and log files are ignored by Git.

## Development Guidelines

- Keep HTTP concerns in handlers.
- Keep business rules in services.
- Keep database access in repositories.
- Keep third-party provider code under `internal/infra/integration`.
- Prefer domain-owned interfaces when business logic needs an external capability.
- Use `context.Context` for request-scoped work and cancellation.
- Add tests for new endpoints and business behavior.
- Run `make test` and `make lint` before publishing changes.

## Contributing

1. Create a feature branch.
2. Make the change.
3. Add or update tests where appropriate.
4. Run `make test` and `make lint`.
5. Open a pull request targeting `main`.
6. Ensure the automated `Test` and `Lint` checks pass.

For questions or issues, use the repository's GitHub issue tracker.
