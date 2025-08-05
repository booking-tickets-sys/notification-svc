# Makefile for notification-svc

# Variables
BINARY_NAME=notification-svc
BUILD_DIR=bin
MAIN_PATH=cmd/notification-svc
CONFIG_DIR=configs

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linting
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Vet code
.PHONY: vet
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Development mode (starts Redis and runs the service)
.PHONY: dev
dev:
	@echo "Starting development environment..."
	@echo "Starting Redis..."
	@docker run -d --name redis-dev -p 6379:6379 redis:alpine || echo "Redis container already running"
	@echo "Starting notification service..."
	@$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Stop development environment
.PHONY: dev-stop
dev-stop:
	@echo "Stopping development environment..."
	@docker stop redis-dev || echo "Redis container not running"
	@docker rm redis-dev || echo "Redis container not found"

# Install development tools
.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run security scan
.PHONY: security
security:
	@echo "Running security scan..."
	gosec ./...

# Run all checks (lint, test, security)
.PHONY: check
check: lint test security

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t notification-svc:latest .

# Run with Docker Compose
.PHONY: docker-up
docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

# Stop Docker Compose services
.PHONY: docker-down
docker-down:
	@echo "Stopping Docker Compose services..."
	docker-compose down

# Show help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Build and run the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  security       - Run security scan"
	@echo "  check          - Run all checks (lint, test, security)"
	@echo "  deps           - Download dependencies"
	@echo "  dev            - Start development environment"
	@echo "  dev-stop       - Stop development environment"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-up      - Start services with Docker Compose"
	@echo "  docker-down    - Stop Docker Compose services"
	@echo "  install-tools  - Install development tools"
	@echo "  help           - Show this help message"