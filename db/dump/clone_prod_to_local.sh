#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROD_DATABASE_URL="postgresql://postgres:QuGQnAvhEZOEfSfMdjVzbyhTtYpVUxBx@yamanote.proxy.rlwy.net:48247/railway?sslmode=require"
PG_DOCKER_IMAGE="${PG_DOCKER_IMAGE:-postgres:17}"
PROD_DUMP_FILE="${PROD_DUMP_FILE:-$SCRIPT_DIR/prod.dump}"
PROD_DUMP_BASENAME="$(basename "${PROD_DUMP_FILE}")"
LOCAL_DB_NAME="${LOCAL_DB_NAME:-dragonki}"
LOCAL_DB_ADMIN_URL="${LOCAL_DB_ADMIN_URL:-postgres://dragonki:dragonki@host.docker.internal:5432/postgres?sslmode=disable}"
LOCAL_DB_URL="${LOCAL_DB_URL:-postgres://dragonki:dragonki@host.docker.internal:5432/${LOCAL_DB_NAME}?sslmode=disable}"

if ! command -v docker >/dev/null 2>&1; then
  echo "docker nao encontrado no PATH." >&2
  exit 1
fi

if [[ -z "${PROD_DATABASE_URL:-}" ]]; then
  echo "PROD_DATABASE_URL nao definido." >&2
  echo "Exemplo:" >&2
  echo "  export PROD_DATABASE_URL='postgres://user:pass@host:5432/db?sslmode=require'" >&2
  exit 1
fi

echo "== Baixando dump de producao =="
docker run --rm \
  -e "PROD_DATABASE_URL=${PROD_DATABASE_URL}" \
  -v "${SCRIPT_DIR}:/work" \
  "${PG_DOCKER_IMAGE}" \
  sh -lc "pg_dump --dbname=\"\$PROD_DATABASE_URL\" --format=custom --no-owner --no-privileges --file=/work/${PROD_DUMP_BASENAME}"

echo "== Dropando base local ${LOCAL_DB_NAME} =="
docker run --rm \
  "${PG_DOCKER_IMAGE}" \
  psql "${LOCAL_DB_ADMIN_URL}" \
  -v ON_ERROR_STOP=1 \
  -c "DROP DATABASE IF EXISTS ${LOCAL_DB_NAME};"

echo "== Criando base local ${LOCAL_DB_NAME} =="
docker run --rm \
  "${PG_DOCKER_IMAGE}" \
  psql "${LOCAL_DB_ADMIN_URL}" \
  -v ON_ERROR_STOP=1 \
  -c "CREATE DATABASE ${LOCAL_DB_NAME};"

echo "== Restaurando dump na base local =="
docker run --rm \
  -v "${SCRIPT_DIR}:/work" \
  "${PG_DOCKER_IMAGE}" \
  pg_restore \
  --dbname="${LOCAL_DB_URL}" \
  --clean \
  --if-exists \
  --no-owner \
  --no-privileges \
  "/work/${PROD_DUMP_BASENAME}"

echo "== Restore concluido =="
echo "Dump salvo em: ${PROD_DUMP_FILE}"
echo "Base local restaurada em: ${LOCAL_DB_URL}"
