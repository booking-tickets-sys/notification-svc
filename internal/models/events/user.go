package events

import (
	"time"
)

type LoginEvent struct {
	EventMetadata EventMetadata `json:"eventMetadata"`
	UserID        string        `json:"userId"`
	Email         string        `json:"email"`
	Username      string        `json:"username"`
	LoginAt       time.Time     `json:"loginAt"`
}