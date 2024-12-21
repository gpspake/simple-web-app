# Variables
DOCKER_COMPOSE = docker-compose
APP_NAME = go_app
POSTGRES_CONTAINER = go_app_postgres

# Default target
.DEFAULT_GOAL := help

# Help target (lists all commands)
help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Targets
run: ## Run the application in production mode
	@echo "Starting the application in production mode..."
	@bash start.sh

dev: ## Run the application in development mode with live reloading
	@echo "Starting the application in development mode..."
	@bash start-dev.sh

test: ## Run tests with Docker Compose
	@echo "Running tests..."
	@bash test.sh

docker-up: ## Start services using Docker Compose
	@echo "Starting services with Docker Compose..."
	@$(DOCKER_COMPOSE) up -d

docker-down: ## Stop and remove Docker Compose services
	@echo "Stopping and removing services..."
	@$(DOCKER_COMPOSE) down --volumes

docker-restart: docker-down docker-up ## Restart services with Docker Compose

clean: ## Clean up Docker containers, images, and volumes
	@echo "Cleaning up unused Docker resources..."
	@docker stop $(POSTGRES_CONTAINER) 2>/dev/null || true
	@docker rm $(POSTGRES_CONTAINER) 2>/dev/null || true
	@docker image prune -f
	@docker volume prune -f

logs: ## Show logs for the Docker Compose app service
	@echo "Fetching logs for app service..."
	@$(DOCKER_COMPOSE) logs -f app

build: ## Build the production Docker image
	@echo "Building the production Docker image..."
	@docker build -t $(APP_NAME):latest .

migrate-up: ## Run database migrations (up)
	@echo "Running database migrations..."
	@migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" up

migrate-down: ## Roll back database migrations (down)
	@echo "Rolling back database migrations..."
	@migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" down

seed-db: ## Run the database seed command
	@echo "Seeding the database..."
	@go run cmd/seed_db/seed_db.go

populate-fts: ## Populate the Full-Text Search table
	@echo "Populating the Full-Text Search table..."
	@go run cmd/populate_fts/populate_fts.go

coverage: ## Generate and display test coverage report
	@echo "Running tests and generating coverage report..."
	@docker-compose -f docker-compose.test.yml run test-runner

format: ## Format Go code
	@echo "Formatting Go code..."
	@gofmt -s -w .

lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run ./...

check: format lint test ## Format, lint, and test the code

.PHONY: help run dev test docker-up docker-down docker-restart clean logs build migrate-up migrate-down coverage format lint check
