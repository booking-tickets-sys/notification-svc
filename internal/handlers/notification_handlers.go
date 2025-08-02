package handlers

import (
	"net/http"

	"notification-svc/internal/client"
	"notification-svc/internal/tasks"

	"github.com/gin-gonic/gin"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	client *client.Client
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(client *client.Client) *NotificationHandler {
	return &NotificationHandler{
		client: client,
	}
}

// SendEmailNotification handles email notification requests
func (h *NotificationHandler) SendEmailNotification(c *gin.Context) {
	var request struct {
		To         string            `json:"to" binding:"required,email"`
		Subject    string            `json:"subject" binding:"required"`
		Body       string            `json:"body" binding:"required"`
		TemplateID string            `json:"templateId,omitempty"`
		Data       map[string]string `json:"data,omitempty"`
		Priority   string            `json:"priority,omitempty"`
		UserID     string            `json:"userId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default priority if not provided
	if request.Priority == "" {
		request.Priority = "default"
	}

	payload := tasks.EmailNotificationPayload{
		To:         request.To,
		Subject:    request.Subject,
		Body:       request.Body,
		TemplateID: request.TemplateID,
		Data:       request.Data,
		Priority:   request.Priority,
		UserID:     request.UserID,
	}

	if err := h.client.EnqueueEmailNotification(payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue email notification"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Email notification queued successfully",
		"type":    "email",
		"to":      request.To,
	})
}

// SendSMSNotification handles SMS notification requests
func (h *NotificationHandler) SendSMSNotification(c *gin.Context) {
	var request struct {
		To       string `json:"to" binding:"required"`
		Message  string `json:"message" binding:"required"`
		Priority string `json:"priority,omitempty"`
		UserID   string `json:"userId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default priority if not provided
	if request.Priority == "" {
		request.Priority = "default"
	}

	payload := tasks.SMSNotificationPayload{
		To:       request.To,
		Message:  request.Message,
		Priority: request.Priority,
		UserID:   request.UserID,
	}

	if err := h.client.EnqueueSMSNotification(payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue SMS notification"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "SMS notification queued successfully",
		"type":    "sms",
		"to":      request.To,
	})
}

// SendPushNotification handles push notification requests
func (h *NotificationHandler) SendPushNotification(c *gin.Context) {
	var request struct {
		DeviceToken string            `json:"deviceToken" binding:"required"`
		Title       string            `json:"title" binding:"required"`
		Body        string            `json:"body" binding:"required"`
		Data        map[string]string `json:"data,omitempty"`
		Priority    string            `json:"priority,omitempty"`
		UserID      string            `json:"userId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default priority if not provided
	if request.Priority == "" {
		request.Priority = "default"
	}

	payload := tasks.PushNotificationPayload{
		DeviceToken: request.DeviceToken,
		Title:       request.Title,
		Body:        request.Body,
		Data:        request.Data,
		Priority:    request.Priority,
		UserID:      request.UserID,
	}

	if err := h.client.EnqueuePushNotification(payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue push notification"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":     "Push notification queued successfully",
		"type":        "push",
		"deviceToken": request.DeviceToken,
	})
}

// SendWebhookNotification handles webhook notification requests
func (h *NotificationHandler) SendWebhookNotification(c *gin.Context) {
	var request struct {
		URL      string            `json:"url" binding:"required"`
		Method   string            `json:"method" binding:"required"`
		Headers  map[string]string `json:"headers,omitempty"`
		Body     interface{}       `json:"body" binding:"required"`
		Priority string            `json:"priority,omitempty"`
		UserID   string            `json:"userId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default priority if not provided
	if request.Priority == "" {
		request.Priority = "default"
	}

	// Set default method if not provided
	if request.Method == "" {
		request.Method = "POST"
	}

	payload := tasks.WebhookNotificationPayload{
		URL:      request.URL,
		Method:   request.Method,
		Headers:  request.Headers,
		Body:     request.Body,
		Priority: request.Priority,
		UserID:   request.UserID,
	}

	if err := h.client.EnqueueWebhookNotification(payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue webhook notification"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Webhook notification queued successfully",
		"type":    "webhook",
		"url":     request.URL,
	})
}

// SendBulkNotification handles bulk notification requests
func (h *NotificationHandler) SendBulkNotification(c *gin.Context) {
	var request struct {
		Type       string            `json:"type" binding:"required,oneof=email sms push webhook"`
		Recipients []string          `json:"recipients" binding:"required"`
		Subject    string            `json:"subject,omitempty"`
		Message    string            `json:"message" binding:"required"`
		Priority   string            `json:"priority,omitempty"`
		UserID     string            `json:"userId" binding:"required"`
		Data       map[string]string `json:"data,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default priority if not provided
	if request.Priority == "" {
		request.Priority = "default"
	}

	successCount := 0
	errorCount := 0

	for _, recipient := range request.Recipients {
		var err error

		switch request.Type {
		case "email":
			payload := tasks.EmailNotificationPayload{
				To:       recipient,
				Subject:  request.Subject,
				Body:     request.Message,
				Data:     request.Data,
				Priority: request.Priority,
				UserID:   request.UserID,
			}
			err = h.client.EnqueueEmailNotification(payload)

		case "sms":
			payload := tasks.SMSNotificationPayload{
				To:       recipient,
				Message:  request.Message,
				Priority: request.Priority,
				UserID:   request.UserID,
			}
			err = h.client.EnqueueSMSNotification(payload)

		case "push":
			payload := tasks.PushNotificationPayload{
				DeviceToken: recipient,
				Title:       request.Subject,
				Body:        request.Message,
				Data:        request.Data,
				Priority:    request.Priority,
				UserID:      request.UserID,
			}
			err = h.client.EnqueuePushNotification(payload)

		case "webhook":
			payload := tasks.WebhookNotificationPayload{
				URL:      recipient,
				Method:   "POST",
				Headers:  request.Data,
				Body:     request.Message,
				Priority: request.Priority,
				UserID:   request.UserID,
			}
			err = h.client.EnqueueWebhookNotification(payload)
		}

		if err != nil {
			errorCount++
		} else {
			successCount++
		}
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Bulk notification queued successfully",
		"type":    request.Type,
		"total":   len(request.Recipients),
		"success": successCount,
		"errors":  errorCount,
	})
}

// HealthCheck handles health check requests
func (h *NotificationHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "notification-service",
	})
}
