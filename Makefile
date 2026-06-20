# Makefile for Gin Skeleton Application
-include .env

.PHONY: help migrate-up migrate-down migrate-status migrate-create migrate-baseline migrate-fresh swagger scaffold

_CLEAN = find database/migrations -name "._*" -delete 2>/dev/null; true
_MIGRATE_ENV = DB_USER='$(DB_USER)' DB_PASSWORD='$(DB_PASSWORD)' DB_HOST='$(DB_HOST)' DB_PORT='$(DB_PORT)' DB_NAME='$(DB_NAME)' DB_SSL_MODE='$(DB_SSL_MODE)'

help: ## Show available commands
	@grep -Eh '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'

migrate-up: ## Apply all pending migrations
	@$(_CLEAN); \
	$(_MIGRATE_ENV) go run ./cmd/migrate up; \
	echo "Migrations applied."

migrate-down: ## Roll back the most recent migration
	@$(_CLEAN); \
	$(_MIGRATE_ENV) go run ./cmd/migrate down

migrate-status: ## Show current migration status
	@$(_CLEAN); \
	$(_MIGRATE_ENV) go run ./cmd/migrate status

migrate-create: ## Create a new SQL migration (usage: make migrate-create NAME=add_foo)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: usage: make migrate-create NAME=add_foo"; exit 1; \
	fi; \
	$(_CLEAN); \
	$(_MIGRATE_ENV) go run ./cmd/migrate create "$(NAME)" sql; \
	$(_CLEAN); \
	echo "Migration file created in database/migrations/"

migrate-baseline: ## Mark embedded migrations as applied without running them
	@$(_CLEAN); \
	$(_MIGRATE_ENV) go run ./cmd/migrate baseline; \
	$(_CLEAN)

migrate-fresh: ## Drop public schema and re-apply all migrations from scratch (destructive)
	@echo "WARNING: This will delete ALL data in $(DB_NAME) on $(DB_HOST):$(DB_PORT)."; \
	printf "Continue? (y/N): "; read confirm; \
	[ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ] || { echo "Cancelled."; exit 1; }; \
	PGPASSWORD='$(DB_PASSWORD)' psql -h '$(DB_HOST)' -p '$(DB_PORT)' -U '$(DB_USER)' -d '$(DB_NAME)' \
		-c "DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;" >/dev/null; \
	$(_CLEAN); \
	$(_MIGRATE_ENV) go run ./cmd/migrate up; \
	echo "Fresh migration completed."

swagger: ## Regenerate Swagger artifacts from code annotations
	@go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g ./cmd/api/main.go -o ./docs --parseDependency --parseInternal
	@echo "Swagger docs regenerated."
	@echo "Access docs at: http://localhost:$${SERVER_PORT:-8000}/swagger/index.html"

scaffold: ## Scaffold a new domain (usage: make scaffold name=book)
	@if [ -z "$(name)" ]; then \
		echo "Error: usage: make scaffold name=book"; exit 1; \
	fi
	@DOMAIN_PKG=$$(echo "$(name)" | tr 'A-Z' 'a-z'); \
	DOMAIN_PASCAL=$$(echo "$(name)" | awk -F'[-_]' '{for(i=1;i<=NF;i++){$$i=toupper(substr($$i,1,1)) substr($$i,2)}}1' OFS=''); \
	DOMAIN_DIR=internal/domain/$$DOMAIN_PKG; \
	REPO_DIR=$$DOMAIN_DIR/repository; \
	SVC_DIR=$$DOMAIN_DIR/service; \
	MODULE_DIR=internal/infra/bootstrap/modules; \
	mkdir -p $$REPO_DIR $$SVC_DIR $$MODULE_DIR; \
	sed -e "s/{{name}}/$$DOMAIN_PKG/g" -e "s/{{Name}}/$$DOMAIN_PASCAL/g" templates/repository/base_repository.go.stub > $$REPO_DIR/$$DOMAIN_PKG\_repository.go; \
	sed -e "s/{{name}}/$$DOMAIN_PKG/g" -e "s/{{Name}}/$$DOMAIN_PASCAL/g" templates/repository/base_repository_interface.go.stub > $$REPO_DIR/$$DOMAIN_PKG\_repository_interface.go; \
	sed -e "s/{{name}}/$$DOMAIN_PKG/g" -e "s/{{Name}}/$$DOMAIN_PASCAL/g" templates/service/base_service.go.stub > $$SVC_DIR/$$DOMAIN_PKG\_service.go; \
	sed -e "s/{{name}}/$$DOMAIN_PKG/g" -e "s/{{Name}}/$$DOMAIN_PASCAL/g" templates/service/base_service_interface.go.stub > $$SVC_DIR/$$DOMAIN_PKG\_service_interface.go; \
	sed -e "s/{{name}}/$$DOMAIN_PKG/g" -e "s/{{Name}}/$$DOMAIN_PASCAL/g" templates/bootstrap/module.go.stub > $$MODULE_DIR/$$DOMAIN_PKG\_module.go; \
	gofmt -w $$REPO_DIR $$SVC_DIR $$MODULE_DIR; \
	echo "Scaffolded $$DOMAIN_PASCAL in internal/domain/$$DOMAIN_PKG"
