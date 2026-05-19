package filter

import (
	"testing"

	"go.uber.org/zap"
)

// mockBlocker simulates blocklist.Manager
type mockBlocker struct {
	blocked map[string]string // domain -> category
}

func (m *mockBlocker) IsBlocked(domain string) (bool, string) {
	cat, ok := m.blocked[domain]
	return ok, cat
}

func newTestEngine(blocked map[string]string, allowlist []string) *Engine {
	log, _ := zap.NewNop(), struct{}{}
	_ = log
	logger := zap.NewNop()
	cfg := FilterConfig{
		BlockPage: "0.0.0.0",
		Allowlist: allowlist,
	}
	bl := &mockBlocker{blocked: blocked}
	return NewEngine(bl, cfg, logger)
}

func TestBlockDecision(t *testing.T) {
	engine := newTestEngine(map[string]string{
		"evil.com": "malware",
	}, nil)

	d := engine.Evaluate("evil.com")
	if d.Action != "block" {
		t.Errorf("expected block, got %s", d.Action)
	}
	if d.Category != "malware" {
		t.Errorf("expected category malware, got %s", d.Category)
	}
}

func TestAllowDecision(t *testing.T) {
	engine := newTestEngine(map[string]string{}, nil)

	d := engine.Evaluate("canada.ca")
	if d.Action != "allow" {
		t.Errorf("expected allow, got %s", d.Action)
	}
}

func TestAllowlistOverridesBlock(t *testing.T) {
	engine := newTestEngine(
		map[string]string{"trusted.ca": "malware"}, // in blocklist
		[]string{"trusted.ca"},                      // but also in allowlist
	)

	d := engine.Evaluate("trusted.ca")
	if d.Action != "allow" {
		t.Errorf("allowlist should override blocklist, got %s", d.Action)
	}
}

func TestSubdomainBlock(t *testing.T) {
	// sub.evil.com should be blocked if evil.com is in blocklist
	// (handled by blocklist.Manager — here we just test filter pass-through)
	engine := newTestEngine(map[string]string{
		"sub.evil.com": "phishing",
	}, nil)

	d := engine.Evaluate("sub.evil.com.")
	if d.Action != "block" {
		t.Errorf("subdomain should be blocked, got %s", d.Action)
	}
}
