package redis_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat"
	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat/integration/redis"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}

	ctx := context.Background()
	redisClient := goredis.NewClient(&goredis.Options{
		Addr: "localhost:6379",
	})

	repo := redis.NewRepository(redisClient, time.Second*30)

	beat := heartbeat.Heartbeat{
		ID:        uuid.New().String(),
		Status:    "online",
		Timestamp: time.Now().UTC(),
	}

	err := repo.Upsert(ctx, beat)
	require.NoError(t, err)

	storedBeat, err := repo.Get(ctx, beat.ID)
	require.NoError(t, err)
	assert.Equal(t, beat, storedBeat)
}
