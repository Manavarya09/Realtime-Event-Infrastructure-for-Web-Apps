package services

import (
	"testing"

	"realtime-events/internal/models"
)

func TestValidationService_ValidateEventRequest(t *testing.T) {
	vs := NewValidationService()

	tests := []struct {
		name    string
		req     models.EventRequest
		wantErr bool
	}{
		{
			name: "valid event",
			req: models.EventRequest{
				EventName: "user_signup",
				UserID:    stringPtr("user123"),
			},
			wantErr: false,
		},
		{
			name: "invalid event name - starts with number",
			req: models.EventRequest{
				EventName: "123invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid event name - special characters",
			req: models.EventRequest{
				EventName: "user-signup!",
			},
			wantErr: true,
		},
		{
			name: "invalid user ID",
			req: models.EventRequest{
				EventName: "user_signup",
				UserID:    stringPtr("user@123"),
			},
			wantErr: true,
		},
		{
			name: "valid metadata",
			req: models.EventRequest{
				EventName: "user_signup",
				Metadata: map[string]interface{}{
					"plan":   "premium",
					"source": "web",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid metadata - nested object",
			req: models.EventRequest{
				EventName: "user_signup",
				Metadata: map[string]interface{}{
					"nested": map[string]interface{}{"key": "value"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vs.ValidateEventRequest(&tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEventRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}