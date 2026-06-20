# Base image for both environments
FROM golang:1.25.0-alpine AS base
WORKDIR /app
RUN apk add --no-cache bash git make postgresql-client
ENV PATH="/go/bin:/usr/local/bin:${PATH}"

# Development stage with Air for hot reloading
FROM base AS development
RUN go install github.com/air-verse/air@latest
RUN ln -sf /go/bin/air /usr/local/bin/air
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN chmod +x /app/docker/entrypoint.sh
EXPOSE 8000
ENTRYPOINT ["/app/docker/entrypoint.sh"]
CMD ["air", "-c", ".air.toml"]

# Build stage for production
FROM base AS build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /app/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /app/migrate ./cmd/migrate

# Production stage with minimal image
FROM alpine:latest AS production
RUN apk --no-cache add bash ca-certificates postgresql-client
WORKDIR /app
COPY --from=build /app/api /app/api
COPY --from=build /app/migrate /app/migrate
COPY --from=build /app/docs /app/docs
COPY --from=build /app/docker/entrypoint.sh /app/docker/entrypoint.sh
RUN chmod +x /app/docker/entrypoint.sh
EXPOSE 8000
ENTRYPOINT ["/app/docker/entrypoint.sh"]
CMD ["/app/api"]
