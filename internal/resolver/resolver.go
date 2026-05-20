package resolver

import (
	"context"
	"net"
	"time"

	"github.com/24kode-team/DNS-SHIELD/internal/config"
	"github.com/24kode-team/DNS-SHIELD/internal/filter"
	"github.com/24kode-team/DNS-SHIELD/internal/metrics"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

// Resolver handles DNS queries over UDP/TCP (port 53) and DoT (port 853).
type Resolver struct {
	cfg       config.ResolverConfig
	engine    *filter.Engine
	metrics   *metrics.Metrics
	log       *zap.Logger
	upstreams []string

	udpServer *dns.Server
	tcpServer *dns.Server
	dotServer *dns.Server
}

func New(cfg config.ResolverConfig, engine *filter.Engine, m *metrics.Metrics, log *zap.Logger) (*Resolver, error) {
	return &Resolver{
		cfg:       cfg,
		engine:    engine,
		metrics:   m,
		log:       log,
		upstreams: cfg.Upstreams,
	}, nil
}

func (r *Resolver) ListenAndServe(ctx context.Context) error {
	mux := dns.NewServeMux()
	mux.HandleFunc(".", r.handleQuery)

	r.udpServer = &dns.Server{
		Addr:         r.cfg.ListenAddr,
		Net:          "udp",
		Handler:      mux,
		ReadTimeout:  r.cfg.ReadTimeout,
		WriteTimeout: r.cfg.ReadTimeout,
	}
	r.tcpServer = &dns.Server{
		Addr:         r.cfg.ListenAddr,
		Net:          "tcp",
		Handler:      mux,
		ReadTimeout:  r.cfg.ReadTimeout,
		WriteTimeout: r.cfg.ReadTimeout,
	}

	errCh := make(chan error, 3)
	go func() { errCh <- r.udpServer.ListenAndServe() }()
	go func() { errCh <- r.tcpServer.ListenAndServe() }()

	if r.cfg.DoTAddr != "" && r.cfg.TLSCert != "" {
		tlsCfg, err := loadTLS(r.cfg.TLSCert, r.cfg.TLSKey)
		if err != nil {
			r.log.Warn("DoT disabled — TLS load failed", zap.Error(err))
		} else {
			r.dotServer = &dns.Server{
				Addr:      r.cfg.DoTAddr,
				Net:       "tcp-tls",
				Handler:   mux,
				TLSConfig: tlsCfg,
			}
			go func() { errCh <- r.dotServer.ListenAndServe() }()
		}
	}

	r.log.Info("DNS resolver listening",
		zap.String("udp/tcp", r.cfg.ListenAddr),
		zap.String("dot", r.cfg.DoTAddr),
	)

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func (r *Resolver) handleQuery(w dns.ResponseWriter, req *dns.Msg) {
	start := time.Now()

	if len(req.Question) == 0 {
		resp := new(dns.Msg)
		resp.SetRcode(req, dns.RcodeFormatError)
		if err := w.WriteMsg(resp); err != nil {
			r.log.Error("write failed", zap.Error(err))
		}
		return
	}

	q := req.Question[0]
	decision := r.engine.Evaluate(q.Name)
	r.metrics.RecordQuery(decision.Action, decision.Category)

	var resp *dns.Msg
	if decision.Action == "block" {
		resp = r.buildBlockResponse(req, decision.BlockIP)
	} else {
		upstream, err := r.forward(req)
		if err != nil {
			r.log.Error("upstream failed", zap.Error(err))
			resp = new(dns.Msg)
			resp.SetRcode(req, dns.RcodeServerFailure)
		} else {
			resp = upstream
		}
	}

	r.metrics.RecordLatency(time.Since(start))
	if err := w.WriteMsg(resp); err != nil {
		r.log.Error("write failed", zap.Error(err))
	}
}

func (r *Resolver) buildBlockResponse(req *dns.Msg, blockIP string) *dns.Msg {
	resp := new(dns.Msg)
	resp.SetReply(req)
	resp.RecursionAvailable = true

	q := req.Question[0]
	if q.Qtype == dns.TypeA && blockIP != "" {
		ip := net.ParseIP(blockIP).To4()
		if ip != nil {
			resp.Answer = append(resp.Answer, &dns.A{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				A: ip,
			})
			return resp
		}
	}
	resp.SetRcode(req, dns.RcodeNameError)
	return resp
}

func (r *Resolver) forward(req *dns.Msg) (*dns.Msg, error) {
	c := &dns.Client{Timeout: 3 * time.Second}
	var lastErr error
	for _, upstream := range r.upstreams {
		resp, _, err := c.Exchange(req, upstream)
		if err == nil {
			return resp, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func (r *Resolver) Shutdown() {
	if r.udpServer != nil {
		if err := r.udpServer.Shutdown(); err != nil {
			_ = err
		}
	}
	if r.tcpServer != nil {
		if err := r.tcpServer.Shutdown(); err != nil {
			_ = err
		}
	}
	if r.dotServer != nil {
		if err := r.dotServer.Shutdown(); err != nil {
			_ = err
		}
	}
}
