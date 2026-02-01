package models

import (
	"time"
)

type Event struct {
	ID             string                 `json:"id" db:"id"`
	ProjectID      string                 `json:"project_id" db:"project_id"`
	EventName      string                 `json:"event_name" db:"event_name" validate:"required,max=100"`
	UserID         *string                `json:"user_id,omitempty" db:"user_id"`
	Timestamp      time.Time              `json:"timestamp" db:"timestamp"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	ReceivedAt     time.Time              `json:"received_at" db:"received_at"`
	IPAddress      *string                `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent      *string                `json:"user_agent,omitempty" db:"user_agent"`
	IdempotencyKey *string                `json:"idempotency_key,omitempty" db:"idempotency_key"`
}

type EventRequest struct {
	EventName      string                 `json:"event_name" binding:"required,min=1,max=100" validate:"required,min=1,max=100"`
	UserID         *string                `json:"user_id,omitempty" validate:"omitempty,max=100"`
	Timestamp      *time.Time             `json:"timestamp,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" validate:"omitempty,dive,keys,keymax=100,endkeys,valuemax=1000"`
	IdempotencyKey *string                `json:"idempotency_key,omitempty" validate:"omitempty,max=255"`
}

type BatchEventRequest struct {
	Events []EventRequest `json:"events" binding:"required,min=1,max=100" validate:"required,min=1,max=100,dive"`
}

type APIKey struct {
	ID        string     `db:"id"`
	ProjectID string     `db:"project_id"`
	KeyHash   string     `db:"key_hash"`
	Name      string     `db:"name"`
	ExpiresAt *time.Time `db:"expires_at"`
}
