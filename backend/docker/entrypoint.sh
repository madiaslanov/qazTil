#!/bin/sh
set -eu

data_dir=$(dirname "$DB_PATH")
mkdir -p "$data_dir"
chown -R app:app "$data_dir"
exec su-exec app /app/qaztil
