package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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
	// Normalize event name to lowercase
	event.EventName = strings.ToLower(strings.TrimSpace(event.EventName))

	// Validate event name format
	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`).MatchString(event.EventName) {
		return fmt.Errorf("invalid event name format")
	}

	// Ensure metadata is not nil
	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}

	return nil
}

func (p *EventProcessor) updateAggregates(ctx context.Context, event *models.Event) error {
	// In real implementation, update time-series aggregates
	// For now, just increment counters
	p.logger.Infow("Updating aggregates", "event_name", event.EventName, "user_id", event.UserID)

	// TODO: Implement actual aggregation logic with Redis/TimescaleDB
	return nil
}

func (p *EventProcessor) evaluateRules(ctx context.Context, event *models.Event) error {
	// Simple rule evaluation - in real implementation, fetch rules from DB
	rules := []map[string]interface{}{
		{
			"event_name": "user_signup",
			"conditions": map[string]interface{}{
				"metadata.plan": "premium",
			},
			"actions": []map[string]interface{}{
				{
					"type": "webhook",
					"url":  "https://api.example.com/webhooks/premium-signup",
				},
			},
		},
	}

	for _, rule := range rules {
		if p.matchesRule(event, rule) {
			if err := p.executeActions(ctx, event, rule["actions"].([]map[string]interface{})); err != nil {
				p.logger.Errorw("Failed to execute rule actions", "error", err, "rule", rule)
			}
		}
	}

	return nil
}

func (p *EventProcessor) matchesRule(event *models.Event, rule map[string]interface{}) bool {
	// Check event name match
	if ruleEventName, ok := rule["event_name"].(string); ok {
		if event.EventName != ruleEventName {
			return false
		}
	}

	// Check conditions
	if conditions, ok := rule["conditions"].(map[string]interface{}); ok {
		for key, expectedValue := range conditions {
			actualValue := p.getNestedValue(event, key)
			if actualValue != expectedValue {
				return false
			}
		}
	}

	return true
}

func (p *EventProcessor) getNestedValue(event *models.Event, key string) interface{} {
	parts := strings.Split(key, ".")
	if len(parts) == 1 {
		switch parts[0] {
		case "event_name":
			return event.EventName
		case "user_id":
			return event.UserID
		case "project_id":
			return event.ProjectID
		}
	} else if parts[0] == "metadata" && len(parts) > 1 {
		if event.Metadata != nil {
			return event.Metadata[parts[1]]
		}
	}
	return nil
}

func (p *EventProcessor) executeActions(ctx context.Context, event *models.Event, actions []map[string]interface{}) error {
	for _, action := range actions {
		actionType, ok := action["type"].(string)
		if !ok {
			continue
		}

		switch actionType {
		case "webhook":
			if err := p.sendWebhook(ctx, event, action); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *EventProcessor) sendWebhook(ctx context.Context, event *models.Event, action map[string]interface{}) error {
	url, ok := action["url"].(string)
	if !ok {
		return fmt.Errorf("webhook URL not specified")
	}

	// In real implementation, use HTTP client with retries, circuit breaker, etc.
	p.logger.Infow("Sending webhook", "url", url, "event_id", event.ID)

	// TODO: Implement actual HTTP request with proper error handling and retries
	return nil
}
