#!/usr/bin/env bash
set -e

export SCRIPT_PATH=/docker-entrypoint-initdb.d/
export PGPASSWORD=test
#psql -U program -d art -f "$SCRIPT_PATH/schemes/schema.sql"
