#!/bin/sh

set -eu

run_migrations() {
  if [ "${RUN_MIGRATIONS:-false}" != "true" ]; then
    echo "Skipping migrations because RUN_MIGRATIONS=${RUN_MIGRATIONS:-false}"
    return
  fi

  echo "Applying goose migrations to ${DB_NAME}..."
  if [ -x /app/migrate ]; then
    /app/migrate up
  else
    go run ./cmd/migrate up
  fi
}

run_migrations

exec "$@"
