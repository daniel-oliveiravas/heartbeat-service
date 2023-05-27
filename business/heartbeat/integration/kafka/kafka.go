package kafka

import (
	"context"

	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat"
)

type Publisher struct {
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Publish(ctx context.Context, heartbeat heartbeat.Heartbeat) error {
	//TODO: Implement publishing event to kafka
	return nil
}
