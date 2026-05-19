package metrics

import (
	"testing"
	"time"
)

func TestRecordAndSnapshot(t *testing.T) {
	m := New()

	m.RecordQuery("block", "malware")
	m.RecordQuery("block", "malware")
	m.RecordQuery("block", "phishing")
	m.RecordQuery("allow", "")

	snap := m.Snapshot()

	if snap.Total != 4 {
		t.Errorf("expected total 4, got %d", snap.Total)
	}
	if snap.Blocked != 3 {
		t.Errorf("expected blocked 3, got %d", snap.Blocked)
	}
	if snap.Allowed != 1 {
		t.Errorf("expected allowed 1, got %d", snap.Allowed)
	}
	if snap.Categories["malware"] != 2 {
		t.Errorf("expected malware count 2, got %d", snap.Categories["malware"])
	}
	if snap.Categories["phishing"] != 1 {
		t.Errorf("expected phishing count 1, got %d", snap.Categories["phishing"])
	}
}

func TestLatencyBuckets(t *testing.T) {
	m := New()

	m.RecordLatency(500 * time.Microsecond) // <1ms
	m.RecordLatency(2 * time.Millisecond)   // <5ms
	m.RecordLatency(7 * time.Millisecond)   // <10ms
	m.RecordLatency(20 * time.Millisecond)  // <50ms
	m.RecordLatency(100 * time.Millisecond) // 50ms+

	snap := m.Snapshot()
	if snap.Latency.Under1ms != 1 {
		t.Errorf("expected Under1ms=1, got %d", snap.Latency.Under1ms)
	}
	if snap.Latency.Over50ms != 1 {
		t.Errorf("expected Over50ms=1, got %d", snap.Latency.Over50ms)
	}
}
