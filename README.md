# Realtime Event Infrastructure for Web Apps

A realtime event infrastructure platform similar to Segment, designed for modern SaaS and web applications. Built with Go backend, React frontend, and enterprise-grade observability.

## 🏗️ Architecture Overview

The system consists of microservices:

- **Event Ingestion Service**: High-throughput HTTP API for event collection
- **Event Processing Engine**: Concurrent pipeline for event processing, rules evaluation, and webhook triggering
- **Webhook Delivery System**: Reliable webhook engine with retries and dead-letter queues
- **Realtime Analytics Engine**: Live metrics and time-series aggregation
- **Storage Layer**: PostgreSQL + TimescaleDB for optimized event storage
- **Frontend Dashboard**: Developer-first SaaS dashboard with real-time updates

### Data Flow
1. Client sends events to Ingestion API
2. Events validated, enriched, and queued (Redis Streams)
3. Processing service consumes, stores, aggregates, evaluates rules
4. Rules trigger webhooks via Webhook service
5. Analytics service computes live metrics
6. Dashboard displays real-time data via WebSocket

### Scaling Strategy
- Horizontal scaling via Kubernetes
- Database sharding for events
- Redis clustering for queues
- CDN for frontend assets

### Failure Scenarios & Resilience
- Queue backpressure: Rate limiting and buffering
- Service failures: At-least-once delivery with retries
- Dead-letter queues for unprocessable events
- Circuit breakers for downstream services

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Node.js 18+

### Setup
```bash
# Clone and setup
git clone <repo>
cd realtime-event-infrastructure

# Start services
docker-compose up -d

# Run backend
go run cmd/ingestion/main.go

# Run frontend
cd web && npm install && npm start
```

## 📊 API Contract

### Event Ingestion
```http
POST /api/v1/events
Authorization: Bearer <api_key>
Content-Type: application/json

{
  "event_name": "user_signed_up",
  "user_id": "uuid",
  "timestamp": "2024-01-30T10:00:00Z",
  "metadata": {
    "plan": "pro",
    "source": "landing_page"
  }
}
```

Response: 202 Accepted with event ID

### Batch Events
```http
POST /api/v1/events/batch
Authorization: Bearer <api_key>

[{"event_name": "..."}, ...]
```

## 🗄️ Database Schema

### Events Table (TimescaleDB)
```sql
CREATE TABLE events (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  event_name TEXT NOT NULL,
  user_id TEXT,
  timestamp TIMESTAMPTZ NOT NULL,
  metadata JSONB,
  received_at TIMESTAMPTZ NOT NULL,
  ip_address INET,
  user_agent TEXT
);
```

### Aggregates Table
```sql
CREATE TABLE aggregates (
  project_id UUID,
  event_name TEXT,
  time_bucket TIMESTAMPTZ,
  count BIGINT,
  PRIMARY KEY (project_id, event_name, time_bucket)
);
```

## 🔐 Authentication

- API Key per project
- JWT for dashboard access
- RBAC: admin, developer, viewer

## 📈 Observability

- Metrics: Prometheus (/metrics)
- Tracing: Jaeger
- Logging: Structured JSON logs
- Health: /health endpoint

## 🧪 Testing

```bash
# Unit tests
go test ./...

# Integration tests
docker-compose -f docker-compose.test.yml up
```

## 📚 Documentation

- [API Docs](./docs/api.md)
- [Architecture](./docs/architecture.md)
- [Deployment](./docs/deployment.md)

## 🤝 Contributing

1. Fork the repo
2. Create feature branch
3. Add tests
4. Submit PR

## 📄 License

MIT License