# API Documentation Guide

## Overview

This project uses **Swagger/OpenAPI** for API documentation, powered by [swaggo/swag](https://github.com/swaggo/swag).

Swagger generation is anchored on the API binary entrypoint at `cmd/api/main.go`. Runtime code is organized under `internal/domain`, `internal/infra`, and `internal/shared`, so generated docs should parse internal packages when scanning annotations.

## Quick Start

### 1. Generate Documentation

```bash
# Using Makefile (recommended)
make swagger

# Or manually
swag init -g cmd/api/main.go -o ./docs --parseDependency --parseInternal
```

### 2. Start the Server

```bash
go run ./cmd/api
```

### 3. Access Swagger UI

Open your browser and navigate to:
```
http://localhost:8000/swagger/index.html
```

## What is Swagger?

**Swagger** (now OpenAPI) is a specification for describing REST APIs. It provides:

- 📖 **Interactive Documentation**: Try API endpoints directly from the browser
- 🔍 **Auto-generated**: Documentation generated from code annotations
- ✅ **Type-safe**: Request/response schemas with validation
- 🔐 **Authentication**: Built-in support for JWT Bearer tokens
- 📱 **Client Generation**: Generate client SDKs in multiple languages

## How It Works

### 1. Code Annotations

Swagger uses special comments (annotations) in your Go code:

```go
// @Summary      Create a new user
// @Description  Create a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      user.UserCreateRequest  true  "User creation data"
// @Success      201   {object}  response.Response{data=user.UserDTO}
// @Failure      422   {object}  response.ErrorResponse
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // Handler implementation
}
```

### 2. Generated Files

After running `make swagger`, the following files are created in `./docs/`:

- `docs.go` - Go code with Swagger definitions
- `swagger.json` - OpenAPI specification (JSON)
- `swagger.yaml` - OpenAPI specification (YAML)

### 3. Swagger UI

The Swagger UI is served at `/swagger/index.html` and automatically reads the generated documentation.

## Project Layout Notes

Swagger annotations should stay close to the HTTP handlers in `internal/domain/<domain>/handler`. Request and response types should be referenced from their domain package, while shared response and error envelopes should come from `internal/shared/response`.

Use `cmd/api/main.go` as the generator entrypoint for local development, CI, and Docker builds:

```bash
swag init -g cmd/api/main.go -o ./docs --parseDependency --parseInternal
```

The top-level `templates/` directory contains scaffolding templates. When adding a new domain from templates, add Swagger annotations to the generated or new handlers before regenerating `docs/docs.go`, `docs/swagger.json`, and `docs/swagger.yaml`.

## Annotation Reference

### Basic Annotations

```go
// @Summary      Short description (required)
// @Description  Detailed description
// @Tags         tag1,tag2        // Group endpoints
// @Accept       json              // Request content type
// @Produce      json              // Response content type
// @Router       /path [method]    // Route definition
```

### Parameters

```go
// Query parameter
// @Param        page  query     int  false  "Page number"  default(1)

// Path parameter
// @Param        id    path      string  true  "User ID"

// Body parameter
// @Param        user  body      user.UserCreateRequest  true  "User data"

// Header parameter
// @Param        Authorization  header    string  true  "Bearer token"
```

### Responses

```go
// Success response
// @Success      200  {object}  response.Response{data=user.UserDTO}

// Error responses
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      422  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
```

### Security

```go
// Require authentication
// @Security     BearerAuth
```

## Current Endpoints Documented

### Authentication (`/api/auth`)
- ✅ `POST /api/auth/login` - User login
- ✅ `POST /api/auth/refresh` - Refresh access token

### Users (`/api/users`)
- ✅ `GET /api/users` - List users (paginated)
- ✅ `GET /api/users/:id` - Get user by ID
- ✅ `PUT /api/users/:id` - Update user (protected)
- ✅ `DELETE /api/users/:id` - Delete user (protected)

### Health (`/health`, `/api/health`)
- ✅ `GET /health` - Health check
- ✅ `GET /api/health` - Health check

## Using Swagger UI

### 1. Try It Out

1. Open `http://localhost:8000/swagger/index.html`
2. Click on an endpoint
3. Click "Try it out"
4. Fill in parameters
5. Click "Execute"
6. See the response!

### 2. Authentication

For protected endpoints:

1. Click the **"Authorize"** button (top right)
2. Enter your JWT token: `Bearer <your-token>`
3. Click "Authorize"
4. Now you can test protected endpoints!

### 3. Example: Login Flow

1. **Login** (`POST /api/auth/login`):
   ```json
   {
     "email": "user@example.com",
     "password": "password123"
   }
   ```
   - Copy the `accessToken` from response

2. **Authorize**:
   - Click "Authorize" button
   - Enter: `Bearer <accessToken>`
   - Click "Authorize"

3. **Test Protected Endpoint**:
   - Try `GET /api/users/:id` or `PUT /api/users/:id`

## Adding Documentation to New Endpoints

### Step 1: Add Annotations

```go
// @Summary      Your endpoint summary
// @Description  Detailed description
// @Tags         your-tag
// @Accept       json
// @Produce      json
// @Param        param  body      yourdomain.YourRequest  true  "Description"
// @Success      200    {object}  response.Response{data=yourdomain.YourDTO}
// @Failure      400    {object}  response.ErrorResponse
// @Router       /your-path [post]
func (h *YourHandler) YourMethod(c *gin.Context) {
    // Implementation
}
```

### Step 2: Regenerate Documentation

```bash
make swagger
```

### Step 3: Verify

- Check `http://localhost:8000/swagger/index.html`
- Your new endpoint should appear!

## Best Practices

### 1. Keep Annotations Updated
- Update annotations when you change endpoints
- Regenerate docs after changes: `make swagger`

### 2. Use Descriptive Summaries
```go
// ✅ Good
// @Summary      Create a new user account

// ❌ Bad
// @Summary      Create user
```

### 3. Document All Parameters
```go
// ✅ Good - All parameters documented
// @Param        page      query     int  false  "Page number"  default(1)
// @Param        per_page  query     int  false  "Items per page"  default(10)

// ❌ Bad - Missing parameters
// @Param        page  query  int  false  "Page"
```

### 4. Use Proper Response Types
```go
// ✅ Good - Specific response type
// @Success      200  {object}  response.Response{data=user.UserDTO}

// ❌ Bad - Generic type
// @Success      200  {object}  map[string]interface{}
```

### 5. Group Related Endpoints
```go
// Use tags to group endpoints
// @Tags         users
// @Tags         auth
// @Tags         admin
```

## Troubleshooting

### Documentation Not Updating?

1. **Regenerate docs**:
   ```bash
   make swagger
   ```

2. **Restart server**:
   ```bash
   go run ./cmd/api
   ```

3. **Clear browser cache** or use incognito mode

### Missing Endpoints?

- Check that annotations are correct
- Verify the handler function is exported (capital letter)
- Ensure the route is registered in router
- Ensure the generator is using `cmd/api/main.go` and `--parseInternal`

### Build Errors?

```bash
# Clean and rebuild
go clean
go mod tidy
go build ./cmd/api
```

### Migration-related Startup Issues?

Swagger generation does not run database migrations. Database schema changes are handled separately with Goose files in `database/migrations`. In Docker, migrations are opt-in; enable the migration mode for environments that should apply schema changes during container startup.

## Alternative Tools

While `swaggo/swag` is the most popular for Gin, other options include:

1. **go-swagger** - More features, steeper learning curve
2. **OpenAPI Generator** - Language-agnostic, more complex setup
3. **Postman Collections** - Manual documentation

For Gin projects, `swaggo/swag` is recommended for its simplicity and Gin integration.

## Resources

- [Swagger/OpenAPI Specification](https://swagger.io/specification/)
- [swaggo/swag Documentation](https://github.com/swaggo/swag)
- [Swagger Annotations Guide](https://github.com/swaggo/swag#declarative-comments-format)

## Summary

✅ **Swagger is set up and ready!**

- Generate docs: `make swagger`
- Access UI: `http://localhost:8000/swagger/index.html`
- Add annotations to new endpoints
- Regenerate after changes
