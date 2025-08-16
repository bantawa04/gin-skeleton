# Makefile for Mitho Go Application
include .env
.PHONY: help migrate-create migrate-up migrate-down migrate-down-all migrate-fresh migrate-status migrate-force install-migrate

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

migrate-fresh: ## Drop all tables and rerun all migrations (fresh start)
	@if [ -z "$(DB_USER)" ] || [ -z "$(DB_PASSWORD)" ] || [ -z "$(DB_HOST)" ] || [ -z "$(DB_PORT)" ] || [ -z "$(DB_NAME)" ] || [ -z "$(DB_SSL_MODE)" ]; then \
		echo "Error: Database environment variables not set. Please check your .env file."; \
		exit 1; \
	fi
	@echo "⚠️  WARNING: This will delete ALL data in the database!"
	@echo "Database: $(DB_NAME) on $(DB_HOST):$(DB_PORT)"
	@read -p "Are you sure you want to continue? (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "Proceeding with fresh migration..."; \
		echo "Rolling back all migrations..."; \
		migrate -path database/migrations -database "$(DB_URL)" down -all; \
		echo "Running all migrations..."; \
		migrate -path database/migrations -database "$(DB_URL)" up; \
		echo "✅ Fresh migration completed successfully!"; \
	else \
		echo "❌ Fresh migration cancelled."; \
		exit 1; \
	fi

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
install-migrate: ## Install migrate CLI based on OS
	@echo "Detecting OS and installing migrate CLI..."
	@if [ "$(OS)" = "Windows_NT" ] || [ "$(shell uname -s)" = "MINGW64_NT" ] || [ "$(shell uname -s)" = "MSYS_NT" ]; then \
		echo "Windows detected, installing via scoop..."; \
		scoop install migrate; \
	elif [ "$(shell uname -s)" = "Darwin" ]; then \
		echo "macOS detected, installing via Homebrew..."; \
		brew install golang-migrate; \
	elif [ "$(shell uname -s)" = "Linux" ]; then \
		echo "Linux detected, installing via package manager..."; \
		if command -v apt-get >/dev/null 2>&1; then \
			echo "Debian/Ubuntu detected, installing via apt..."; \
			curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | sudo apt-key add -; \
			echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $$(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/migrate.list; \
			sudo apt-get update; \
			sudo apt-get install -y migrate; \
		elif command -v yum >/dev/null 2>&1; then \
			echo "RHEL/CentOS detected, installing via yum..."; \
			curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | sudo rpm --import -; \
			echo "[migrate]" | sudo tee /etc/yum.repos.d/migrate.repo; \
			echo "name=migrate" | sudo tee -a /etc/yum.repos.d/migrate.repo; \
			echo "baseurl=https://packagecloud.io/golang-migrate/migrate/el/$$(rpm -E %rhel)/$$(rpm -E %dist)/" | sudo tee -a /etc/yum.repos.d/migrate.repo; \
			echo "gpgcheck=1" | sudo tee -a /etc/yum.repos.d/migrate.repo; \
			echo "gpgkey=https://packagecloud.io/golang-migrate/migrate/gpgkey" | sudo tee -a /etc/yum.repos.d/migrate.repo; \
			sudo yum install -y migrate; \
		else \
			echo "Unsupported Linux distribution. Please install manually."; \
			exit 1; \
		fi; \
	else \
		echo "Unsupported OS. Please install migrate CLI manually."; \
		exit 1; \
	fi
