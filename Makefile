# Makefile for Mitho Go Application
include .env
.PHONY: help migrate-create migrate-up migrate-down migrate-down-all migrate-status migrate-force install-migrate

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Migration commands
migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: Please provide a migration name. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir database/migrations $(NAME)

# Database connection string
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

migrate-up: ## Run all pending migrations
	@if [ -z "$(DB_USER)" ] || [ -z "$(DB_PASSWORD)" ] || [ -z "$(DB_HOST)" ] || [ -z "$(DB_PORT)" ] || [ -z "$(DB_NAME)" ] || [ -z "$(DB_SSL_MODE)" ]; then \
		echo "Error: Database environment variables not set. Please check your .env file."; \
		echo "Required: DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME, DB_SSL_MODE"; \
		exit 1; \
	fi
	migrate -path database/migrations -database "$(DB_URL)" up

migrate-down: ## Rollback last migration
	@if [ -z "$(DB_USER)" ] || [ -z "$(DB_PASSWORD)" ] || [ -z "$(DB_HOST)" ] || [ -z "$(DB_PORT)" ] || [ -z "$(DB_NAME)" ] || [ -z "$(DB_SSL_MODE)" ]; then \
		echo "Error: Database environment variables not set. Please check your .env file."; \
		exit 1; \
	fi
	migrate -path database/migrations -database "$(DB_URL)" down 1

migrate-down-all: ## Rollback all migrations
	@if [ -z "$(DB_USER)" ] || [ -z "$(DB_PASSWORD)" ] || [ -z "$(DB_HOST)" ] || [ -z "$(DB_PORT)" ] || [ -z "$(DB_NAME)" ] || [ -z "$(DB_SSL_MODE)" ]; then \
		echo "Error: Database environment variables not set. Please check your .env file."; \
		exit 1; \
	fi
	migrate -path database/migrations -database "$(DB_URL)" down

migrate-status: ## Show migration status
	@if [ -z "$(DB_USER)" ] || [ -z "$(DB_PASSWORD)" ] || [ -z "$(DB_HOST)" ] || [ -z "$(DB_PORT)" ] || [ -z "$(DB_NAME)" ] || [ -z "$(DB_SSL_MODE)" ]; then \
		echo "Error: Database environment variables not set. Please check your .env file."; \
		exit 1; \
	fi
	migrate -path database/migrations -database "$(DB_URL)" version

migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: Please provide a version. Usage: make migrate-force VERSION=1"; \
		exit 1; \
	fi
	@if [ -z "$(DB_USER)" ] || [ -z "$(DB_PASSWORD)" ] || [ -z "$(DB_HOST)" ] || [ -z "$(DB_PORT)" ] || [ -z "$(DB_NAME)" ] || [ -z "$(DB_SSL_MODE)" ]; then \
		echo "Error: Database environment variables not set. Please check your .env file."; \
		exit 1; \
	fi
	migrate -path database/migrations -database "$(DB_URL)" force $(VERSION)

# Utility commands
install-migrate: ## Install migrate CLI via Homebrew
	brew install golang-migrate
