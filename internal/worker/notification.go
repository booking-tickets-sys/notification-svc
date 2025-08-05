package worker

import (
	"context"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"

	"notification-svc/internal/handler"
	"notification-svc/internal/models/events"
)

// Config holds worker configuration
type Config struct {
	RedisAddr   string
	Concurrency int
}

// NotificationWorker listens to the Redis stream and enqueues Asynq tasks
type NotificationWorker struct {
	cfg       Config
	logger    *logrus.Logger
	server    *asynq.Server
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	mu        sync.RWMutex
	isRunning bool
}

// NewNotificationWorker creates a new notification worker
func NewNotificationWorker(cfg Config, logger *logrus.Logger) *NotificationWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &NotificationWorker{
		cfg:    cfg,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start starts the notification worker
func (w *NotificationWorker) Start() {
	w.mu.Lock()
	if w.isRunning {
		w.mu.Unlock()
		w.logger.Warn("NotificationWorker is already running")
		return
	}
	w.isRunning = true
	w.mu.Unlock()

	w.logger.Info("NotificationWorker started")

	w.server = asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: w.cfg.RedisAddr,
		},
		asynq.Config{
			Concurrency:     w.cfg.Concurrency,
			ShutdownTimeout: 30 * time.Second,
		},
	)

	userHandler := handler.NewUserHandler(w.logger)
	orderHandler := handler.NewOrderHandler(w.logger)

	mux := asynq.NewServeMux()
	mux.HandleFunc(string(events.OrderCreatedEventType), orderHandler.HandleOrderCreated)
	mux.HandleFunc(string(events.LoginEventType), userHandler.HandleLoginEvent)

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		if err := w.server.Run(mux); err != nil {
			select {
			case <-w.ctx.Done():
				w.logger.Info("Server stopped due to context cancellation")
			default:
				w.logger.Error("Server error: ", err)
			}
		}
	}()

	<-w.ctx.Done()
	w.logger.Info("NotificationWorker received shutdown signal")
}

// Stop stops the notification worker gracefully
func (w *NotificationWorker) Stop() {
	w.mu.Lock()
	if !w.isRunning {
		w.mu.Unlock()
		w.logger.Warn("NotificationWorker is not running")
		return
	}
	w.isRunning = false
	w.mu.Unlock()

	w.logger.Info("Stopping NotificationWorker...")
	w.cancel()

	if w.server != nil {
		w.logger.Info("Shutting down Asynq server...")
		w.server.Shutdown()
	}

	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		w.logger.Info("NotificationWorker stopped gracefully")
	case <-time.After(10 * time.Second): // Adjusted timeout
		w.logger.Warn("NotificationWorker shutdown timeout reached")
	}
}

// IsRunning returns true if the worker is running
func (w *NotificationWorker) IsRunning() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.isRunning
}
