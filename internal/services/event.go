package services

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"realtime-events/internal/models"
	"realtime-events/pkg/queue"
	"realtime-events/pkg/storage"
)

type EventService struct {
	store  storage.EventStore
	queue  queue.EventQueue
	logger *zap.SugaredLogger
}

func NewEventService(store storage.EventStore, queue queue.EventQueue, logger *zap.SugaredLogger) *EventService {
	return &EventService{
		store:  store,
		queue:  queue,
		logger: logger,
	}
}

func (s *EventService) ProcessEvent(ctx context.Context, req *models.EventRequest, projectID string, ip net.IP, userAgent string) (*models.Event, error) {
	// Create event
	ipStr := ip.String()
	event := &models.Event{
		ID:             uuid.New().String(),
		ProjectID:      projectID,
		EventName:      req.EventName,
		UserID:         req.UserID,
		Metadata:       req.Metadata,
		ReceivedAt:     time.Now(),
		IPAddress:      &ipStr,
		UserAgent:      &userAgent,
		IdempotencyKey: req.IdempotencyKey,
	}

	// Set timestamp
	if req.Timestamp != nil {
		event.Timestamp = *req.Timestamp
	} else {
		event.Timestamp = time.Now()
	}

	// Validate and enrich
	if err := s.validateEvent(event); err != nil {
		return nil, err
	}

	// Store event
	if err := s.store.InsertEvent(ctx, event); err != nil {
		s.logger.Errorw("Failed to store event", "error", err, "event_id", event.ID)
		return nil, err
	}

	// Queue for processing
	if err := s.queue.PublishEvent(ctx, event); err != nil {
		s.logger.Errorw("Failed to queue event", "error", err, "event_id", event.ID)
		// Don't return error, event is already stored
	}

	s.logger.Infow("Event processed", "event_id", event.ID, "event_name", event.EventName)
	return event, nil
}

func (s *EventService) validateEvent(event *models.Event) error {
	if event.EventName == "" {
		return fmt.Errorf("event_name is required")
	}
	if len(event.EventName) > 100 {
		return fmt.Errorf("event_name too long")
	}
	// Add more validation as needed
	return nil
}

func (s *EventService) AuthenticateAPIKey(ctx context.Context, apiKey string) (string, error) {
	hash := sha256.Sum256([]byte(apiKey))
	hashStr := fmt.Sprintf("%x", hash)

	key, err := s.store.GetAPIKeyByHash(ctx, hashStr)
	if err != nil {
		return "", fmt.Errorf("invalid API key")
	}

	if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("API key expired")
	}

	return key.ProjectID, nil
}