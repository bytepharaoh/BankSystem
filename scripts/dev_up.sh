#!/usr/bin/env bash
set -euo pipefail

# Boots local dependencies and prepares app for development.
# Usage:
#   ./scripts/dev_up.sh            # setup only
#   ./scripts/dev_up.sh --run      # setup + run server

RUN_SERVER=false
if [[ "${1:-}" == "--run" ]]; then
  RUN_SERVER=true
fi

POSTGRES_CONTAINER="postgres17"
DB_NAME="simple_bank"
DB_USER="root"

if ! command -v docker >/dev/null 2>&1; then
  echo "[error] docker is not installed"
  exit 1
fi

echo "==> Starting PostgreSQL container if needed"
if ! docker ps -a --format '{{.Names}}' | grep -q "^${POSTGRES_CONTAINER}$"; then
  make postgres
else
  docker start "${POSTGRES_CONTAINER}" >/dev/null 2>&1 || true
fi

echo "==> Waiting for PostgreSQL readiness"
for _ in {1..30}; do
  if docker exec "${POSTGRES_CONTAINER}" pg_isready -U "${DB_USER}" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

echo "==> Ensuring database exists"
if ! docker exec "${POSTGRES_CONTAINER}" psql -U "${DB_USER}" -tAc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'" | grep -q 1; then
  docker exec "${POSTGRES_CONTAINER}" createdb --username="${DB_USER}" --owner="${DB_USER}" "${DB_NAME}"
fi

echo "==> Running migrations"
make migrateup

echo "==> Regenerating SQLC and mocks"
make sqlc
GOCACHE=$(pwd)/.gocache make mock

if [[ "$RUN_SERVER" == "true" ]]; then
  echo "==> Starting API server"
  make server
else
  echo "Done. Environment is ready. Run 'make server' to start the API."
fi
