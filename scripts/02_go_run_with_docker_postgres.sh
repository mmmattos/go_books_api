#!/usr/bin/env bash
# Start docker Postgres if needed, then run 'go run' on host, pointing to that Postgres.
set -euo pipefail
. "$(dirname "$0")/_helpers.sh" || true

# Configurable
DB_CONTAINER="${DB_CONTAINER:-books-db}"
DB_USER="${DB_USER:-user}"
DB_PASS="${DB_PASS:-password}"
DB_NAME="${DB_NAME:-booksdb}"
HOST_PORT="${HOST_PORT:-5432}"
SCHEMA_FILE="${SCHEMA_FILE:-schema.sql}"
PORT="${PORT:-8080}"

# Ensure Postgres running (reuse script 01)
if ! docker ps --format '{{.Names}}' | grep -q "^${DB_CONTAINER}$"; then
  info "No running Postgres container named ${DB_CONTAINER}. Starting one..."
  bash "$(dirname "$0")/01_start_postgres_and_seed.sh"
else
  info "Found running Postgres container ${DB_CONTAINER}."
fi

# Run the app on host using go run
export DB_CONN="postgres://${DB_USER}:${DB_PASS}@localhost:${HOST_PORT}/${DB_NAME}?sslmode=disable"
info "Running 'go run' with DB_CONN=${DB_CONN} PORT=${PORT}"
PORT="${PORT}" DB_CONN="${DB_CONN}" go run ./cmd/api-service