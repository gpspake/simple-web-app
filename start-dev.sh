#!/bin/bash

# Run postgres in a container and run the go app locally using air for auto reloading

# PostgreSQL container details
POSTGRES_CONTAINER_NAME="go_app_postgres"
POSTGRES_IMAGE="postgres:15"
POSTGRES_PORT="5432"
POSTGRES_USER="your_user"
POSTGRES_PASSWORD="your_password"
POSTGRES_DB="your_db"

# Start PostgreSQL container
echo "Starting PostgreSQL container..."
docker run --name "$POSTGRES_CONTAINER_NAME" \
  -e POSTGRES_USER="$POSTGRES_USER" \
  -e POSTGRES_PASSWORD="$POSTGRES_PASSWORD" \
  -e POSTGRES_DB="$POSTGRES_DB" \
  -p "$POSTGRES_PORT:5432" \
  -d "$POSTGRES_IMAGE"

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until docker exec "$POSTGRES_CONTAINER_NAME" pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" >/dev/null 2>&1; do
    printf '.'
    sleep 1
done

echo "PostgreSQL is ready!"

# Set environment variables for the Go application
export POSTGRES_USER="$POSTGRES_USER"
export POSTGRES_PASSWORD="$POSTGRES_PASSWORD"
export POSTGRES_DB="$POSTGRES_DB"
export POSTGRES_HOST="localhost"
export POSTGRES_PORT="$POSTGRES_PORT"
export TEMPLATE_DIR="$(pwd)/internal/templates" # Set the templates path for local runs
export MIGRATIONS_PATH="file://$(pwd)/migrations" # Set the migrations path for local runs



# Run the Go application
echo "Starting the Go application..."
air -c .air.toml & npm run tailwind:watch

# Cleanup: Stop and remove the PostgreSQL container when the script exits
cleanup() {
    echo "Stopping and removing PostgreSQL container..."
    docker stop "$POSTGRES_CONTAINER_NAME" && docker rm "$POSTGRES_CONTAINER_NAME"
}
trap cleanup EXIT
