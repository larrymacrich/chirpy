package metrics

import "sync/atomic"

type MetricsConfig struct {
	FileserverRequests atomic.Int32
}
