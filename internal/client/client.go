package client

import (
	"log"

	"notification-svc/internal/config"
	"notification-svc/internal/tasks"

	"github.com/hibiken/asynq"
)

// Client represents the Asynq client for enqueueing tasks
type Client struct {
	client *asynq.Client
}

// NewClient creates a new Asynq client
func NewClient(cfg *config.Config) *Client {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	client := asynq.NewClient(redisOpt)

	return &Client{
		client: client,
	}
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.client.Close()
}

// EnqueueEmailNotification enqueues an email notification task
func (c *Client) EnqueueEmailNotification(payload tasks.EmailNotificationPayload) error {
	task, err := tasks.NewEmailNotificationTask(payload)
	if err != nil {
		return err
	}

	info, err := c.client.Enqueue(task)
	if err != nil {
		return err
	}

	log.Printf("Enqueued email notification task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueueSMSNotification enqueues an SMS notification task
func (c *Client) EnqueueSMSNotification(payload tasks.SMSNotificationPayload) error {
	task, err := tasks.NewSMSNotificationTask(payload)
	if err != nil {
		return err
	}

	info, err := c.client.Enqueue(task)
	if err != nil {
		return err
	}

	log.Printf("Enqueued SMS notification task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueuePushNotification enqueues a push notification task
func (c *Client) EnqueuePushNotification(payload tasks.PushNotificationPayload) error {
	task, err := tasks.NewPushNotificationTask(payload)
	if err != nil {
		return err
	}

	info, err := c.client.Enqueue(task)
	if err != nil {
		return err
	}

	log.Printf("Enqueued push notification task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueueWebhookNotification enqueues a webhook notification task
func (c *Client) EnqueueWebhookNotification(payload tasks.WebhookNotificationPayload) error {
	task, err := tasks.NewWebhookNotificationTask(payload)
	if err != nil {
		return err
	}

	info, err := c.client.Enqueue(task)
	if err != nil {
		return err
	}

	log.Printf("Enqueued webhook notification task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}
