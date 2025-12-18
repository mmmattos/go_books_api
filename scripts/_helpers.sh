#!/usr/bin/env bash
# small helper functions used across scripts
set -euo pipefail

info(){ echo -e "\033[1;34m[INFO]\033[0m $*"; }
warn(){ echo -e "\033[1;33m[WARN]\033[0m $*"; }
err(){ echo -e "\033[1;31m[ERROR]\033[0m $*" >&2; }

# Wait for Postgres to be ready in a container
wait_for_pg_container() {
  local container=$1
  info "Waiting for Postgres container '$container' to be ready..."
  local retries=30
  local i=0
  until docker exec -i "$container" pg_isready -U "${DB_USER:-user}" -d "${DB_NAME:-booksdb}" >/dev/null 2>&1; do
    i=$((i+1))
    if [ $i -gt $retries ]; then
      err "Postgres did not become ready in time."
      return 1
    fi
    sleep 1
  done
  info "Postgres is ready."
}