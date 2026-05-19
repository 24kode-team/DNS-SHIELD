package blocklist

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dns-shield/shield/internal/config"
	"go.uber.org/zap"
)

// Manager holds all blocklists per category and refreshes them periodically.
// No query-level logging — only aggregate counts stored in memory.
type Manager struct {
	mu      sync.RWMutex
	sets    map[string]map[string]struct{} // category -> domain set
	feeds   []config.FeedConfig
	dataDir string
	refresh time.Duration
	log     *zap.Logger
	client  *http.Client
}

func NewManager(cfg config.BlocklistsConfig, log *zap.Logger) (*Manager, error) {
	m := &Manager{
		sets:    make(map[string]map[string]struct{}),
		feeds:   cfg.Feeds,
		dataDir: cfg.DataDir,
		refresh: cfg.RefreshEvery,
		log:     log,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
	if err := m.fetchAll(); err != nil {
		return nil, err
	}
	return m, nil
}

// IsBlocked returns (blocked bool, category string).
// Checks exact domain and all parent domains (subdomain coverage).
func (m *Manager) IsBlocked(domain string) (bool, string) {
	domain = strings.ToLower(strings.TrimSuffix(domain, "."))
	m.mu.RLock()
	defer m.mu.RUnlock()

	parts := strings.Split(domain, ".")
	for cat, set := range m.sets {
		// Exact match
		if _, ok := set[domain]; ok {
			return true, cat
		}
		// Parent domain match: sub.evil.com -> evil.com
		for i := 1; i < len(parts)-1; i++ {
			parent := strings.Join(parts[i:], ".")
			if _, ok := set[parent]; ok {
				return true, cat
			}
		}
	}
	return false, ""
}

// Stats returns domain count per category — safe to expose publicly.
func (m *Manager) Stats() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]int, len(m.sets))
	for cat, set := range m.sets {
		out[cat] = len(set)
	}
	return out
}

// StartRefreshLoop runs fetchAll on the configured interval.
func (m *Manager) StartRefreshLoop(ctx context.Context) {
	ticker := time.NewTicker(m.refresh)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.fetchAll(); err != nil {
				m.log.Error("blocklist refresh failed", zap.Error(err))
			}
		}
	}
}

func (m *Manager) fetchAll() error {
	newSets := make(map[string]map[string]struct{})

	for _, feed := range m.feeds {
		if !feed.Enabled {
			continue
		}
		domains, err := m.fetchFeed(feed)
		if err != nil {
			m.log.Warn("failed to fetch feed",
				zap.String("name", feed.Name),
				zap.Error(err),
			)
			continue
		}
		if _, ok := newSets[feed.Category]; !ok {
			newSets[feed.Category] = make(map[string]struct{})
		}
		for _, d := range domains {
			newSets[feed.Category][d] = struct{}{}
		}
		m.log.Info("loaded feed",
			zap.String("name", feed.Name),
			zap.String("category", feed.Category),
			zap.Int("count", len(domains)),
		)
	}

	m.mu.Lock()
	m.sets = newSets
	m.mu.Unlock()
	return nil
}

func (m *Manager) fetchFeed(feed config.FeedConfig) ([]string, error) {
	resp, err := m.client.Get(feed.URL)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	switch feed.Format {
	case "hosts":
		return parseHosts(resp.Body), nil
	case "abp":
		return parseABP(resp.Body), nil
	default:
		return parseDomains(resp.Body), nil
	}
}

// parseHosts handles /etc/hosts format: "0.0.0.0 domain.com"
func parseHosts(r io.Reader) []string {
	var out []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			d := strings.ToLower(fields[1])
			if isValidDomain(d) {
				out = append(out, d)
			}
		}
	}
	return out
}

// parseDomains handles one domain per line
func parseDomains(r io.Reader) []string {
	var out []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}
		d := strings.ToLower(line)
		if isValidDomain(d) {
			out = append(out, d)
		}
	}
	return out
}

// parseABP handles AdBlock Plus format: "||domain.com^"
func parseABP(r io.Reader) []string {
	var out []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if !strings.HasPrefix(line, "||") {
			continue
		}
		line = strings.TrimPrefix(line, "||")
		line = strings.Split(line, "^")[0]
		line = strings.ToLower(line)
		if isValidDomain(line) {
			out = append(out, line)
		}
	}
	return out
}

func isValidDomain(d string) bool {
	if d == "" || d == "localhost" {
		return false
	}
	if strings.ContainsAny(d, " /:\\") {
		return false
	}
	return strings.Contains(d, ".")
}
