.PHONY: help build run test test-unit test-coverage mocks migrate migrate-up migrate-down migrate-force clean docker-up docker-down

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the REST API binary
	go build -o bin/api cmd/rest/main.go

build-docker: ## Build Docker image for the application
	docker build -t task-management-api .

# Run targets
run: ## Run the REST API server locally
	go run cmd/rest/main.go

run-dev: ## Run with air for hot reloading (requires air: go install github.com/air-verse/air@latest)
	air -c .air.toml

# Test targets
test: mocks ## Run all tests
	go test ./... -v

test-unit: mocks ## Run only unit tests
	go test ./... -v -short

test-coverage: mocks ## Generate test coverage report
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: mocks ## Run tests with race detection
	go test ./... -race

# Mock generation
mocks: ## Generate mocks using mockery
	@if command -v mockery >/dev/null 2>&1; then \
		mockery; \
	else \
		echo "Installing mockery..."; \
		go install github.com/vektra/mockery/v2@latest; \
		mockery; \
	fi

# Database migration targets
migrate-install: ## Install golang-migrate tool
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up: ## Run database migrations up
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down: ## Run database migrations down
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-force: ## Force migration version (use with VERSION=N)
	migrate -path migrations -database "$(DATABASE_URL)" force $(VERSION)

migrate-version: ## Check current migration version
	migrate -path migrations -database "$(DATABASE_URL)" version

migrate-create: ## Create a new migration file (use with NAME=migration_name)
	migrate create -ext sql -dir migrations -seq $(NAME)

# Docker targets
docker-up: ## Start development environment with docker-compose
	cd .devcontainer && docker-compose up -d

docker-down: ## Stop development environment
	cd .devcontainer && docker-compose down

docker-logs: ## Show logs from all containers
	cd .devcontainer && docker-compose logs -f

docker-db: ## Connect to PostgreSQL database
	cd .devcontainer && docker-compose exec db psql -U user -d taskdb

docker-redis: ## Connect to Redis CLI
	cd .devcontainer && docker-compose exec redis redis-cli

# Development setup
setup: ## Set up development environment
	go mod tidy
	cp .env.example .env
	@echo "Development environment setup complete!"
	@echo "1. Update .env file with your configuration"
	@echo "2. Run 'make docker-up' to start services"
	@echo "3. Run 'make migrate-up' to set up database"
	@echo "4. Run 'make run' to start the API server"

setup-devcontainer: ## Set up devcontainer environment
	go mod tidy
	cp .env.example .devcontainer/.env
	@echo "Devcontainer setup complete!"

# Code quality
lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

# Clean targets
clean: ## Clean build artifacts and generated files
	rm -rf bin/
	rm -f coverage.out coverage.html
	find . -name "mocks" -type d -exec rm -rf {} +
	go clean -cache
	go clean -testcache

clean-docker: ## Clean Docker containers and volumes
	cd .devcontainer && docker-compose down -v
	docker system prune -f

# Environment variables for common commands
DATABASE_URL ?= postgres://user:password@localhost:5432/taskdb?sslmode=disable
TEST_DATABASE_URL ?= postgres://user:password@localhost:5432/taskdb_test?sslmode=disable

# Load environment variables from .env if it exists
ifneq (,$(wildcard ./.env))
    include ./.env
    export
endif