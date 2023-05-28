package heartbeat

import (
	"time"
)

type Heartbeat struct {
	ID        string    `json:"id,omitempty"`
	Status    string    `json:"status,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
