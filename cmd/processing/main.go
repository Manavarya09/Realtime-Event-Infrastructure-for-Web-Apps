package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"realtime-events/internal/config"
	"realtime-events/internal/models"
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

	// Initialize processing service
	processor := services.NewEventProcessor(db, sugar)

	// Start processing
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		sugar.Info("Shutting down processor...")
		cancel()
	}()

	sugar.Info("Starting event processor...")

	// Process events
	if err := eventQueue.ConsumeEvents(ctx, func(event *models.Event) error {
		return processor.ProcessEvent(ctx, event)
	}); err != nil && err != context.Canceled {
		sugar.Fatalw("Processor failed", "error", err)
	}

	sugar.Info("Processor stopped")
}
