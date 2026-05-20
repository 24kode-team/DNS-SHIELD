package api

import (
	"fmt"

	"github.com/24kode-team/DNS-SHIELD/internal/api/profiles"
	"github.com/gofiber/fiber/v2"
)

// registerSetup adds all /setup/* download endpoints.
// These serve ready-to-use config files for each platform — no App Store needed.
func (s *Server) registerSetup(app *fiber.App) {

	// ── iOS / macOS — .mobileconfig (one tap install) ────────────────────
	app.Get("/setup/ios", s.handleSetupIOS)
	app.Get("/setup/macos", s.handleSetupMacOS)    // same file, different name
	app.Get("/setup/ios-dot", s.handleSetupIOSDoT) // DoT variant

	// ── Windows — PowerShell scripts ──────────────────────────────────────
	app.Get("/setup/windows", s.handleSetupWindows)
	app.Get("/setup/windows-remove", s.handleSetupWindowsRemove)

	// ── Linux — shell scripts ─────────────────────────────────────────────
	app.Get("/setup/linux", s.handleSetupLinux)
	app.Get("/setup/linux-remove", s.handleSetupLinuxRemove)

	// ── Android — plain text hostname for Private DNS ─────────────────────
	app.Get("/setup/android", s.handleSetupAndroid)

	// ── Setup page (HTML) — shows all options ─────────────────────────────
	app.Get("/setup", s.handleSetupPage)
}

// ── iOS / macOS ───────────────────────────────────────────────────────────────

func (s *Server) handleSetupIOS(c *fiber.Ctx) error {
	host, ip, doh := s.serverAddrs(c)
	content := profiles.MobileConfig(host, ip, doh)
	c.Set("Content-Type", "application/x-apple-aspen-config")
	c.Set("Content-Disposition", `attachment; filename="dns-shield.mobileconfig"`)
	return c.SendString(content)
}

func (s *Server) handleSetupMacOS(c *fiber.Ctx) error {
	return s.handleSetupIOS(c) // same file format
}

func (s *Server) handleSetupIOSDoT(c *fiber.Ctx) error {
	host, ip, _ := s.serverAddrs(c)
	content := profiles.MobileConfigDoT(host, ip)
	c.Set("Content-Type", "application/x-apple-aspen-config")
	c.Set("Content-Disposition", `attachment; filename="dns-shield-dot.mobileconfig"`)
	return c.SendString(content)
}

// ── Windows ───────────────────────────────────────────────────────────────────

func (s *Server) handleSetupWindows(c *fiber.Ctx) error {
	ip, _, _ := s.serverAddrs(c)
	content := profiles.WindowsPowerShell(ip, "9.9.9.9")
	c.Set("Content-Type", "text/plain; charset=utf-8")
	c.Set("Content-Disposition", `attachment; filename="dns-shield-setup.ps1"`)
	return c.SendString(content)
}

func (s *Server) handleSetupWindowsRemove(c *fiber.Ctx) error {
	content := profiles.WindowsPowerShellRemove()
	c.Set("Content-Type", "text/plain; charset=utf-8")
	c.Set("Content-Disposition", `attachment; filename="dns-shield-remove.ps1"`)
	return c.SendString(content)
}

// ── Linux ─────────────────────────────────────────────────────────────────────

func (s *Server) handleSetupLinux(c *fiber.Ctx) error {
	ip, _, _ := s.serverAddrs(c)
	content := profiles.LinuxScript(ip, "9.9.9.9")
	c.Set("Content-Type", "text/plain; charset=utf-8")
	c.Set("Content-Disposition", `attachment; filename="dns-shield-setup.sh"`)
	return c.SendString(content)
}

func (s *Server) handleSetupLinuxRemove(c *fiber.Ctx) error {
	content := profiles.LinuxScriptRemove()
	c.Set("Content-Type", "text/plain; charset=utf-8")
	c.Set("Content-Disposition", `attachment; filename="dns-shield-remove.sh"`)
	return c.SendString(content)
}

// ── Android ───────────────────────────────────────────────────────────────────

func (s *Server) handleSetupAndroid(c *fiber.Ctx) error {
	host, _, _ := s.serverAddrs(c)
	// Just return the hostname as plain text — user copies into Private DNS
	c.Set("Content-Type", "text/plain; charset=utf-8")
	return c.SendString(host)
}

// ── Setup page ────────────────────────────────────────────────────────────────
// handleSetupPage is defined in setup_page.go

// serverAddrs extracts server host, IP, and DoH URL from the request.
func (s *Server) serverAddrs(c *fiber.Ctx) (host, ip, doh string) {
	host = c.Hostname()
	ip = host
	doh = fmt.Sprintf("%s/dns-query", c.BaseURL())
	return
}
