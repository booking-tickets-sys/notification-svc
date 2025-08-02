.PHONY: build build-server build-worker run-server run-worker test clean docker-build docker-run

# Build variables
BINARY_SERVER=notification-server
BINARY_WORKER=notification-worker
BUILD_DIR=build

# Build the server binary
build-server:
	@echo "Building notification server..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_SERVER) ./cmd/server

# Build the worker binary
build-worker:
	@echo "Building notification worker..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_WORKER) ./cmd/worker

# Build both binaries
build: build-server build-worker

# Run the server
run-server: build-server
	@echo "Starting notification server..."
	./$(BUILD_DIR)/$(BINARY_SERVER)

# Run the worker
run-worker: build-worker
	@echo "Starting notification worker..."
	./$(BUILD_DIR)/$(BINARY_WORKER)

# Run both server and worker (requires tmux)
run: build
	@echo "Starting notification service (server + worker)..."
	@if command -v tmux >/dev/null 2>&1; then \
		tmux new-session -d -s notification-svc './$(BUILD_DIR)/$(BINARY_SERVER)' \; \
		split-window -h './$(BUILD_DIR)/$(BINARY_WORKER)' \; \
		attach-session -d notification-svc; \
	else \
		echo "tmux not found. Please install tmux or run server and worker separately."; \
		echo "Run 'make run-server' in one terminal and 'make run-worker' in another."; \
	fi

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t notification-svc .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file config.env notification-svc

# Install Asynq CLI tool
install-asynq-cli:
	@echo "Installing Asynq CLI tool..."
	go install github.com/hibiken/asynq/tools/asynq@latest

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build both server and worker binaries"
	@echo "  build-server   - Build server binary only"
	@echo "  build-worker   - Build worker binary only"
	@echo "  run-server     - Run the notification server"
	@echo "  run-worker     - Run the notification worker"
	@echo "  run            - Run both server and worker (requires tmux)"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  install-asynq-cli - Install Asynq CLI tool"
	@echo "  help           - Show this help message" 