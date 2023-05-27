package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/daniel-oliveiravas/heartbeat-service/business/event"
	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/kafkapubsub"
	"google.golang.org/protobuf/proto"
)

type PublisherConfig struct {
	TopicURL       string
	KafkaAddresses []string
}

const messageKeyName = "key"

type Publisher struct {
	cfg   PublisherConfig
	topic *pubsub.Topic
}

func NewPublisher(cfg PublisherConfig) (*Publisher, error) {
	config := kafkapubsub.MinimalConfig()
	kafkaTopic, err := kafkapubsub.OpenTopic(cfg.KafkaAddresses, config, cfg.TopicURL, &kafkapubsub.TopicOptions{
		KeyName: messageKeyName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open topic with kafkapubsub. %w", err)
	}

	return &Publisher{
		cfg:   cfg,
		topic: kafkaTopic,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, heartbeat heartbeat.Heartbeat) error {
	heartbeatStatus, err := toHeartbeatStatusEvent(heartbeat.Status)
	if err != nil {
		return err
	}

	heartbeatEvent := event.Heartbeat{
		Id:        heartbeat.ID,
		Status:    heartbeatStatus,
		Timestamp: heartbeat.Timestamp.Format(time.RFC3339),
	}

	b, err := proto.Marshal(&heartbeatEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal event.Heartbeat. %w", err)
	}

	msg := pubsub.Message{
		Body: b,
		Metadata: map[string]string{
			messageKeyName: heartbeat.ID,
		},
	}

	err = p.topic.Send(ctx, &msg)
	if err != nil {
		return fmt.Errorf("failed to send message to topic. %w", err)
	}

	return nil
}

func toHeartbeatStatusEvent(status string) (event.HeartbeatStatus, error) {
	switch status {
	case "online":
		return event.HeartbeatStatus_ONLINE, nil
	case "offline":
		return event.HeartbeatStatus_OFFLINE, nil
	default:
		return 0, fmt.Errorf("invalid status %s", status)
	}
}
