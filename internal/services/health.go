package services

import (
	"context"
	"time"

	"go.uber.org/zap"
	"realtime-events/pkg/storage"
)

type HealthService struct {
	store  storage.EventStore
	logger *zap.SugaredLogger
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

func NewHealthService(store storage.EventStore, logger *zap.SugaredLogger) *HealthService {
	return &HealthService{
		store:  store,
		logger: logger,
	}
}

func (h *HealthService) CheckHealth(ctx context.Context) *HealthStatus {
	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
	}

	// Check database by attempting a simple query
	if err := h.checkDatabase(ctx); err != nil {
		status.Services["database"] = "unhealthy"
		status.Status = "unhealthy"
		h.logger.Errorw("Database health check failed", "error", err)
	} else {
		status.Services["database"] = "healthy"
	}

	// Check Redis (simplified - in real implementation, ping Redis)
	status.Services["redis"] = "healthy" // TODO: implement actual Redis check

	// Check other services
	status.Services["processing"] = "healthy"

	return status
}

func (h *HealthService) checkDatabase(ctx context.Context) error {
	// Try to get an event by a non-existent ID to test DB connectivity
	_, err := h.store.GetEventByID(ctx, "health-check-id")
	// We expect this to return an error (event not found), but not a connection error
	if err != nil && err.Error() != "sql: no rows in result set" {
		return err
	}
	return nil
}