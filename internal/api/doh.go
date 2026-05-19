package api

import (
	"encoding/base64"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

// registerDoH adds DNS-over-HTTPS endpoints (RFC 8484).
//
//	GET  /dns-query?dns=<base64url>
//	POST /dns-query  Content-Type: application/dns-message
func (s *Server) registerDoH(app *fiber.App) {
	app.Get("/dns-query", s.handleDoHGet)
	app.Post("/dns-query", s.handleDoHPost)
}

func (s *Server) handleDoHGet(c *fiber.Ctx) error {
	param := c.Query("dns")
	if param == "" {
		return c.Status(fiber.StatusBadRequest).SendString("missing dns parameter")
	}
	wire, err := base64.RawURLEncoding.DecodeString(param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid base64url")
	}
	return s.processDoH(c, wire)
}

func (s *Server) handleDoHPost(c *fiber.Ctx) error {
	if c.Get("Content-Type") != "application/dns-message" {
		return c.Status(fiber.StatusUnsupportedMediaType).
			SendString("Content-Type must be application/dns-message")
	}
	return s.processDoH(c, c.Body())
}

func (s *Server) processDoH(c *fiber.Ctx, wire []byte) error {
	start := time.Now()

	req := new(dns.Msg)
	if err := req.Unpack(wire); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid DNS message")
	}

	var respMsg *dns.Msg

	if len(req.Question) == 0 {
		respMsg = new(dns.Msg)
		respMsg.SetRcode(req, dns.RcodeFormatError)
	} else {
		q := req.Question[0]
		decision := s.engine.Evaluate(q.Name)
		s.metrics.RecordQuery(decision.Action, decision.Category)

		if decision.Action == "block" {
			s.log.Debug("DoH blocked", zap.String("domain", q.Name), zap.String("category", decision.Category))
			respMsg = dohBlockResponse(req, decision.BlockIP)
		} else {
			upstream, err := s.dohForward(req)
			if err != nil {
				s.log.Error("DoH upstream failed", zap.Error(err))
				respMsg = new(dns.Msg)
				respMsg.SetRcode(req, dns.RcodeServerFailure)
			} else {
				respMsg = upstream
			}
		}
	}

	s.metrics.RecordLatency(time.Since(start))

	packed, err := respMsg.Pack()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("pack error")
	}

	c.Set("Content-Type", "application/dns-message")
	c.Set("Cache-Control", "no-store")
	return c.Status(fiber.StatusOK).Send(packed)
}

func (s *Server) dohForward(req *dns.Msg) (*dns.Msg, error) {
	c := &dns.Client{Timeout: 3 * time.Second}
	var lastErr error
	for _, upstream := range s.upstreams {
		resp, _, err := c.Exchange(req, upstream)
		if err == nil {
			return resp, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

func dohBlockResponse(req *dns.Msg, blockIP string) *dns.Msg {
	resp := new(dns.Msg)
	resp.SetReply(req)
	resp.RecursionAvailable = true

	if len(req.Question) > 0 {
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
	}
	return resp
}
