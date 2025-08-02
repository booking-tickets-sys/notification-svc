# Notification Service

A high-performance, asynchronous notification service built with Go, Redis, and Asynq. This service provides a scalable solution for sending various types of notifications (email, SMS, push, webhook) through a distributed task queue system.

## Features

- **Asynchronous Processing**: Uses Redis and Asynq for reliable, distributed task processing
- **Multiple Notification Types**: Support for email, SMS, push notifications, and webhooks
- **Priority Queues**: Different priority levels (critical, default, low) for task processing
- **RESTful API**: Clean HTTP API for enqueueing notifications
- **Scalable Architecture**: Separate server and worker processes for horizontal scaling
- **Docker Support**: Complete containerization with docker-compose
- **Monitoring**: Optional Asynqmon integration for queue monitoring

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Server   │───▶│   Redis Queue   │───▶│   Worker(s)     │
│   (API)         │    │   (Asynq)       │    │   (Processors)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Prerequisites

- Go 1.24 or higher
- Redis 4.0 or higher
- Docker and Docker Compose (optional)

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd notification-svc
   ```

2. **Start the services**
   ```bash
   docker-compose up -d
   ```

3. **Test the service**
   ```bash
   # Send an email notification
   curl -X POST http://localhost:8080/api/v1/notifications/email \
     -H "Content-Type: application/json" \
     -d '{
       "to": "user@example.com",
       "subject": "Welcome!",
       "body": "Welcome to our service!",
       "priority": "high",
       "user_id": "123"
     }'
   ```

### Manual Setup

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Start Redis**
   ```bash
   # Using Docker
   docker run -d -p 6379:6379 redis:7-alpine
   
   # Or install Redis locally
   # brew install redis (macOS)
   # sudo apt-get install redis-server (Ubuntu)
   ```

3. **Configure the service**
   ```bash
   cp config.env.example config.env
   # Edit config.env with your settings
   ```

4. **Build and run**
   ```bash
   # Build both server and worker
   make build
   
   # Run server (in one terminal)
   make run-server
   
   # Run worker (in another terminal)
   make run-worker
   ```

## API Endpoints

### Health Check
```
GET /health
```

### Email Notifications
```
POST /api/v1/notifications/email
Content-Type: application/json

{
  "to": "user@example.com",
  "subject": "Welcome!",
  "body": "Welcome to our service!",
  "templateId": "welcome_template",
  "data": {"name": "John"},
  "priority": "high",
  "userId": "123"
}
```

### SMS Notifications
```
POST /api/v1/notifications/sms
Content-Type: application/json

{
  "to": "+1234567890",
  "message": "Your verification code is 123456",
  "priority": "high",
  "userId": "123"
}
```

### Push Notifications
```
POST /api/v1/notifications/push
Content-Type: application/json

{
  "deviceToken": "fcm_token_here",
  "title": "New Message",
  "body": "You have a new message",
  "data": {"messageId": "456"},
  "priority": "default",
  "userId": "123"
}
```

### Webhook Notifications
```
POST /api/v1/notifications/webhook
Content-Type: application/json

{
  "url": "https://api.example.com/webhook",
  "method": "POST",
  "headers": {"Authorization": "Bearer token"},
  "body": {"event": "user_registered"},
  "priority": "low",
  "userId": "123"
}
```

### Bulk Notifications
```
POST /api/v1/notifications/bulk
Content-Type: application/json

{
  "type": "email",
  "recipients": ["user1@example.com", "user2@example.com"],
  "subject": "Bulk Notification",
  "message": "This is a bulk notification",
  "priority": "default",
  "userId": "123"
}
```

## Configuration

The service can be configured using environment variables or the `config.env` file:

```env
# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Asynq Configuration
ASYNQ_CONCURRENCY=10
ASYNQ_QUEUE_CRITICAL=6
ASYNQ_QUEUE_DEFAULT=3
ASYNQ_QUEUE_LOW=1

# Notification Configuration
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_FROM=noreply@example.com
EMAIL_PASSWORD=your-email-password

# Logging
LOG_LEVEL=info
```

## Queue Priorities

The service uses three priority queues:

- **Critical** (6 workers): High-priority notifications
- **Default** (3 workers): Standard notifications
- **Low** (1 worker): Low-priority notifications

## Monitoring

### Using Asynqmon

Asynqmon provides a web-based dashboard for monitoring queues and tasks:

1. **Enable in docker-compose.yml**
   ```yaml
   asynqmon:
     image: hibiken/asynqmon:latest
     ports:
       - "8081:8080"
     environment:
       - REDIS_ADDR=redis:6379
   ```

2. **Access the dashboard**
   ```
   http://localhost:8081
   ```

### Using Asynq CLI

Install the Asynq CLI tool for command-line monitoring:

```bash
make install-asynq-cli

# View queue statistics
asynq stats --redis-addr=localhost:6379

# View tasks in a queue
asynq queue --redis-addr=localhost:6379 default

# View task details
asynq task --redis-addr=localhost:6379 <task-id>
```

## Development

### Project Structure

```
notification-svc/
├── cmd/
│   ├── server/          # HTTP server entry point
│   └── worker/          # Worker entry point
├── internal/
│   ├── client/          # Asynq client
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── server/          # HTTP server
│   ├── tasks/           # Task definitions and handlers
│   └── worker/          # Worker server
├── config.env           # Configuration file
├── docker-compose.yml   # Docker services
├── Dockerfile           # Container definition
├── go.mod              # Go modules
├── Makefile            # Build and run commands
└── README.md           # This file
```

### Building

```bash
# Build both binaries
make build

# Build server only
make build-server

# Build worker only
make build-worker
```

### Running

```bash
# Run both server and worker (requires tmux)
make run

# Run server only
make run-server

# Run worker only
make run-worker
```

### Testing

```bash
# Run tests
make test
```

## Scaling

### Horizontal Scaling

The service is designed for horizontal scaling:

1. **Multiple Workers**: Run multiple worker instances to process more tasks
2. **Load Balancing**: Use a load balancer in front of multiple server instances
3. **Redis Cluster**: Use Redis Cluster for high availability

### Example: Multiple Workers

```bash
# Start multiple worker instances
make run-worker  # Terminal 1
make run-worker  # Terminal 2
make run-worker  # Terminal 3
```

## Production Deployment

### Environment Variables

Set appropriate environment variables for production:

```bash
export REDIS_ADDR=your-redis-cluster:6379
export REDIS_PASSWORD=your-redis-password
export ASYNQ_CONCURRENCY=50
export LOG_LEVEL=warn
```

### Security Considerations

1. **Redis Security**: Use Redis ACLs and strong passwords
2. **API Security**: Implement authentication and rate limiting
3. **Network Security**: Use TLS for all communications
4. **Container Security**: Run containers as non-root users

### Monitoring and Logging

1. **Application Metrics**: Implement Prometheus metrics
2. **Logging**: Use structured logging with correlation IDs
3. **Health Checks**: Implement comprehensive health checks
4. **Alerting**: Set up alerts for queue backlogs and failures

## Troubleshooting

### Common Issues

1. **Redis Connection Failed**
   - Check Redis is running: `redis-cli ping`
   - Verify connection settings in config.env

2. **Tasks Not Processing**
   - Ensure worker is running: `make run-worker`
   - Check Redis connection from worker
   - Verify task handlers are registered

3. **High Memory Usage**
   - Adjust `ASYNQ_CONCURRENCY` setting
   - Monitor Redis memory usage
   - Implement task cleanup

### Debug Mode

Enable debug logging by setting `LOG_LEVEL=debug` in your configuration.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## References

- [Asynq Documentation](https://github.com/hibiken/asynq)
- [Redis Documentation](https://redis.io/documentation)
- [Gin Framework](https://gin-gonic.com/) 