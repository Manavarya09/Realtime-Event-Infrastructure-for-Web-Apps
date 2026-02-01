package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"realtime-events/internal/models"
)

type EventStore interface {
	InsertEvent(ctx context.Context, event *models.Event) error
	GetEventByID(ctx context.Context, id string) (*models.Event, error)
	GetAPIKeyByHash(ctx context.Context, hash string) (*models.APIKey, error)
}

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgres(url string) (*PostgresStore, error) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

func (s *PostgresStore) InsertEvent(ctx context.Context, event *models.Event) error {
	query := `
		INSERT INTO events (id, project_id, event_name, user_id, timestamp, metadata, received_at, ip_address, user_agent, idempotency_key)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := s.pool.Exec(ctx, query,
		event.ID, event.ProjectID, event.EventName, event.UserID,
		event.Timestamp, event.Metadata, event.ReceivedAt,
		event.IPAddress, event.UserAgent, event.IdempotencyKey)
	return err
}

func (s *PostgresStore) GetEventByID(ctx context.Context, id string) (*models.Event, error) {
	query := `SELECT id, project_id, event_name, user_id, timestamp, metadata, received_at, ip_address, user_agent FROM events WHERE id = $1`
	row := s.pool.QueryRow(ctx, query, id)
	var event models.Event
	err := row.Scan(&event.ID, &event.ProjectID, &event.EventName, &event.UserID,
		&event.Timestamp, &event.Metadata, &event.ReceivedAt, &event.IPAddress, &event.UserAgent)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *PostgresStore) GetAPIKeyByHash(ctx context.Context, hash string) (*models.APIKey, error) {
	query := `SELECT id, project_id, key_hash, name, expires_at FROM api_keys WHERE key_hash = $1`
	row := s.pool.QueryRow(ctx, query, hash)
	var key models.APIKey
	err := row.Scan(&key.ID, &key.ProjectID, &key.KeyHash, &key.Name, &key.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &key, nil
}