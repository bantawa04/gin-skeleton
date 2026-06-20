# Gin Skeleton - Go API Boilerplate

A production-ready Go API boilerplate built with Gin, GORM, and Uber FX. This skeleton provides a clean architecture, dependency injection, comprehensive middleware, and best practices for building scalable REST APIs.

> 📖 **See [ARCHITECTURE.md](./ARCHITECTURE.md)** for detailed architectural patterns and best practices.

## Features

- 🏗️ **Clean Architecture**: Layered architecture with Handler → Service → Repository pattern
- 🔌 **Dependency Injection**: Uber FX for dependency management and lifecycle
- 🔐 **Authentication**: JWT-based authentication with access and refresh tokens
- 🛡️ **Security**: Input sanitization, CORS, rate limiting, XSS protection
- 📊 **Database**: PostgreSQL with GORM, Goose migrations, and transaction support
- 📝 **Logging**: Structured logging with logrus and file rotation
- ✅ **Validation**: Request validation with go-playground/validator
- 🔄 **Case Conversion**: Automatic camelCase ↔ snake_case conversion
- 🏥 **Health Checks**: Database connectivity health endpoint
- 🎯 **Request Tracing**: Request ID middleware for distributed tracing
- ⚡ **Performance**: Connection pooling, optimized middleware, buffer pooling
- 📚 **API Documentation**: Interactive Swagger/OpenAPI documentation

## Project Structure

```
gin-skeleton/
├── cmd/
│   └── api/
│       └── main.go              # API binary entrypoint
├── database/
│   └── migrations/              # Goose SQL migrations
├── docker/
│   └── web.Dockerfile           # Docker configuration
├── internal/
│   ├── domain/                  # Business domains (auth, user, refresh token, health)
│   │   ├── auth/                # Auth handlers, requests, DTOs, service
│   │   ├── user/                # User handlers, requests, DTOs, model, service, repository
│   │   ├── refresh_token/       # Refresh token model, service, repository
│   │   └── health/              # Health handler
│   ├── infra/                   # Infrastructure adapters and wiring
│   │   ├── bootstrap/           # Uber FX application and module composition
│   │   ├── config/              # Environment and runtime configuration
│   │   ├── database/            # GORM connection and database setup
│   │   ├── logger/              # Logging adapter
│   │   ├── middleware/          # HTTP middleware
│   │   └── router/              # Gin route registration
│   └── shared/                  # Cross-domain primitives and helpers
│       ├── constant/            # Shared constants
│       ├── exception/           # Error types and error helpers
│       ├── response/            # Response envelopes
│       ├── utils/               # Shared utilities
│       └── validator/           # Validation helpers
├── templates/                   # Top-level scaffolding templates
├── docker-compose.yml           # Docker Compose configuration
├── env.example                  # Environment variables template
├── go.mod                       # Go module dependencies
├── Makefile                     # Build and migration commands
└── README.md                    # This file
```

## Architecture

The application follows a clean architecture pattern:

```
┌─────────────────────────────────────────┐
│         HTTP Handlers (API Layer)        │
│  - Request validation                    │
│  - Response formatting                   │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Services (Business Logic)          │
│  - Business rules                       │
│  - Transaction management               │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│    Repositories (Data Access)          │
│  - Database operations                  │
│  - Query optimization                   │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Database (PostgreSQL)           │
└─────────────────────────────────────────┘
```

## Prerequisites

- Go 1.25.0 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd gin-skeleton
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

4. **Run database migrations**
   ```bash
   make migrate-up
   ```

5. **Start the application**
   ```bash
   go run ./cmd/api
   ```

The API will be available at `http://localhost:8000`

## Scaffolding a new domain

Generate repository, service, and Fx module wiring from top-level templates:

```bash
make scaffold name=book
```

This creates:
- `internal/domain/book/repository/book_repository.go` (+ interface)
- `internal/domain/book/service/book_service.go` (+ interface)
- `internal/infra/bootstrap/modules/book_module.go`

Notes:
- The `name` argument is converted to lower-case for packages and PascalCase for types.
- Templates live under `templates/` so scaffolding assets are not mixed into runtime packages.
- You still need to add the model, handlers, requests, and DTOs under `internal/domain/<name>/` when needed.

## Configuration

The application uses environment variables for configuration. Copy `env.example` to `.env` and configure:

### Database Configuration
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSL_MODE=disable
```

### Server Configuration
```env
SERVER_PORT=8000
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
```

### JWT Configuration
```env
JWT_SECRET_KEY=your-super-secret-jwt-key-change-in-production
JWT_ACCESS_EXPIRY=168h      # 7 days
JWT_REFRESH_EXPIRY=720h     # 30 days
```

## API Documentation

### Swagger UI

Interactive API documentation is available at:
```
http://localhost:8000/swagger/index.html
```

The Swagger UI provides:
- Complete API endpoint documentation
- Request/response schemas
- Try-it-out functionality
- Authentication support (JWT Bearer tokens)

### Generate Documentation

```bash
# Generate Swagger documentation
make swagger

# Or manually
swag init -g cmd/api/main.go -o ./docs --parseDependency --parseInternal
```

### API Endpoints

#### Public Endpoints

- `GET /ping` - Health check (simple)
- `GET /health` - Health check with database connectivity
- `GET /swagger/*any` - Swagger API documentation (basic auth)
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh access token
- `GET /api/users` - List users (paginated)
- `GET /api/users/:id` - Get user by ID

#### Protected Endpoints (Require JWT)

- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

## Database Migrations

Migrations use [Goose](https://github.com/pressly/goose) with SQL files stored in `database/migrations`.

The recommended workflow uses the embedded migration runner:

```bash
make migrate-up
```

### Create a new migration
```bash
make migrate-create NAME=create_products_table
```

### Run migrations
```bash
make migrate-up
```

### Rollback last migration
```bash
make migrate-down
```

### Check migration status
```bash
make migrate-status
```

### Baseline an existing database
```bash
make migrate-baseline
```

### Fresh migration (drop all and rerun)
```bash
make migrate-fresh
```

## Middleware

The application includes the following middleware (executed in order):

1. **CORS Middleware** - Handles Cross-Origin Resource Sharing
2. **Request ID Middleware** - Generates unique request IDs for tracing
3. **Logging Middleware** - Structured request/response logging
4. **Sanitize Middleware** - XSS prevention through input sanitization
5. **Case Converter Middleware** - Converts camelCase ↔ snake_case
6. **Error Handler Middleware** - Centralized error handling
7. **Rate Limiting** - Applied to authentication endpoints (10 req/min)
8. **JWT Auth Middleware** - Validates JWT tokens for protected routes

## Authentication

The application uses JWT tokens for authentication:

1. **Login**: `POST /api/auth/login` with email and password
   - Returns access token and refresh token
   
2. **Protected Routes**: Include token in Authorization header
   ```
   Authorization: Bearer <access_token>
   ```

3. **Refresh Token**: `POST /api/auth/refresh` with refresh token
   - Returns new access token

## Development

### Running in Development Mode
```bash
go run ./cmd/api
```

### Building
```bash
go build -o bin/api ./cmd/api
```

### Running Tests
```bash
go test ./...
```

### Generate API Documentation
```bash
make swagger
# Or
swag init -g cmd/api/main.go -o ./docs --parseDependency --parseInternal
```

After generating, access the documentation at `http://localhost:8000/swagger/index.html`

### Code Formatting
```bash
go fmt ./...
```

### Linting
```bash
golangci-lint run
```

## Docker

### Build Docker Image
```bash
docker build -f docker/web.Dockerfile -t gin-skeleton .
```

### Run with Docker Compose
```bash
docker compose up
```

### Opt-in Docker migrations

Docker startup does not run migrations by default. Apply migrations explicitly before starting the app, or opt in with the migration mode provided by the deployment configuration:

```bash
RUN_MIGRATIONS=true docker compose up
```

Keep automatic migrations disabled for normal application startup unless the environment is prepared for schema changes during container boot.

## Logging

Logs are written to the `logs/` directory:
- `app.log` - General application logs
- `errors.log` - Error logs only

Logs are rotated daily and kept for 30 days. Logs are also rotated when they exceed 100MB.

## Error Handling

The application uses structured error handling:

- **Validation Errors** (422): Invalid input data
- **Not Found** (404): Resource not found
- **Unauthorized** (401): Authentication required
- **Forbidden** (403): Insufficient permissions
- **Internal Error** (500): Server errors

All errors follow a consistent format:
```json
{
  "success": false,
  "message": "Error message",
  "description": "Detailed description",
  "data": {}
}
```

## Security Features

- **Password Hashing**: bcrypt with default cost
- **JWT Tokens**: HS256 signing with configurable expiry
- **Input Sanitization**: XSS prevention via bluemonday
- **Rate Limiting**: Protection against brute force attacks
- **CORS**: Configurable cross-origin resource sharing
- **Request ID**: Distributed tracing support

## Best Practices

1. **Always use context.Context** for cancellation and timeouts
2. **Use transactions** for multi-step database operations
3. **Validate all inputs** using the validator package
4. **Handle errors properly** using the exception package
5. **Use structured logging** with appropriate log levels
6. **Follow the repository pattern** for data access
7. **Keep business logic in services**, not handlers

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the MIT License.

## Support

For issues and questions, please open an issue on the repository.
