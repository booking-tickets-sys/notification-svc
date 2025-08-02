package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"notification-svc/internal/client"
	"notification-svc/internal/config"
	"notification-svc/internal/handlers"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	router *gin.Engine
	server *http.Server
	client *client.Client
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, client *client.Client) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := &Server{
		config: cfg,
		router: router,
		client: client,
	}

	server.setupRoutes()

	return server
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	// Create handler
	handler := handlers.NewNotificationHandler(s.client)

	// Health check
	s.router.GET("/health", handler.HealthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Notification routes
		notifications := v1.Group("/notifications")
		{
			notifications.POST("/email", handler.SendEmailNotification)
			notifications.POST("/sms", handler.SendSMSNotification)
			notifications.POST("/push", handler.SendPushNotification)
			notifications.POST("/webhook", handler.SendWebhookNotification)
			notifications.POST("/bulk", handler.SendBulkNotification)
		}
	}

	// Add CORS middleware
	s.router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting HTTP server on %s", addr)
	return s.server.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop() error {
	log.Println("Stopping HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
