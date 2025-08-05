#!/bin/bash

echo "🧪 Testing Graceful Shutdown"

# Check if Redis is running
if ! redis-cli ping > /dev/null 2>&1; then
    echo "❌ Redis is not running. Please start Redis first."
    exit 1
fi

echo "✅ Redis is running"

# Start the notification service
echo "🚀 Starting notification service..."
go run main.go &
SERVICE_PID=$!

# Wait for service to start
sleep 3

echo "📋 Service is running (PID: $SERVICE_PID)"
echo "⏰ Sending SIGTERM for graceful shutdown..."

# Send SIGTERM for graceful shutdown
kill -TERM $SERVICE_PID

# Wait for graceful shutdown
echo "⏳ Waiting for graceful shutdown..."
wait $SERVICE_PID

echo "✅ Graceful shutdown completed!"
echo ""
echo "💡 Expected behavior:"
echo "   - Service should log 'Shutting down notification service...'"
echo "   - Asynq server should log 'Starting graceful shutdown'"
echo "   - Service should log 'NotificationWorker stopped gracefully'"
echo "   - Service should exit cleanly" 