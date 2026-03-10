#!/usr/bin/env bash
set -euo pipefail

POSTGRES_CONTAINER="postgres17"

if command -v docker >/dev/null 2>&1; then
  docker stop "${POSTGRES_CONTAINER}" >/dev/null 2>&1 || true
fi

echo "PostgreSQL container stopped (if it was running)."
