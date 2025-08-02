package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

// Task types
const (
	TypeEmailNotification   = "email:notification"
	TypeSMSNotification     = "sms:notification"
	TypePushNotification    = "push:notification"
	TypeWebhookNotification = "webhook:notification"
)

// EmailNotificationPayload represents the payload for email notifications
type EmailNotificationPayload struct {
	To         string            `json:"to"`
	Subject    string            `json:"subject"`
	Body       string            `json:"body"`
	TemplateID string            `json:"templateId,omitempty"`
	Data       map[string]string `json:"data,omitempty"`
	Priority   string            `json:"priority"`
	UserID     string            `json:"userId"`
}

// SMSNotificationPayload represents the payload for SMS notifications
type SMSNotificationPayload struct {
	To       string `json:"to"`
	Message  string `json:"message"`
	Priority string `json:"priority"`
	UserID   string `json:"userId"`
}

// PushNotificationPayload represents the payload for push notifications
type PushNotificationPayload struct {
	DeviceToken string            `json:"deviceToken"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Data        map[string]string `json:"data,omitempty"`
	Priority    string            `json:"priority"`
	UserID      string            `json:"userId"`
}

// WebhookNotificationPayload represents the payload for webhook notifications
type WebhookNotificationPayload struct {
	URL      string            `json:"url"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers,omitempty"`
	Body     interface{}       `json:"body"`
	Priority string            `json:"priority"`
	UserID   string            `json:"userId"`
}

// NewEmailNotificationTask creates a new email notification task
func NewEmailNotificationTask(payload EmailNotificationPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Timeout(30 * time.Second),
	}

	// Set queue based on priority
	switch payload.Priority {
	case "high":
		opts = append(opts, asynq.Queue("critical"))
	case "low":
		opts = append(opts, asynq.Queue("low"))
	default:
		opts = append(opts, asynq.Queue("default"))
	}

	return asynq.NewTask(TypeEmailNotification, data, opts...), nil
}

// NewSMSNotificationTask creates a new SMS notification task
func NewSMSNotificationTask(payload SMSNotificationPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Timeout(15 * time.Second),
	}

	switch payload.Priority {
	case "high":
		opts = append(opts, asynq.Queue("critical"))
	case "low":
		opts = append(opts, asynq.Queue("low"))
	default:
		opts = append(opts, asynq.Queue("default"))
	}

	return asynq.NewTask(TypeSMSNotification, data, opts...), nil
}

// NewPushNotificationTask creates a new push notification task
func NewPushNotificationTask(payload PushNotificationPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Timeout(20 * time.Second),
	}

	switch payload.Priority {
	case "high":
		opts = append(opts, asynq.Queue("critical"))
	case "low":
		opts = append(opts, asynq.Queue("low"))
	default:
		opts = append(opts, asynq.Queue("default"))
	}

	return asynq.NewTask(TypePushNotification, data, opts...), nil
}

// NewWebhookNotificationTask creates a new webhook notification task
func NewWebhookNotificationTask(payload WebhookNotificationPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	opts := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.Timeout(60 * time.Second),
	}

	switch payload.Priority {
	case "high":
		opts = append(opts, asynq.Queue("critical"))
	case "low":
		opts = append(opts, asynq.Queue("low"))
	default:
		opts = append(opts, asynq.Queue("default"))
	}

	return asynq.NewTask(TypeWebhookNotification, data, opts...), nil
}

// HandleEmailNotificationTask handles email notification tasks
func HandleEmailNotificationTask(ctx context.Context, t *asynq.Task) error {
	var p EmailNotificationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Processing email notification: to=%s, subject=%s, user_id=%s", p.To, p.Subject, p.UserID)

	// TODO: Implement actual email sending logic
	// For now, just simulate email sending
	time.Sleep(2 * time.Second)

	log.Printf("Email notification sent successfully: to=%s", p.To)
	return nil
}

// HandleSMSNotificationTask handles SMS notification tasks
func HandleSMSNotificationTask(ctx context.Context, t *asynq.Task) error {
	var p SMSNotificationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Processing SMS notification: to=%s, user_id=%s", p.To, p.UserID)

	// TODO: Implement actual SMS sending logic
	// For now, just simulate SMS sending
	time.Sleep(1 * time.Second)

	log.Printf("SMS notification sent successfully: to=%s", p.To)
	return nil
}

// HandlePushNotificationTask handles push notification tasks
func HandlePushNotificationTask(ctx context.Context, t *asynq.Task) error {
	var p PushNotificationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Processing push notification: device_token=%s, title=%s, user_id=%s", p.DeviceToken, p.Title, p.UserID)

	// TODO: Implement actual push notification logic
	// For now, just simulate push notification sending
	time.Sleep(1 * time.Second)

	log.Printf("Push notification sent successfully: device_token=%s", p.DeviceToken)
	return nil
}

// HandleWebhookNotificationTask handles webhook notification tasks
func HandleWebhookNotificationTask(ctx context.Context, t *asynq.Task) error {
	var p WebhookNotificationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Processing webhook notification: url=%s, method=%s, user_id=%s", p.URL, p.Method, p.UserID)

	// TODO: Implement actual webhook sending logic
	// For now, just simulate webhook call
	time.Sleep(3 * time.Second)

	log.Printf("Webhook notification sent successfully: url=%s", p.URL)
	return nil
}
