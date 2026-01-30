-- Database Schema for Realtime Event Infrastructure
-- Uses PostgreSQL with TimescaleDB extension for time-series optimization

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Organizations table
CREATE TABLE organizations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Projects table
CREATE TABLE projects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id),
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- API Keys table
CREATE TABLE api_keys (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  project_id UUID NOT NULL REFERENCES projects(id),
  key_hash TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ
);

-- Events table (hypertable for time-series)
CREATE TABLE events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  project_id UUID NOT NULL REFERENCES projects(id),
  event_name TEXT NOT NULL,
  user_id TEXT,
  timestamp TIMESTAMPTZ NOT NULL,
  metadata JSONB,
  received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  ip_address INET,
  user_agent TEXT,
  idempotency_key TEXT
);

-- Convert to hypertable
SELECT create_hypertable('events', 'timestamp', if_not_exists => TRUE);

-- Indexes for events
CREATE INDEX idx_events_project_timestamp ON events (project_id, timestamp DESC);
CREATE INDEX idx_events_name ON events (event_name);
CREATE INDEX idx_events_user ON events (user_id);
CREATE INDEX idx_events_idempotency ON events (idempotency_key);

-- Aggregates table (materialized view for analytics)
CREATE TABLE aggregates (
  project_id UUID NOT NULL REFERENCES projects(id),
  event_name TEXT NOT NULL,
  time_bucket TIMESTAMPTZ NOT NULL,
  count BIGINT NOT NULL DEFAULT 0,
  unique_users BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (project_id, event_name, time_bucket)
);

-- Convert to hypertable
SELECT create_hypertable('aggregates', 'time_bucket', if_not_exists => TRUE);

-- Rules table
CREATE TABLE rules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  project_id UUID NOT NULL REFERENCES projects(id),
  name TEXT NOT NULL,
  conditions JSONB NOT NULL,
  actions JSONB NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Webhooks table
CREATE TABLE webhooks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  project_id UUID NOT NULL REFERENCES projects(id),
  url TEXT NOT NULL,
  secret TEXT NOT NULL,
  events TEXT[] NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Webhook attempts table
CREATE TABLE webhook_attempts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  webhook_id UUID NOT NULL REFERENCES webhooks(id),
  event_id UUID NOT NULL REFERENCES events(id),
  status TEXT NOT NULL,
  response_code INTEGER,
  response_body TEXT,
  attempt_number INTEGER NOT NULL DEFAULT 1,
  next_retry_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Dead letter queue for failed events
CREATE TABLE dead_letter_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  original_event_id UUID REFERENCES events(id),
  project_id UUID NOT NULL REFERENCES projects(id),
  payload JSONB NOT NULL,
  error_message TEXT NOT NULL,
  failed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Users table for dashboard auth
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id),
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL CHECK (role IN ('admin', 'developer', 'viewer')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);