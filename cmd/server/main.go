package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"notification-svc/internal/client"
	"notification-svc/internal/config"
	"notification-svc/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create Asynq client
	asynqClient := client.NewClient(cfg)
	defer asynqClient.Close()

	// Create HTTP server
	srv := server.NewServer(cfg, asynqClient)

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

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	log.Println("Shutting down server...")
	if err := srv.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}
	log.Println("Server stopped")
}
