# Base image for both environments
FROM golang:1.25.0-alpine AS base
WORKDIR /app
RUN apk add --no-cache git

# Development stage with Air for hot reloading
FROM base AS development
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8000
CMD ["air", "-c", ".air.toml"]

# Build stage for production
FROM base AS build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Production stage with minimal image
FROM alpine:latest AS production
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/main .
EXPOSE 8000
CMD ["/app/main"]