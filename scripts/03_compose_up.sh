#!/usr/bin/env bash
# Bring up docker-compose, build images, and tail logs
set -euo pipefail
. "$(dirname "$0")/_helpers.sh" || true

COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.yml}"

if ! command -v docker-compose >/dev/null 2>&1; then
  warn "docker-compose not found. Use 'docker compose' (integrated) as fallback."
  docker compose up --build
  exit 0
fi

info "Starting docker-compose using ${COMPOSE_FILE}..."
docker-compose -f "${COMPOSE_FILE}" up --build -d
info "Services started. Tailing logs (ctrl-c to stop)."
docker-compose -f "${COMPOSE_FILE}" logs -f