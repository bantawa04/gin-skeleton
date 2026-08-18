# Agent Guidelines

This file defines repository-specific rules for AI coding agents working on this project. Read `README.md` and `ARCHITECTURE.md` before making structural changes.

## Architecture

- Keep business/domain code under `internal/domain`.
- Keep framework, infrastructure, and third-party implementation details under `internal/infra`.
- Third-party providers belong under `internal/infra/integration/<provider>`.
- Do not call Stripe, S3, Resend, or other vendor SDKs directly from domain code.
- Keep handlers focused on HTTP concerns: request binding, validation, calling services, and returning responses.
- Keep business rules in services.
- Keep database access in repositories.
- Prefer domain-owned interfaces when business logic depends on an external capability.
- Do not create generic `services/`, `helpers/`, `common/`, or similar dumping-ground packages without a clear architectural reason.
- Keep `internal/shared` small and limited to functionality that is genuinely shared across domains.
- Preserve the existing Handler -> Service -> Repository dependency flow unless the task explicitly requires an architectural change.

## JSON Case Conversion

The global `CaseConverterMiddleware` already handles JSON key casing.

- Incoming JSON request keys are normalized from `camelCase` to `snake_case` before handlers bind the request.
- Outgoing JSON response keys are converted recursively from `snake_case` to `camelCase` before they are sent to the client.
- Conversion applies recursively to nested objects and arrays.
- Do not add manual camelCase/snake_case conversion in handlers, services, DTO mappers, or response helpers unless there is a specific exception that the middleware cannot handle.
- Do not remove or bypass the case-conversion middleware for normal JSON API routes without an explicit reason.

## Dependencies

- Do not add a third-party dependency unless it is actually needed.
- Keep placeholder integrations dependency-free until a real provider implementation is required.
- Prefer the Go standard library where practical.
- Avoid introducing overlapping libraries for functionality the repository already provides.

## Database

- Schema changes must use Goose migrations under `database/migrations`.
- Do not auto-migrate schemas during normal API startup.
- Do not modify an existing migration that may already have been applied; create a new migration instead.
- Use the existing transaction middleware for write operations that require transactional behavior.

## API and Middleware

- Follow the existing response envelope and centralized error-handling conventions.
- Use application errors from `internal/shared/exception` instead of ad-hoc error JSON where applicable.
- Protected endpoints must use the existing JWT middleware.
- Reuse existing middleware before creating endpoint-specific alternatives for concerns such as authentication, transactions, rate limiting, sanitization, request IDs, or case conversion.
- Keep route registration under `internal/infra/router`.

## External Integrations

- Put provider-specific implementations under `internal/infra/integration/<provider>`.
- Keep vendor-specific request/response types and SDK details inside the provider adapter where possible.
- Let domain/application code depend on capabilities/interfaces rather than provider SDK types.
- Register concrete integration implementations through the Fx bootstrap wiring when they become active dependencies.

## Testing

- Add or update tests whenever API behavior or business behavior changes.
- Prefer Go's standard `testing` package and `net/http/httptest` for HTTP tests.
- Follow `internal/infra/router/web_test.go` as the reference pattern for router-level API tests.
- Test both successful behavior and important failure/authorization paths when relevant.
- Run:

```bash
make test
```

- Generate coverage when useful with:

```bash
make test-cover
```

## Formatting and Linting

- Format Go code with `gofmt`.
- Run linting before considering a change complete:

```bash
make lint
```

- Use `make lint-fix` only when automatic fixes are appropriate and review the resulting diff.

## Swagger and Generated Files

- Regenerate Swagger artifacts with `make swagger` when API contracts or Swagger annotations change.
- Avoid manually editing generated Swagger files when regeneration can produce the correct output.

## CI

- Pull requests targeting `main` run the repository's GitHub Actions test and lint checks.
- Keep changes compatible with `make test` and the configured `golangci-lint` checks.
- Do not add or replace CI providers unless explicitly requested.

## Scope and Change Discipline

- Avoid unrelated refactoring.
- Preserve existing naming, package boundaries, and conventions unless the requested change requires otherwise.
- Prefer the smallest coherent change that solves the task.
- Update documentation when repository structure, developer workflow, API behavior, or architectural conventions change.
- Do not commit secrets, local environment files, logs, generated coverage output, or OS metadata.
