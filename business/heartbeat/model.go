package heartbeat

import (
	"time"
)

type Heartbeat struct {
	ID        string    `json:"ID,omitempty"`
	Status    string    `json:"status,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
