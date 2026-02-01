package services

import (
	"fmt"
	"regexp"
	"strings"

	"realtime-events/internal/models"
)

type ValidationService struct{}

func NewValidationService() *ValidationService {
	return &ValidationService{}
}

func (v *ValidationService) ValidateEventRequest(req *models.EventRequest) error {
	// Validate event name format
	if !v.isValidEventName(req.EventName) {
		return fmt.Errorf("event_name must be alphanumeric with underscores, starting with a letter")
	}

	// Validate user ID if provided
	if req.UserID != nil && !v.isValidUserID(*req.UserID) {
		return fmt.Errorf("user_id must be alphanumeric with underscores and hyphens")
	}

	// Validate idempotency key if provided
	if req.IdempotencyKey != nil && !v.isValidIdempotencyKey(*req.IdempotencyKey) {
		return fmt.Errorf("idempotency_key must be alphanumeric")
	}

	// Validate metadata
	if err := v.validateMetadata(req.Metadata); err != nil {
		return err
	}

	return nil
}

func (v *ValidationService) isValidEventName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, name)
	return matched
}

func (v *ValidationService) isValidUserID(userID string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, userID)
	return matched
}

func (v *ValidationService) isValidIdempotencyKey(key string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, key)
	return matched
}

func (v *ValidationService) validateMetadata(metadata map[string]interface{}) error {
	if metadata == nil {
		return nil
	}

	for key, value := range metadata {
		if len(key) > 100 {
			return fmt.Errorf("metadata key too long: %s", key)
		}

		// Convert value to string for length check
		valueStr := fmt.Sprintf("%v", value)
		if len(valueStr) > 1000 {
			return fmt.Errorf("metadata value too long for key: %s", key)
		}

		// Check for nested objects (simplified)
		if strings.Contains(valueStr, "map[") || strings.Contains(valueStr, "{") {
			return fmt.Errorf("nested objects not allowed in metadata for key: %s", key)
		}
	}

	return nil
}