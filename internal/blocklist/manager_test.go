package blocklist

import (
	"strings"
	"testing"
)

func TestParseHosts(t *testing.T) {
	input := `# comment
0.0.0.0 evil.com
0.0.0.0 malware.net
# another comment
127.0.0.1 localhost
`
	r := strings.NewReader(input)
	domains := parseHosts(r)

	if len(domains) != 2 {
		t.Fatalf("expected 2 domains, got %d", len(domains))
	}
}

func TestParseDomains(t *testing.T) {
	input := `# header
evil.com
phishing.net
! adblock comment
valid.org
`
	r := strings.NewReader(input)
	domains := parseDomains(r)

	if len(domains) != 3 {
		t.Fatalf("expected 3 domains, got %d", len(domains))
	}
}

func TestParseABP(t *testing.T) {
	input := `[Adblock Plus]
||evil.com^
||malware.net^$third-party
##.ad-class
`
	r := strings.NewReader(input)
	domains := parseABP(r)

	if len(domains) != 2 {
		t.Fatalf("expected 2 domains, got %d", len(domains))
	}
}

func TestIsValidDomain(t *testing.T) {
	cases := []struct {
		domain string
		valid  bool
	}{
		{"evil.com", true},
		{"sub.domain.ca", true},
		{"localhost", false},
		{"", false},
		{"no-dot", false},
		{"has space.com", false},
		{"has/slash.com", false},
	}
	for _, c := range cases {
		got := isValidDomain(c.domain)
		if got != c.valid {
			t.Errorf("isValidDomain(%q) = %v, want %v", c.domain, got, c.valid)
		}
	}
}
