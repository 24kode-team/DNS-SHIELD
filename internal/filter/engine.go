package filter

import (
	"strings"

	"github.com/dns-shield/shield/internal/config"
	"go.uber.org/zap"
)

// Blocker is satisfied by blocklist.Manager
type Blocker interface {
	IsBlocked(domain string) (bool, string)
}

// Engine decides BLOCK or ALLOW for each DNS query.
type Engine struct {
	bl        Blocker
	cfg       config.FilterConfig
	allowlist map[string]struct{}
	log       *zap.Logger
}

// Decision is the result of evaluating a domain.
type Decision struct {
	Action   string // "allow" | "block"
	Category string // e.g. "malware"
	BlockIP  string // IP to return for blocked A queries (empty = NXDOMAIN)
}

func NewEngine(bl Blocker, cfg config.FilterConfig, log *zap.Logger) *Engine {
	al := make(map[string]struct{}, len(cfg.Allowlist))
	for _, d := range cfg.Allowlist {
		al[strings.ToLower(d)] = struct{}{}
	}
	return &Engine{bl: bl, cfg: cfg, allowlist: al, log: log}
}

// Evaluate returns a Decision for the given FQDN.
func (e *Engine) Evaluate(domain string) Decision {
	domain = strings.ToLower(strings.TrimSuffix(domain, "."))

	// 1. Allowlist wins first
	if e.isAllowed(domain) {
		return Decision{Action: "allow"}
	}

	// 2. Check blocklist
	if blocked, cat := e.bl.IsBlocked(domain); blocked {
		e.log.Debug("blocked query",
			zap.String("domain", domain),
			zap.String("category", cat),
		)
		return Decision{
			Action:   "block",
			Category: cat,
			BlockIP:  e.cfg.BlockPage,
		}
	}

	return Decision{Action: "allow"}
}

func (e *Engine) isAllowed(domain string) bool {
	if _, ok := e.allowlist[domain]; ok {
		return true
	}
	parts := strings.Split(domain, ".")
	for i := 1; i < len(parts); i++ {
		parent := strings.Join(parts[i:], ".")
		if _, ok := e.allowlist[parent]; ok {
			return true
		}
	}
	return false
}

// AddAllowlist adds a domain to the in-memory allowlist (hot reload — no restart needed).
func (e *Engine) AddAllowlist(domain string) {
	e.allowlist[strings.ToLower(domain)] = struct{}{}
}

// RemoveAllowlist removes a domain from the in-memory allowlist.
func (e *Engine) RemoveAllowlist(domain string) {
	delete(e.allowlist, strings.ToLower(domain))
}
