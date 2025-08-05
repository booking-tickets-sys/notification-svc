# Notification Service

A high-performance notification service built with Go that processes order events using Redis and Asynq for reliable async job processing. The service is designed to handle order creation events and send appropriate notifications to users.

## Features

- **Event-Driven Architecture**: Processes order creation events from Redis streams
- **Async Job Processing**: Uses Asynq for reliable Redis-based job queue processing
- **Worker Pool**: Configurable number of concurrent workers for processing notifications
- **Order Event Handling**: Specialized handling for order creation events
- **Graceful Shutdown**: Proper cleanup and shutdown handling with context cancellation
- **Configuration Management**: Uses Viper for flexible configuration (files, environment variables)
- **Structured Logging**: Uses Logrus for structured JSON logging with contextual fields
- **Health Monitoring**: Built-in health monitoring and worker status tracking

## Prerequisites

- Go 1.24.4 or later
- Redis server running locally or remotely

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd notification-svc
```

2. Install dependencies:
```bash
go mod tidy
```

3. Start Redis server (if not already running):
```bash
# Using Docker
docker run -d -p 6379:6379 redis:alpine

# Or using Homebrew (macOS)
brew install redis
brew services start redis
```

## Configuration

The service uses Viper for configuration management, supporting multiple configuration sources:

### Configuration File

The service looks for configuration in the `configs` directory. Create a `configs/config.yaml` file:

```yaml
# Notification Service Configuration

server:
  port: "8080"
  host: "0.0.0.0"
  read_timeout: 30
  write_timeout: 30
  idle_timeout: 60

redis:
  addr: "localhost:6379"
  password: ""
  db: 0
  pool_size: 10

queue:
  name: "notifications"

logging:
  level: "info"  # debug, info, warn, error, fatal, panic
  format: "json" # json, text
  output: "stdout"

workers:
  count: 5
```

### Environment Variables

You can override any configuration using environment variables:

```bash
# Redis configuration
export REDIS_ADDR=redis.example.com:6379
export REDIS_PASSWORD=your_password
export REDIS_DB=1
export REDIS_POOL_SIZE=20

# Logging configuration
export LOGGING_LEVEL=debug
export LOGGING_FORMAT=text

# Worker configuration
export WORKERS_COUNT=10
```

### Configuration Priority

1. Environment variables (highest priority)
2. Configuration file (`configs/config.yaml`)
3. Default values (lowest priority)

## Running the Service

### Using Make

```bash
# Build and run
make run

# Development mode (starts Redis automatically)
make dev

# Stop development environment
make dev-stop

# See all available commands
make help
```

### Direct Execution

```bash
# Run with default configuration
go run cmd/notification-svc/main.go

# Run with custom config directory
CONFIG_DIR=/path/to/configs go run cmd/notification-svc/main.go
```

## Architecture

### Event Processing Flow

1. **Event Reception**: The service listens for order creation events from Redis streams
2. **Job Enqueuing**: Events are converted to Asynq tasks and enqueued in Redis
3. **Worker Processing**: Multiple workers process tasks concurrently
4. **Notification Sending**: Based on the event type, appropriate notifications are sent
5. **Error Handling**: Failed tasks are retried with exponential backoff

### Supported Events

#### Order Created Event

**Event Name**: `order_created`

**Payload Structure**:
```json
{
  "eventMetadata": {
    "eventID": "evt_123456789",
    "eventName": "order_created",
    "publishedAt": 1642234567890
  },
  "orderId": "ord_123456789",
  "userId": "user_123456789",
  "ticketId": "tkt_123456789"
}
```

**Processing**: When an order is created, the service sends confirmation notifications to the user.

## Project Structure

```
notification-svc/
├── cmd/
│   └── notification-svc/
│       └── main.go              # Main application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management with Viper
│   ├── logger/
│   │   └── logger.go            # Structured logging with Logrus
│   ├── handler/
│   │   └── order.go             # Order event handlers
│   ├── worker/
│   │   ├── notification.go      # Notification worker implementation
│   │   └── notification_test.go # Worker tests
│   └── models/
│       └── events/
│           ├── events.go        # Event metadata structures
│           └── order.go         # Order-related event payloads
├── pkg/
│   └── events/                  # Public event packages
├── configs/
│   ├── config.yaml              # Sample configuration file
│   └── env.example              # Environment variables example
├── scripts/                     # Build and utility scripts
├── docs/                        # Documentation
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── Makefile                     # Build and development commands
└── README.md                    # Documentation
```

## Development

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v -run TestNotificationWorker
```

### Code Quality

```bash
# Format code
make fmt

# Lint code
make lint

# Vet code
make vet
```

### Building and Development

```bash
# Build binary
make build

# Install dependencies
make deps

# Install development tools
make install-tools
```

### Code Quality

```bash
# Format code
make fmt

# Lint code
make lint

# Vet code
make vet
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

## Available Make Commands

The project includes a comprehensive Makefile with the following commands:

### Core Commands
- `make build` - Build the application
- `make run` - Build and run the application
- `make clean` - Clean build artifacts

### Development Commands
- `make dev` - Start development environment (Redis + service)
- `make dev-stop` - Stop development environment
- `make deps` - Download and tidy dependencies

### Testing Commands
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report

### Code Quality Commands
- `make fmt` - Format code
- `make lint` - Run linter
- `make vet` - Vet code
- `make install-tools` - Install development tools

### Help
- `make help` - Show all available commands

## Monitoring and Observability

### Logging

The service uses Logrus for structured logging with the following features:

#### Log Levels

- `debug`: Detailed debug information
- `info`: General information about application flow
- `warn`: Warning messages for potentially harmful situations
- `error`: Error messages for error conditions
- `fatal`: Fatal errors that cause the application to exit
- `panic`: Panic messages

#### Log Formats

- `json`: Structured JSON logging (default)
- `text`: Human-readable text format

#### Structured Fields

The logger automatically adds contextual fields:

```json
{
  "level": "info",
  "msg": "Processing order created event",
  "worker_id": 1,
  "order_id": "ord_123456789",
  "user_id": "user_123456789",
  "time": "2024-01-15T10:30:00Z"
}
```

### Worker Status

The service provides methods to check worker status:

```go
// Check if worker is running
isRunning := notificationWorker.IsRunning()
```

## Production Considerations

1. **Redis Configuration**: Use Redis Cluster or Redis Sentinel for high availability
2. **Monitoring**: Add metrics collection (Prometheus, etc.)
3. **Logging**: Configure log aggregation (ELK stack, etc.)
4. **Error Handling**: Add retry mechanisms and dead letter queues
5. **Scaling**: Run multiple instances behind a load balancer
6. **Configuration**: Use environment variables for sensitive configuration
7. **Health Checks**: Implement comprehensive health checks
8. **Security**: Follow security best practices for production deployments
9. **Graceful Shutdown**: Ensure proper cleanup of resources
10. **Circuit Breakers**: Implement circuit breakers for external service calls

## Troubleshooting

### Common Issues

1. **Redis Connection Error**: Ensure Redis is running and accessible
2. **Worker Not Processing**: Check worker logs and Redis connection
3. **Configuration Not Loading**: Verify config directory path and format
4. **Graceful Shutdown Issues**: Check for proper context cancellation

### Logs

The service provides detailed structured logging for debugging:

```bash
# View logs in JSON format
go run cmd/notification-svc/main.go | jq '.'

# View logs in text format
LOGGING_FORMAT=text go run cmd/notification-svc/main.go
```

### Health Checks

```bash
# Check Redis connection
redis-cli ping

# Check worker status (programmatically)
# Use the IsRunning() method on the worker instance
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

[Add your license here] 