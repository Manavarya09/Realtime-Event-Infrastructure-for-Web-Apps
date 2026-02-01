package queue

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"realtime-events/internal/models"
)

type EventQueue interface {
	PublishEvent(ctx context.Context, event *models.Event) error
	ConsumeEvents(ctx context.Context, handler func(*models.Event) error) error
	Close() error
}

type RedisQueue struct {
	client *redis.Client
	stream string
}

func NewRedisQueue(url, stream string) (*RedisQueue, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return &RedisQueue{client: client, stream: stream}, nil
}

func (q *RedisQueue) PublishEvent(ctx context.Context, event *models.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: q.stream,
		Values: map[string]interface{}{"event": data},
	}).Err()
}

func (q *RedisQueue) ConsumeEvents(ctx context.Context, handler func(*models.Event) error) error {
	lastID := "0"
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			streams, err := q.client.XRead(ctx, &redis.XReadArgs{
				Streams: []string{q.stream, lastID},
				Count:   10,
				Block:   0,
			}).Result()
			if err != nil {
				return err
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					var event models.Event
					if err := json.Unmarshal([]byte(message.Values["event"].(string)), &event); err != nil {
						continue // Skip invalid messages
					}
					if err := handler(&event); err != nil {
						// Handle processing error (could send to dead letter)
						continue
					}
					lastID = message.ID
				}
			}
		}
	}
}

func (q *RedisQueue) Close() error {
	return q.client.Close()
}