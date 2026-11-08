#!/usr/bin/env sh
set -eu

cd "$(dirname "$0")/../backend/api"
STORAGE_DRIVER="${STORAGE_DRIVER:-memory}" go run ./cmd/api
