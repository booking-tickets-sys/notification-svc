#!/bin/bash

# Test script for notification service
echo "🧪 Testing Notification Service with Asynq"

# Check if Redis is running
echo "📡 Checking Redis connection..."
if ! redis-cli ping > /dev/null 2>&1; then
    echo "❌ Redis is not running. Please start Redis first:"
    echo "   brew services start redis"
    echo "   or"
    echo "   docker run -d -p 6379:6379 redis:alpine"
    exit 1
fi
echo "✅ Redis is running"

# Start the notification service in background
echo "🚀 Starting notification service..."
go run main.go &
SERVICE_PID=$!

# Wait for service to start
sleep 2

# Test 1: Send a test notification task
echo ""
echo "📧 Test 1: Sending order_created notification task..."

# Create test payload
cat > /tmp/test_payload.json << EOF
{
  "eventMetadata": {
    "eventID": "test-event-$(date +%s)",
    "eventName": "order_created",
    "publishedAt": $(date +%s)
  },
  "orderId": "order-$(date +%s)",
  "userId": "user-123",
  "ticketId": "ticket-456"
}
EOF

# Enqueue task using Asynq CLI (if available) or Redis directly
echo "📤 Enqueueing task to Redis..."
redis-cli LPUSH asynq:default '{"type":"order_created","payload":"'$(cat /tmp/test_payload.json | tr -d '\n' | sed 's/"/\\"/g')'","retry":0,"queue":"default","deadline":0}'

echo "✅ Task enqueued successfully"

# Test 2: Check Redis queue status
echo ""
echo "📊 Test 2: Checking queue status..."
echo "Pending tasks:"
redis-cli LLEN asynq:default

echo "Processing tasks:"
redis-cli LLEN asynq:processing

echo "Failed tasks:"
redis-cli LLEN asynq:failed

# Test 3: Monitor logs for processing
echo ""
echo "📋 Test 3: Monitoring service logs..."
echo "Check the service logs above for processing messages"
echo "You should see: 'Processing order_created event' and 'Sending notification for order: ...'"

# Wait a bit for processing
sleep 3

# Test 4: Check final queue status
echo ""
echo "📊 Test 4: Final queue status..."
echo "Pending tasks:"
redis-cli LLEN asynq:default

echo "Processing tasks:"
redis-cli LLEN asynq:processing

echo "Failed tasks:"
redis-cli LLEN asynq:failed

# Cleanup
echo ""
echo "🧹 Cleaning up..."
kill $SERVICE_PID 2>/dev/null
rm -f /tmp/test_payload.json

echo ""
echo "🎉 Testing completed!"
echo ""
echo "💡 To test manually:"
echo "1. Start the service: go run main.go"
echo "2. In another terminal, enqueue tasks:"
echo "   redis-cli LPUSH asynq:default '{\"type\":\"order_created\",\"payload\":\"...\",\"retry\":0,\"queue\":\"default\"}'"
echo "3. Watch the service logs for processing" 