#!/bin/bash

# Test script for the notification service
# Make sure the service is running on localhost:8080

BASE_URL="http://localhost:8080/api/v1"

echo "🧪 Testing Notification Service"
echo "================================"

# Test health check
echo "1. Testing health check..."
curl -s -X GET "$BASE_URL/../health" | jq .
echo ""

# Test email notification
echo "2. Testing email notification..."
curl -s -X POST "$BASE_URL/notifications/email" \
  -H "Content-Type: application/json" \
  -d '{
    "to": "test@example.com",
    "subject": "Welcome to our service!",
    "body": "Thank you for joining us. We are excited to have you on board!",
    "templateId": "welcome_template",
    "data": {"name": "John Doe", "company": "Example Corp"},
    "priority": "high",
    "userId": "user_123"
  }' | jq .
echo ""

# Test SMS notification
echo "3. Testing SMS notification..."
curl -s -X POST "$BASE_URL/notifications/sms" \
  -H "Content-Type: application/json" \
  -d '{
    "to": "+1234567890",
    "message": "Your verification code is 123456. Valid for 10 minutes.",
    "priority": "high",
    "userId": "user_123"
  }' | jq .
echo ""

# Test push notification
echo "4. Testing push notification..."
curl -s -X POST "$BASE_URL/notifications/push" \
  -H "Content-Type: application/json" \
  -d '{
    "deviceToken": "fcm_token_example_123456789",
    "title": "New Message",
    "body": "You have received a new message from John",
    "data": {"messageId": "msg_456", "sender": "John Doe"},
    "priority": "default",
    "userId": "user_123"
  }' | jq .
echo ""

# Test webhook notification
echo "5. Testing webhook notification..."
curl -s -X POST "$BASE_URL/notifications/webhook" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://httpbin.org/post",
    "method": "POST",
    "headers": {"Authorization": "Bearer test_token", "X-Custom-Header": "test_value"},
    "body": {"event": "user_registered", "userId": "user_123", "timestamp": "2024-01-01T00:00:00Z"},
    "priority": "low",
    "userId": "user_123"
  }' | jq .
echo ""

# Test bulk email notifications
echo "6. Testing bulk email notifications..."
curl -s -X POST "$BASE_URL/notifications/bulk" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "recipients": ["user1@example.com", "user2@example.com", "user3@example.com"],
    "subject": "Important Announcement",
    "message": "We have an important announcement for all users. Please check your dashboard for details.",
    "priority": "default",
    "userId": "admin_456"
  }' | jq .
echo ""

# Test bulk SMS notifications
echo "7. Testing bulk SMS notifications..."
curl -s -X POST "$BASE_URL/notifications/bulk" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "sms",
    "recipients": ["+1234567890", "+0987654321"],
    "message": "Emergency alert: System maintenance scheduled for tonight at 2 AM.",
    "priority": "high",
    "userId": "admin_456"
  }' | jq .
echo ""

# Test different priority levels
echo "8. Testing different priority levels..."

echo "   - High priority email..."
curl -s -X POST "$BASE_URL/notifications/email" \
  -H "Content-Type: application/json" \
  -d '{
    "to": "urgent@example.com",
    "subject": "URGENT: System Alert",
    "body": "Critical system issue detected. Immediate attention required.",
    "priority": "high",
    "userId": "system_monitor"
  }' | jq .

echo "   - Low priority email..."
curl -s -X POST "$BASE_URL/notifications/email" \
  -H "Content-Type: application/json" \
  -d '{
    "to": "newsletter@example.com",
    "subject": "Weekly Newsletter",
    "body": "Here is your weekly newsletter with the latest updates.",
    "priority": "low",
    "userId": "newsletter_system"
  }' | jq .
echo ""

echo "✅ All tests completed!"
echo ""
echo "📊 To monitor the queues, you can:"
echo "   - Use Asynqmon: http://localhost:8081 (if enabled)"
echo "   - Use Asynq CLI: asynq stats --redis-addr=localhost:6379"
echo ""
echo "🔍 Check the worker logs to see task processing:"
echo "   make run-worker" 