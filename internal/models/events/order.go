package events

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

// OrderCreatedPayload represents the payload for order_created events
type OrderCreatedPayload struct {
	EventMetadata EventMetadata `json:"eventMetadata"`
	OrderID       string        `json:"orderId"`
	UserID        string        `json:"userId"`
	TicketID      string        `json:"ticketId"`
}

// ToTask converts the payload to an Asynq task
func (p *OrderCreatedPayload) ToTask() (*asynq.Task, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(string(OrderCreatedEventType), payload), nil
}