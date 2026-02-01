package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"realtime-events/internal/api/handlers"
	"realtime-events/internal/config"
	"realtime-events/internal/middleware"
	"realtime-events/internal/observability"
	"realtime-events/internal/services"
	"realtime-events/pkg/queue"
	"realtime-events/pkg/storage"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		sugar.Fatalw("Failed to load config", "error", err)
	}

	// Initialize storage
	db, err := storage.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		sugar.Fatalw("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Initialize queue
	eventQueue, err := queue.NewRedisQueue(cfg.RedisURL, "events")
	if err != nil {
		sugar.Fatalw("Failed to connect to queue", "error", err)
	}
	defer eventQueue.Close()

	// Initialize services
	eventService := services.NewEventService(db, eventQueue, sugar)
	healthService := services.NewHealthService(db, sugar)

	// Initialize handlers
	eventHandler := handlers.NewEventHandler(eventService, sugar)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middleware
	router.Use(middleware.Logger(sugar))
	router.Use(middleware.Recovery(sugar))
	router.Use(middleware.RateLimit())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		health := healthService.CheckHealth(c.Request.Context())
		statusCode := 200
		if health.Status != "healthy" {
			statusCode = 503
		}
		c.JSON(statusCode, health)
	})

	// Metrics
	router.GET("/metrics", gin.WrapH(observability.MetricsHandler()))

	// API routes
	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthRequired())
	{
		v1.POST("/events", eventHandler.IngestEvent)
		v1.POST("/events/batch", eventHandler.IngestBatchEvents)
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		sugar.Infow("Starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalw("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugar.Info("Shutting down server...")

	// Context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatalw("Server forced to shutdown", "error", err)
	}

	sugar.Info("Server exited")
}
