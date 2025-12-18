#!/usr/bin/env bash
# Start a local Postgres Docker container and seed schema.sql
set -euo pipefail
. "$(dirname "$0")/_helpers.sh" || true

# Configurable via env or CLI
DB_CONTAINER="${DB_CONTAINER:-books-db}"
DB_USER="${DB_USER:-user}"
DB_PASS="${DB_PASS:-password}"
DB_NAME="${DB_NAME:-booksdb}"
SCHEMA_FILE="${SCHEMA_FILE:-schema.sql}"
IMAGE="${IMAGE:-postgres:16-alpine}"
HOST_PORT="${HOST_PORT:-5432}"

info "Starting Postgres container '${DB_CONTAINER}' (image ${IMAGE})..."

# If container exists but is stopped, start it; if running, skip
if docker ps -a --format '{{.Names}}' | grep -q "^${DB_CONTAINER}$"; then
  status=$(docker inspect -f '{{.State.Status}}' "${DB_CONTAINER}")
  if [ "$status" = "running" ]; then
    info "Container ${DB_CONTAINER} already running."
  else
    info "Starting existing container ${DB_CONTAINER}..."
    docker start "${DB_CONTAINER}"
  fi
else
  docker run --name "${DB_CONTAINER}" \
    -e POSTGRES_USER="${DB_USER}" \
    -e POSTGRES_PASSWORD="${DB_PASS}" \
    -e POSTGRES_DB="${DB_NAME}" \
    -p "${HOST_PORT}:5432" -d "${IMAGE}"
fi

# Wait for readiness then seed
wait_for_pg_container "${DB_CONTAINER}"

if [ -f "${SCHEMA_FILE}" ]; then
  info "Seeding database from ${SCHEMA_FILE}..."
  docker exec -i "${DB_CONTAINER}" psql -U "${DB_USER}" -d "${DB_NAME}" < "${SCHEMA_FILE}"
  info "Seed completed."
else
  warn "Schema file '${SCHEMA_FILE}' not found — skipping seed."
fi

info "Postgres is up: container='${DB_CONTAINER}', user='${DB_USER}', db='${DB_NAME}', host_port='${HOST_PORT}'"