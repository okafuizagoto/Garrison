#!/bin/bash
set -e

# Copy .env if missing
if [ ! -f "/app/${APP_NAME}/.env" ]; then
    echo "[build.sh] Copying .env.example -> .env for ${APP_NAME}"
    cp /app/${APP_NAME}/.env.example /app/${APP_NAME}/.env
fi

# Initialize go.mod if missing
if [ ! -f "/app/${APP_NAME}/go.mod" ]; then
    echo "[build.sh] go mod init for ${APP_NAME}"
    go mod init ${APP_NAME}
fi

echo "[build.sh] Starting ${APP_NAME}..."
go run /app/${APP_NAME}/main.go
