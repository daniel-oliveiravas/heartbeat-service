package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client          *redis.Client
	heartbeatExpiry time.Duration
}

func NewRepository(client *redis.Client, heartbeatExpiry time.Duration) *Repository {
	return &Repository{
		client:          client,
		heartbeatExpiry: heartbeatExpiry,
	}
}

type HeartbeatDB struct {
	ID        string    `redis:"id"`
	Status    string    `redis:"status"`
	Timestamp time.Time `redis:"timestamp"`
}

func (r *Repository) Upsert(ctx context.Context, beat heartbeat.Heartbeat) error {
	dbModel := HeartbeatDB{
		ID:        beat.ID,
		Status:    beat.Status,
		Timestamp: beat.Timestamp,
	}
	err := r.client.HSet(ctx, beat.ID, dbModel).Err()
	if err != nil {
		return fmt.Errorf("failed to set heartbeat to redis. %w", err)
	}

	err = r.client.Expire(ctx, beat.ID, r.heartbeatExpiry).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiry to key %s. %w", beat.ID, err)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id string) (heartbeat.Heartbeat, error) {
	var dbModel HeartbeatDB
	err := r.client.HGetAll(ctx, id).Scan(&dbModel)
	if err != nil {
		return heartbeat.Heartbeat{}, fmt.Errorf("failed to get heartbeat from redis. %w", err)
	}

	return heartbeat.Heartbeat{
		ID:        dbModel.ID,
		Status:    dbModel.Status,
		Timestamp: dbModel.Timestamp,
	}, nil
}
