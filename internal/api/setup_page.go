package api

import (
	_ "embed"
	"strings"

	"github.com/gofiber/fiber/v2"
)

//go:embed setup.html
var setupHTMLTemplate string

func (s *Server) handleSetupPage(c *fiber.Ctx) error {
	host, ip, _ := s.serverAddrs(c)
	html := strings.NewReplacer(
		"{{HOST}}", host,
		"{{IP}}", ip,
		"{{SETUP_BASE}}", c.BaseURL()+"/setup",
	).Replace(setupHTMLTemplate)
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(html)
}
