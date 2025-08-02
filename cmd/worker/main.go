package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"notification-svc/internal/config"
	"notification-svc/internal/worker"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create worker
	w := worker.NewWorker(cfg)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v", sig)
		cancel()
	}()

	// Start worker in a goroutine
	go func() {
		if err := w.Start(); err != nil {
			log.Printf("Worker error: %v", err)
			cancel()
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	log.Println("Shutting down worker...")
	w.Stop()
	log.Println("Worker stopped")
}
