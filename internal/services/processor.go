package services

import (
	"context"
	"realtime-events/internal/models"
	"realtime-events/pkg/storage"

	"go.uber.org/zap"
)

type EventProcessor struct {
	store  storage.EventStore
	logger *zap.SugaredLogger
}

func NewEventProcessor(store storage.EventStore, logger *zap.SugaredLogger) *EventProcessor {
	return &EventProcessor{
		store:  store,
		logger: logger,
	}
}

func (p *EventProcessor) ProcessEvent(ctx context.Context, event *models.Event) error {
	// Normalize event
	if err := p.normalizeEvent(event); err != nil {
		p.logger.Errorw("Failed to normalize event", "error", err, "event_id", event.ID)
		return err
	}

	// Update aggregates (simplified - in real impl, use time-series DB)
	if err := p.updateAggregates(ctx, event); err != nil {
		p.logger.Errorw("Failed to update aggregates", "error", err, "event_id", event.ID)
	}

	// Evaluate rules and trigger webhooks
	if err := p.evaluateRules(ctx, event); err != nil {
		p.logger.Errorw("Failed to evaluate rules", "error", err, "event_id", event.ID)
	}

	p.logger.Infow("Event processed successfully", "event_id", event.ID, "event_name", event.EventName)
	return nil
}

func (p *EventProcessor) normalizeEvent(event *models.Event) error {
	// Add normalization logic here
	// e.g., standardize field names, validate data types, etc.
	return nil
}

func (p *EventProcessor) updateAggregates(ctx context.Context, event *models.Event) error {
	// In real implementation, update time-series aggregates
	// For now, just log
	p.logger.Infow("Updating aggregates", "event_name", event.EventName, "user_id", event.UserID)
	return nil
}

func (p *EventProcessor) evaluateRules(ctx context.Context, event *models.Event) error {
	// In real implementation, fetch rules from DB and evaluate
	// If conditions match, trigger webhooks
	p.logger.Infow("Evaluating rules", "event_name", event.EventName)
	return nil
}