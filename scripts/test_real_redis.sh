#!/bin/bash

echo "🧪 Testing Notification Service with Real Redis"

# Check if Redis is running
if ! redis-cli ping > /dev/null 2>&1; then
    echo "❌ Redis is not running. Please start Redis first."
    exit 1
fi

echo "✅ Redis is running"

# Clean up any existing test data
echo "🧹 Cleaning up existing test data..."
redis-cli DEL asynq:default asynq:processing asynq:failed asynq:retry asynq:deadline

# Start the notification service
echo "🚀 Starting notification service..."
go run main.go &
SERVICE_PID=$!

# Wait for service to start
sleep 3

echo "📋 Service is running (PID: $SERVICE_PID)"

# Create a simple Go program to enqueue tasks
cat > enqueue_tasks.go << 'EOF'
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"notification-svc/models/events"
)

func main() {
	// Create Asynq client
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: "localhost:6379",
	})
	defer client.Close()

	// Enqueue multiple tasks
	for i := 0; i < 5; i++ {
		payload := events.OrderCreatedPayload{
			EventMetadata: events.EventMetadata{
				EventID:     fmt.Sprintf("test-event-%d", i),
				EventName:   "order_created",
				PublishedAt: time.Now().Unix(),
			},
			OrderID:  fmt.Sprintf("order-%d", i),
			UserID:   fmt.Sprintf("user-%d", i),
			TicketID: fmt.Sprintf("ticket-%d", i),
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Failed to marshal payload %d: %v", i, err)
			continue
		}
		task := asynq.NewTask("order_created", payloadBytes)

		info, err := client.EnqueueContext(context.Background(), task)
		if err != nil {
			log.Printf("Failed to enqueue task %d: %v", i, err)
			continue
		}

		fmt.Printf("Task %d enqueued with ID: %s\n", i, info.ID)
	}

	fmt.Println("All tasks enqueued successfully!")
}
EOF

# Test 1: Enqueue tasks using proper Asynq client
echo ""
echo "📧 Test 1: Enqueueing tasks using Asynq client..."
go run enqueue_tasks.go

# Wait for processing
sleep 3

# Check Redis queue status
echo ""
echo "📊 Redis Queue Status:"
echo "Pending tasks: $(redis-cli LLEN asynq:default)"
echo "Processing tasks: $(redis-cli LLEN asynq:processing)"
echo "Failed tasks: $(redis-cli LLEN asynq:failed)"
echo "Retry tasks: $(redis-cli LLEN asynq:retry)"

# Wait a bit more for processing
sleep 2

# Check final status
echo ""
echo "📊 Final Queue Status:"
echo "Pending tasks: $(redis-cli LLEN asynq:default)"
echo "Processing tasks: $(redis-cli LLEN asynq:processing)"
echo "Failed tasks: $(redis-cli LLEN asynq:failed)"
echo "Retry tasks: $(redis-cli LLEN asynq:retry)"

# Stop the service
echo ""
echo "🛑 Stopping notification service..."
kill -TERM $SERVICE_PID

# Wait for graceful shutdown
wait $SERVICE_PID

# Cleanup
echo ""
echo "🧹 Cleaning up..."
rm -f enqueue_tasks.go

echo ""
echo "🎉 Real Redis testing completed!"
echo ""
echo "💡 Expected behavior:"
echo "   - Service should process all valid tasks"
echo "   - All queues should be empty after processing"
echo "   - Service should shutdown gracefully" 