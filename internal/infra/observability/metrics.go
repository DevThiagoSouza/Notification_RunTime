package observability

import "sync/atomic"

type Metrics struct {
	created   uint64
	processed uint64
	sent      uint64
	failed    uint64
	retried   uint64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) IncCreated() {
	atomic.AddUint64(&m.created, 1)
}

func (m *Metrics) IncProcessed() {
	atomic.AddUint64(&m.processed, 1)
}

func (m *Metrics) IncSent() {
	atomic.AddUint64(&m.sent, 1)
}

func (m *Metrics) IncFailed() {
	atomic.AddUint64(&m.failed, 1)
}

func (m *Metrics) IncRetried() {
	atomic.AddUint64(&m.retried, 1)
}

func (m *Metrics) Snapshot() map[string]uint64 {
	return map[string]uint64{
		"created":   atomic.LoadUint64(&m.created),
		"processed": atomic.LoadUint64(&m.processed),
		"sent":      atomic.LoadUint64(&m.sent),
		"failed":    atomic.LoadUint64(&m.failed),
		"retried":   atomic.LoadUint64(&m.retried),
	}
}
