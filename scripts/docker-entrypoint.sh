#!/bin/sh
set -eu

DB_HOST="${DB_HOST:-postgres}"
DB_PORT="${DB_PORT:-5432}"

echo "Waiting for database at ${DB_HOST}:${DB_PORT}..."
while ! nc -z "${DB_HOST}" "${DB_PORT}" >/dev/null 2>&1; do
  sleep 1
done

echo "Running migrations..."
if ! migrate_output="$(migrate -path /app/db/migration -database "${DB_SOURCE}" up 2>&1)"; then
  case "${migrate_output}" in
    *"no change"*)
      echo "Migrations already up to date"
      ;;
    *)
      echo "${migrate_output}"
      echo "Migration failed"
      exit 1
      ;;
  esac
fi

echo "Starting API server..."
exec /app/simplebank
