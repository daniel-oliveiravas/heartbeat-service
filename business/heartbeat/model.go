package heartbeat

import (
	"time"
)

type Heartbeat struct {
	ID        string
	Status    string
	Timestamp time.Time
}
