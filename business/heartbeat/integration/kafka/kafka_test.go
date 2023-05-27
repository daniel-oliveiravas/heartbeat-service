package kafka_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat"
	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat/integration/kafka"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPublisher_Publish(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests: set INTEGRATION environment variable")
	}

	publisherCfg := kafka.PublisherConfig{
		TopicURL:       "heartbeat-events-test",
		KafkaAddresses: []string{"localhost:9092"},
	}
	publisher, err := kafka.NewPublisher(publisherCfg)
	require.NoError(t, err)

	ctx := context.Background()
	beat := heartbeat.Heartbeat{
		ID:        uuid.New().String(),
		Status:    "online",
		Timestamp: time.Now().UTC(),
	}
	err = publisher.Publish(ctx, beat)
	require.NoError(t, err)
}
