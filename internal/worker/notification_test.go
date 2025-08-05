package worker

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"notification-svc/internal/models/events"
)

func TestNotificationWorker_New(t *testing.T) {
	logger := logrus.New()
	cfg := Config{
		RedisAddr:   "localhost:6379",
		Concurrency: 5,
	}

	worker := NewNotificationWorker(cfg, logger)
	assert.NotNil(t, worker)
	assert.Equal(t, cfg, worker.cfg)
	assert.Equal(t, logger, worker.logger)
	assert.False(t, worker.IsRunning())
}

func TestOrderCreatedPayload_ToTask(t *testing.T) {
	payload := events.OrderCreatedPayload{
		EventMetadata: events.EventMetadata{
			EventID:     "test-event-123",
			EventName:   "order_created",
			PublishedAt: time.Now().Unix(),
		},
		OrderID:  "order-123",
		UserID:   "user-456",
		TicketID: "ticket-789",
	}

	task, err := payload.ToTask()
	require.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "order_created", task.Type())
} 