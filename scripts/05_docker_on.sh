#!/usr/bin/env bash
# Run both Postgres and the Go app as containers on a custom docker network (no host Go)
set -euo pipefail
. "$(dirname "$0")/_helpers.sh" || true

NETWORK="${NETWORK:-books-net}"
DB_CONTAINER="${DB_CONTAINER:-books-db}"
DB_USER="${DB_USER:-user}"
DB_PASS="${DB_PASS:-password}"
DB_NAME="${DB_NAME:-booksdb}"
APP_IMAGE="${APP_IMAGE:-books-api}"
PORT="${PORT:-8080}"

# 1) Create network if missing
if ! docker network ls --format '{{.Name}}' | grep -q "^${NETWORK}$"; then
  info "Creating docker network ${NETWORK}..."
  docker network create "${NETWORK}"
fi

# 2) Start Postgres container attached to network (if not exists)
if ! docker ps -a --format '{{.Names}}' | grep -q "^${DB_CONTAINER}$"; then
  info "Running Postgres container ${DB_CONTAINER}..."
  docker run --name "${DB_CONTAINER}" --network "${NETWORK}" \
    -e POSTGRES_USER="${DB_USER}" -e POSTGRES_PASSWORD="${DB_PASS}" -e POSTGRES_DB="${DB_NAME}" \
    -d postgres:16-alpine
  # wait then seed if schema exists
  ./scripts/01_start_postgres_and_seed.sh
else
  status=$(docker inspect -f '{{.State.Status}}' "${DB_CONTAINER}")
  if [ "$status" != "running" ]; then
    info "Starting existing container ${DB_CONTAINER}..."
    docker start "${DB_CONTAINER}"
  else
    info "Postgres container ${DB_CONTAINER} already running."
  fi
fi

# 3) Build app image (multi-stage Dockerfile expected)
info "Building app image ${APP_IMAGE}..."
docker build -t "${APP_IMAGE}" .

# 4) Run app container
if docker ps --format '{{.Names}}' | grep -q "^books-api$"; then
  info "Stopping existing books-api container..."
  docker rm -f books-api
fi

info "Running app container and connecting to postgres via hostname '${DB_CONTAINER}'..."
docker run --name books-api --network "${NETWORK}" \
  -e PORT="${PORT}" \
  -e DB_CONN="postgres://${DB_USER}:${DB_PASS}@${DB_CONTAINER}:5432/${DB_NAME}?sslmode=disable" \
  -p "${PORT}:8080" -d "${APP_IMAGE}"

info "App is accessible at http://localhost:${PORT}"