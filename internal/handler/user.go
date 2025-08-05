package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"

	"notification-svc/internal/models/events"
)

type UserHandler struct {
	logger *logrus.Logger
}

func NewUserHandler(logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		logger: logger,
	}
}

func (h *UserHandler) HandleLoginEvent(ctx context.Context, task *asynq.Task) error {
	var loginEvent events.LoginEvent
	if err := json.Unmarshal(task.Payload(), &loginEvent); err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": loginEvent.UserID,
		"email":    loginEvent.Email,
		"username": loginEvent.Username,
		"login_at": loginEvent.LoginAt.Format(time.RFC3339),
	}).Info("Login event received")

	// TODO: Implement login event handling

	return nil
}