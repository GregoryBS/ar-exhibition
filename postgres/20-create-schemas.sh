#!/usr/bin/env bash
set -e

echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/13/main/pg_hba.conf
echo "listen_addresses='*'" >> /etc/postgresql/13/main/postgresql.conf
export SCRIPT_PATH=/docker-entrypoint-initdb.d/
export PGPASSWORD=test
psql -U program -d art -f "$SCRIPT_PATH/schemes/schema.sql"
