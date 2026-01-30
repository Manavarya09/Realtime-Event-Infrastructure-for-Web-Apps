# API Contract

## Authentication
All requests require API key in Authorization header:
`Authorization: Bearer <api_key>`

## Event Ingestion

### Single Event
```http
POST /api/v1/events
Content-Type: application/json

{
  "event_name": "user_signed_up",
  "user_id": "uuid",
  "timestamp": "2024-01-30T10:00:00Z",
  "metadata": {
    "plan": "pro",
    "source": "landing_page"
  },
  "idempotency_key": "optional-unique-key"
}
```

**Response:**
```json
{
  "event_id": "uuid",
  "status": "accepted"
}
```

**Validation:**
- event_name: required, string, max 100 chars
- user_id: optional, string
- timestamp: optional, ISO8601, defaults to now
- metadata: optional, object, max 10KB
- idempotency_key: optional, prevents duplicate processing

### Batch Events
```http
POST /api/v1/events/batch
Content-Type: application/json

[
  {
    "event_name": "page_view",
    "user_id": "user123",
    "metadata": {"page": "/home"}
  }
]
```

**Limits:** Max 100 events per batch, 1MB total payload.

## Analytics API

### Get Event Counts
```http
GET /api/v1/analytics/events?project_id=uuid&event_name=user_signed_up&period=1h
```

**Response:**
```json
{
  "count": 1250,
  "unique_users": 890,
  "time_buckets": [
    {"timestamp": "2024-01-30T09:00:00Z", "count": 100}
  ]
}
```

## Webhooks

### Create Webhook
```http
POST /api/v1/webhooks
Content-Type: application/json

{
  "url": "https://api.example.com/webhook",
  "events": ["user_signed_up", "purchase_completed"],
  "secret": "webhook-secret"
}
```

### List Webhooks
```http
GET /api/v1/webhooks?project_id=uuid
```

## Rules

### Create Rule
```http
POST /api/v1/rules
Content-Type: application/json

{
  "name": "High Value Purchase",
  "conditions": {
    "event_name": "purchase_completed",
    "metadata.amount": { "$gt": 500 }
  },
  "actions": [
    {
      "type": "webhook",
      "webhook_id": "uuid"
    }
  ]
}
```

## Error Responses
```json
{
  "error": "invalid_request",
  "message": "Event name is required",
  "code": 400
}
```

## Rate Limits
- 1000 events/minute per API key
- 10 requests/minute for analytics API

## Webhook Payload
```json
{
  "event_id": "uuid",
  "event_name": "purchase_completed",
  "user_id": "user123",
  "timestamp": "2024-01-30T10:00:00Z",
  "metadata": {...},
  "signature": "sha256=..."
}
```