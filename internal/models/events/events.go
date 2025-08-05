package events

type EventMetadata struct {
	EventID     string `json:"eventID"`
	EventName   string `json:"eventName"`
	PublishedAt int64 `json:"publishedAt"`
}

type EventType string

const (
	OrderCreatedEventType EventType = "order_created"
	LoginEventType EventType = "login"
)
