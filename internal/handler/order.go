package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"

	"notification-svc/internal/models/events"
)

// OrderHandler handles order-related events
type OrderHandler struct {
	logger *logrus.Logger
}

// NewOrderHandler creates a new order handler with injected logger
func NewOrderHandler(logger *logrus.Logger) *OrderHandler {
	return &OrderHandler{
		logger: logger,
	}
}

// HandleOrderCreated processes order_created events
func (h *OrderHandler) HandleOrderCreated(ctx context.Context, task *asynq.Task) error {
	h.logger.Info("Processing order_created event")

	// parse the task payload
	var payload events.OrderCreatedPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		h.logger.WithError(err).Error("Failed to unmarshal task payload")
		return err
	}

	// send notification
	h.logger.WithFields(logrus.Fields{
		"order_id":  payload.OrderID,
		"user_id":   payload.UserID,
		"ticket_id": payload.TicketID,
	}).Info("Sending notification for order")

	// TODO: Implement actual notification logic
	// For now, just simulate processing time
	time.Sleep(100 * time.Millisecond)

	h.logger.WithField("order_id", payload.OrderID).Info("Notification sent successfully")
	return nil
}