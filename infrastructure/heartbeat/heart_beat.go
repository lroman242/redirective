package heartbeat

import (
	"context"
	"time"
)

// HeartBeat is a service which check
// if parts of application is working fine
type HeartBeat interface {
	Monitor(ctx context.Context, time time.Duration)
	Check() error
}
