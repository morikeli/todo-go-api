#!/usr/bin/env bash

set -euo pipefail

project_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$project_root"

if [[ -f .env ]]; then
  set -a
  # shellcheck disable=SC1091
  source .env
  set +a
fi

: "${DB_USER:?DB_USER must be set in .env or the environment}"
: "${DB_PASSWORD:?DB_PASSWORD must be set in .env or the environment}"
: "${DB_NAME:?DB_NAME must be set in .env or the environment}"

if ! command -v migrate >/dev/null 2>&1; then
  printf '%s\n' 'migrate is not installed or is not on PATH.' >&2
  printf '%s\n' 'Install it with: go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest' >&2
  exit 1
fi

docker compose up -d db

until docker compose exec -T db pg_isready -U "$DB_USER" -d "$DB_NAME" >/dev/null 2>&1; do
  sleep 1
done

database_url="postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT:-5432}/${DB_NAME}?sslmode=disable"

if [[ $# -eq 0 ]]; then
  set -- up
fi

migrate -path "$project_root/migrations" -database "$database_url" "$@"