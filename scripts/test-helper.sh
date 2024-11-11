#!/usr/bin/env bash
DB_DATABASE=$(grep DB_DATABASE test.env | cut -d '=' -f2)
psql -d "${DB_DATABASE}" -c "DROP SCHEMA public CASCADE;"
psql -d "${DB_DATABASE}" -c "CREATE SCHEMA public;"
