package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics holds lightweight in-memory counters.
// Zero PII — no per-query logs, no IP storage, no domain logging.
type Metrics struct {
	TotalQueries   atomic.Int64
	BlockedQueries atomic.Int64
	AllowedQueries atomic.Int64

	categoryMu     sync.RWMutex
	categoryCounts map[string]*atomic.Int64

	// Latency buckets: <1ms | <5ms | <10ms | <50ms | 50ms+
	LatencyBuckets [5]atomic.Int64
}

func New() *Metrics {
	return &Metrics{
		categoryCounts: make(map[string]*atomic.Int64),
	}
}

func (m *Metrics) RecordQuery(action, category string) {
	m.TotalQueries.Add(1)
	if action == "block" {
		m.BlockedQueries.Add(1)
		m.categoryMu.Lock()
		c, ok := m.categoryCounts[category]
		if !ok {
			c = &atomic.Int64{}
			m.categoryCounts[category] = c
		}
		m.categoryMu.Unlock()
		c.Add(1)
	} else {
		m.AllowedQueries.Add(1)
	}
}

func (m *Metrics) RecordLatency(d time.Duration) {
	ms := d.Milliseconds()
	switch {
	case ms < 1:
		m.LatencyBuckets[0].Add(1)
	case ms < 5:
		m.LatencyBuckets[1].Add(1)
	case ms < 10:
		m.LatencyBuckets[2].Add(1)
	case ms < 50:
		m.LatencyBuckets[3].Add(1)
	default:
		m.LatencyBuckets[4].Add(1)
	}
}

type Snapshot struct {
	Total      int64            `json:"total_queries"`
	Blocked    int64            `json:"blocked_queries"`
	Allowed    int64            `json:"allowed_queries"`
	Categories map[string]int64 `json:"blocked_by_category"`
	Latency    LatencySnapshot  `json:"latency_ms"`
}

type LatencySnapshot struct {
	Under1ms  int64 `json:"under_1ms"`
	Under5ms  int64 `json:"under_5ms"`
	Under10ms int64 `json:"under_10ms"`
	Under50ms int64 `json:"under_50ms"`
	Over50ms  int64 `json:"over_50ms"`
}

func (m *Metrics) Snapshot() Snapshot {
	m.categoryMu.RLock()
	cats := make(map[string]int64, len(m.categoryCounts))
	for k, v := range m.categoryCounts {
		cats[k] = v.Load()
	}
	m.categoryMu.RUnlock()

	return Snapshot{
		Total:      m.TotalQueries.Load(),
		Blocked:    m.BlockedQueries.Load(),
		Allowed:    m.AllowedQueries.Load(),
		Categories: cats,
		Latency: LatencySnapshot{
			Under1ms:  m.LatencyBuckets[0].Load(),
			Under5ms:  m.LatencyBuckets[1].Load(),
			Under10ms: m.LatencyBuckets[2].Load(),
			Under50ms: m.LatencyBuckets[3].Load(),
			Over50ms:  m.LatencyBuckets[4].Load(),
		},
	}
}
