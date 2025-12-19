# Architecture & Best Practices

This document outlines the architectural patterns and best practices used in this project.

## Domain-Oriented Structure

We organize code by domain; each domain contains its own handler, requests, DTOs, model, service, and repository.

### Domain folders
- `internal/auth/` – auth handler, requests, DTOs
- `internal/user/` – user handler, DTOs, model, request, service/repository wiring
- `internal/refresh_token/` – refresh token model/DTO + repo/service wiring
- `internal/health/` – health handler
- Cross-cutting: `internal/middleware/`, `internal/response/`, `internal/utils/`, `internal/bootstrap/modules/`

### Layer responsibilities (per domain)
- **Request**: Bind/validate incoming JSON (`binding` tags) within the domain package.
- **DTO**: Response shapes and mappers (`FromUserModel`, etc.) within the domain package.
- **Model**: GORM entities in the domain package (e.g., `internal/user/model.go`).
- **Handler**: HTTP orchestration (bind -> validate -> call service -> respond) in the domain package.
- **Service**: Business logic; calls repositories; no transport logic.
- **Repository**: Data access with GORM; no business rules.

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

## Folder Structure Summary

```
internal/
├── auth/              # Auth domain (handler, requests, DTOs)
├── user/              # User domain (handler, DTOs, request, model)
├── refresh_token/     # Refresh token domain (model, DTO)
├── repository/        # Shared repos (domain-scoped: user/, refresh_token/)
├── service/           # Shared services (domain-scoped: user/, refresh_token/)
├── health/            # Health handler
├── bootstrap/modules/ # Fx modules wiring domains
├── middleware/        # HTTP middleware
├── response/          # Response helpers and formats
├── utils/             # Utilities (binding errors, JWT, etc.)
└── api/exception/     # Error types and middleware
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

