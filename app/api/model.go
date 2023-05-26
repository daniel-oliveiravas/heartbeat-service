package api

import (
	"time"

	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat"
)

type HeartbeatSignal struct {
	Status    string
	Timestamp time.Time
}

func (h HeartbeatSignal) toUsecase(identifier string) heartbeat.Heartbeat {
	return heartbeat.Heartbeat{
		ID:        identifier,
		Status:    h.Status,
		Timestamp: h.Timestamp,
	}
}
