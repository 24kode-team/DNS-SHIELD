package api

import (
	"github.com/dns-shield/shield/internal/blocklist"
	"github.com/dns-shield/shield/internal/config"
	"github.com/dns-shield/shield/internal/filter"
	"github.com/dns-shield/shield/internal/metrics"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type Server struct {
	app       *fiber.App
	cfg       config.APIConfig
	bl        *blocklist.Manager
	engine    *filter.Engine
	metrics   *metrics.Metrics
	log       *zap.Logger
	upstreams []string
}

func New(cfg config.APIConfig, bl *blocklist.Manager, engine *filter.Engine, m *metrics.Metrics, log *zap.Logger, upstreams []string) *Server {
	s := &Server{
		cfg:       cfg,
		bl:        bl,
		engine:    engine,
		metrics:   m,
		log:       log,
		upstreams: upstreams,
	}
	s.app = fiber.New(fiber.Config{DisableStartupMessage: true})
	s.app.Use(recover.New())
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Authorization",
	}))

	s.registerRoutes()
	s.registerDoH(s.app)
	s.registerSetup(s.app)
	return s
}

func (s *Server) registerRoutes() {
	// ── Public ────────────────────────────────────────────────────────────
	s.app.Get("/health", s.handleHealth)
	s.app.Get("/metrics", s.handleMetrics)

	// ── Web Dashboard ─────────────────────────────────────────────────────
	s.app.Get("/", s.handleDashboard)

	// ── Admin (Bearer token required) ─────────────────────────────────────
	admin := s.app.Group("/admin", s.authMiddleware)
	admin.Get("/blocklist/stats", s.handleBlocklistStats)
	admin.Post("/allowlist", s.handleAddAllowlist)
	admin.Delete("/allowlist", s.handleRemoveAllowlist)
}

func (s *Server) handleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok", "service": "dns-shield"})
}

func (s *Server) handleMetrics(c *fiber.Ctx) error {
	return c.JSON(s.metrics.Snapshot())
}

func (s *Server) handleBlocklistStats(c *fiber.Ctx) error {
	return c.JSON(s.bl.Stats())
}

func (s *Server) handleAddAllowlist(c *fiber.Ctx) error {
	var body struct {
		Domain string `json:"domain"`
	}
	if err := c.BodyParser(&body); err != nil || body.Domain == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "domain required"})
	}
	s.engine.AddAllowlist(body.Domain)
	s.log.Info("allowlist: added", zap.String("domain", body.Domain))
	return c.JSON(fiber.Map{"ok": true, "domain": body.Domain})
}

func (s *Server) handleRemoveAllowlist(c *fiber.Ctx) error {
	var body struct {
		Domain string `json:"domain"`
	}
	if err := c.BodyParser(&body); err != nil || body.Domain == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "domain required"})
	}
	s.engine.RemoveAllowlist(body.Domain)
	s.log.Info("allowlist: removed", zap.String("domain", body.Domain))
	return c.JSON(fiber.Map{"ok": true, "domain": body.Domain})
}

func (s *Server) authMiddleware(c *fiber.Ctx) error {
	if c.Get("Authorization") != "Bearer "+s.cfg.AdminToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.Next()
}

func (s *Server) Start() error {
	s.log.Info("API + Dashboard + DoH + Setup listening", zap.String("addr", s.cfg.ListenAddr))
	return s.app.Listen(s.cfg.ListenAddr)
}

func (s *Server) Shutdown() {
	s.app.Shutdown()
}
