package worker

import (
	"log"

	"notification-svc/internal/config"
	"notification-svc/internal/tasks"

	"github.com/hibiken/asynq"
)

// Worker represents the Asynq worker server
type Worker struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

// NewWorker creates a new worker instance
func NewWorker(cfg *config.Config) *Worker {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: cfg.Asynq.Concurrency,
			Queues:      cfg.Asynq.Queues,
		},
	)

	mux := asynq.NewServeMux()

	// Register task handlers
	mux.HandleFunc(tasks.TypeEmailNotification, tasks.HandleEmailNotificationTask)
	mux.HandleFunc(tasks.TypeSMSNotification, tasks.HandleSMSNotificationTask)
	mux.HandleFunc(tasks.TypePushNotification, tasks.HandlePushNotificationTask)
	mux.HandleFunc(tasks.TypeWebhookNotification, tasks.HandleWebhookNotificationTask)

	return &Worker{
		server: server,
		mux:    mux,
	}
}

// Start starts the worker server
func (w *Worker) Start() error {
	log.Println("Starting notification worker...")
	return w.server.Run(w.mux)
}

// Stop gracefully stops the worker server
func (w *Worker) Stop() {
	log.Println("Stopping notification worker...")
	w.server.Stop()
	w.server.Shutdown()
}
