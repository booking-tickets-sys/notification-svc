package events

// OrderCreatedPayload represents the payload for order_created events
type OrderCreatedPayload struct {
	EventMetadata EventMetadata `json:"eventMetadata"`
	OrderID       string        `json:"orderId"`
	UserID        string        `json:"userId"`
	TicketID      string        `json:"ticketId"`
}