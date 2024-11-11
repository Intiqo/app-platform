#!/usr/bin/env bash
if [ ! -f "test.env" ]
then
  touch test.env
fi
set -o allexport
. test.env set
echo "Running migrations for tests"
POSTGRES_URL="host=${DB_HOST} port=${DB_PORT} user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_DATABASE_NAME} sslmode=disable"
goose --dir './internal/database/migrations' postgres "${POSTGRES_URL}" up
