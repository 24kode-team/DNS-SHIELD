package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Resolver   ResolverConfig   `yaml:"resolver"`
	Filter     FilterConfig     `yaml:"filter"`
	Blocklists BlocklistsConfig `yaml:"blocklists"`
	API        APIConfig        `yaml:"api"`
}

type ResolverConfig struct {
	ListenAddr  string        `yaml:"listen_addr"`
	DoTAddr     string        `yaml:"dot_addr"`
	DoHPath     string        `yaml:"doh_path"`
	Upstreams   []string      `yaml:"upstreams"`
	TLSCert     string        `yaml:"tls_cert"`
	TLSKey      string        `yaml:"tls_key"`
	ReadTimeout time.Duration `yaml:"read_timeout"`
}

type FilterConfig struct {
	BlockPage  string   `yaml:"block_page"`
	Categories []string `yaml:"categories"`
	Allowlist  []string `yaml:"allowlist"`
}

type BlocklistsConfig struct {
	DataDir      string        `yaml:"data_dir"`
	RefreshEvery time.Duration `yaml:"refresh_every"`
	Feeds        []FeedConfig  `yaml:"feeds"`
}

type FeedConfig struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Category string `yaml:"category"`
	Format   string `yaml:"format"` // hosts | domains | abp
	Enabled  bool   `yaml:"enabled"`
}

type APIConfig struct {
	ListenAddr string `yaml:"listen_addr"`
	AdminToken string `yaml:"admin_token"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}

	if tok := os.Getenv("SHIELD_ADMIN_TOKEN"); tok != "" {
		cfg.API.AdminToken = tok
	}

	return cfg, nil
}
