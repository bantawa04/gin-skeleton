# Architecture & Best Practices

This document outlines the architectural patterns and best practices used in this project.

## Domain-Oriented Structure

Runtime code is organized behind one executable entrypoint and three internal package groups:

- `cmd/api` owns the HTTP API binary entrypoint.
- `internal/domain` owns business capabilities and domain-specific HTTP/application/data code.
- `internal/infra` owns framework, database, configuration, logging, routing, middleware, and dependency-injection adapters.
- `internal/shared` owns cross-domain primitives that are intentionally reusable.
- `templates` owns top-level scaffolding templates so generated-code assets are not mixed into runtime packages.

### Domain folders
- `internal/domain/auth/` - auth handler, requests, DTOs, and service logic.
- `internal/domain/user/` - user handler, requests, DTOs, model, service, and repository.
- `internal/domain/refresh_token/` - refresh-token model, DTOs, service, and repository.
- `internal/domain/health/` - health-check handler.

### Layer responsibilities (per domain)
- **Request**: Bind/validate incoming JSON (`binding` tags) within the domain package.
- **DTO**: Response shapes and mappers (`FromUserModel`, etc.) within the domain package.
- **Model**: GORM entities in the domain package (for example, `internal/domain/user/model.go`).
- **Handler**: HTTP orchestration (bind -> validate -> call service -> respond) in the domain package.
- **Service**: Business logic; calls repositories; no transport logic.
- **Repository**: Data access with GORM; no business rules.

### Infrastructure folders
- `internal/infra/bootstrap/` - Uber Fx application construction and module composition.
- `internal/infra/config/` - environment loading and typed runtime configuration.
- `internal/infra/logger/` - logging adapter setup.
- `internal/infra/middleware/` - Gin middleware.
- `internal/infra/router/` - route registration and API grouping.

### Shared folders
- `internal/shared/constant/` - constants used by multiple domains.
- `internal/shared/exception/` - application error types and constructors.
- `internal/shared/response/` - response envelope helpers.
- `internal/shared/utils/` - generic helpers such as token and binding utilities.
- `internal/shared/validator/` - validator setup and validation helpers.

### Handler pattern (current)
```go
var req user.UserUpdateRequest
if err := c.ShouldBindJSON(&req); err != nil {
    if ve := utils.ExtractBindingErrors(err); len(ve) > 0 {
        _ = c.Error(exceptions.ValidationError("The given data was invalid.", nil, ve))
        return
    }
    _ = c.Error(exceptions.ValidationError("Invalid request format. Please check your JSON syntax.", nil))
    return
}

updatedUser, err := h.userService.UpdateUser(ctx, updates, req.Password, id)
if err != nil {
    _ = c.Error(err)
    return
}

response.SendResponse(c, user.FromUserModel(*updatedUser), "user updated successfully")
```

## Error Handling Pattern

### Centralized Error Handling
- **Handler**: Use `c.Error(appErr)` to pass errors to middleware
- **Middleware**: Error handler middleware formats all errors consistently
- **Benefits**:
  - Consistent error format across all endpoints
  - Centralized logging
  - Cleaner handler code
  - Easier to maintain

**Example:**
```go
// In handler
if err := h.validator.Struct(req); err != nil {
    validationErrors := h.validator.GenerateValidationErrors(err)
    appErr := exceptions.ValidationError("Validation failed", nil, validationErrors)
    _ = c.Error(appErr)  // Pass to error middleware
    return
}
```

### Custom Error Types
```go
exceptions.ValidationError()   // 422 Unprocessable Entity
exceptions.NotFoundError()     // 404 Not Found
exceptions.UnauthorizedError() // 401 Unauthorized
exceptions.InternalError()     // 500 Internal Server Error
```

## Validation Pattern

### Two-Level Validation
1. **Gin Binding**: Basic JSON structure validation (`binding:"required"`)
2. **Custom Validator**: Business rules validation (`validate:"min=2,max=100"`)

**Example:**
```go
type UserCreateRequest struct {
    Name  string `json:"name" binding:"required" validate:"required,min=2,max=100"`
    Email string `json:"email" binding:"required,email" validate:"required,email"`
}
```

## Response Pattern

### Consistent Response Format
All successful responses use:
```go
response.SendResponse(c, data, message)
```

Output:
```json
{
    "success": true,
    "message": "user created successfully",
    "data": { ... }
}
```

All error responses are formatted by error middleware:
```json
{
    "success": false,
    "message": "Validation failed",
    "errors": { ... }
}
```

## Key Architectural Decisions

### 1. Request vs DTO Separation
- **Why**: Clear separation between API contract (Request) and internal data transfer (DTO)
- **Benefit**: API changes don't affect internal layers

### 2. Error Middleware Pattern
- **Why**: Centralized error handling and formatting
- **Benefit**: Consistent error responses, easier maintenance, cleaner handlers

### 3. Generic Transformers
- **Why**: Reusable transformation logic
- **Benefit**: DRY principle, type-safe transformations

### 4. Transaction Middleware
- **Why**: Automatic transaction management for write operations
- **Benefit**: Data consistency, easier to use

### 5. Explicit API Entrypoint
- **Why**: `cmd/api/main.go` makes the deployable binary explicit and keeps startup code outside reusable packages.
- **Benefit**: Clear build target for local development, Docker, Swagger generation, and future additional commands.

### 6. Infrastructure Isolation
- **Why**: Gin, GORM, logging, configuration, and Fx wiring are framework concerns, not domain concerns.
- **Benefit**: Domains stay focused on application behavior while adapters remain replaceable and easier to test.

### 7. Goose Migrations
- **Why**: Goose provides a simple SQL-first migration workflow with timestamped files and explicit up/down control.
- **Benefit**: Schema changes can be reviewed, run locally, and run in deployment without coupling migrations to app boot.

### 8. Opt-in Container Migrations
- **Why**: Containers should normally start the API, not silently mutate the database.
- **Benefit**: Deployments can choose when to run migrations by enabling the migration mode explicitly.

## Folder Structure Summary

```
cmd/
└── api/
    └── main.go          # API binary entrypoint

internal/
├── domain/
│   ├── auth/
│   ├── user/
│   ├── refresh_token/
│   └── health/
├── infra/
│   ├── bootstrap/
│   ├── config/
│   ├── logger/
│   ├── middleware/
│   └── router/
└── shared/
    ├── constant/
    ├── exception/
    ├── response/
    ├── utils/
    └── validator/

database/
└── migrations/        # Goose SQL migrations

templates/             # Top-level scaffolding templates
```

## Best Practices

1. ✅ **Use Request structs** for HTTP input binding
2. ✅ **Use DTOs** for HTTP responses and internal data transfer
3. ✅ **Use Models** for database operations
4. ✅ **Use error middleware** for consistent error handling
5. ✅ **Validate at handler level** before calling services
6. ✅ **Keep handlers thin** - just coordinate between layers
7. ✅ **Business logic in services** - not in handlers or repositories
8. ✅ **Database operations in repositories** - not in services
9. ✅ **Use dependency injection** (Uber Fx) for all components
10. ✅ **Use context** for request-scoped data (user ID, request ID, transactions)
11. ✅ **Keep migration execution explicit** with Goose commands or opt-in container migration mode
12. ✅ **Keep scaffolding templates top-level** under `templates/`, outside runtime packages
