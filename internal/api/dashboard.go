package api

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"
)

//go:embed dashboard.html
var dashboardHTML string

func (s *Server) handleDashboard(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(dashboardHTML)
}
