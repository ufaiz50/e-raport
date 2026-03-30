#!/usr/bin/env bash
set -euo pipefail

export PGPASSWORD="${POSTGRES_PASSWORD:-password}"
psql \
  -h "${POSTGRES_HOST:-127.0.0.1}" \
  -p "${POSTGRES_PORT:-5435}" \
  -U "${POSTGRES_USER:-docker}" \
  -d "${POSTGRES_DB:-go_app_dev}" \
  -f "$(dirname "$0")/reset_uuid_schema.sql"

echo "Schema reset complete. Start the backend to AutoMigrate UUID-native tables."
