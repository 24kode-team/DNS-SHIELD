# DNS-SHIELD — دستورات Deploy

## ── ۱. Push به GitHub ────────────────────────────────────────────────────────

cd /Users/titosmcbook/TITO/DNS-SHIELD

# اگه remote هنوز set نشده:
git remote add origin https://github.com/24kode-team/DNS-SHIELD.git

# یا اگه قبلاً add شده ولی URL اشتباهه:
git remote set-url origin https://github.com/24kode-team/DNS-SHIELD.git

# Stage همه فایل‌ها
git add .

# Commit
git commit -m "feat: complete DNS-SHIELD v1.0.0

Core:
- DNS filtering engine (UDP/TCP/DoT/DoH)
- 8 block categories: phishing, malware, scam, porn, gambling, predatory, deepfake, sextortion
- Filter engine with allowlist (canada.ca, gc.ca)
- Blocklist manager with 24h auto-refresh

API & Dashboard:
- Built-in web dashboard at :8080
- REST API (health, metrics, admin)
- DoH handler (RFC 8484) — GET + POST
- Device setup endpoints: /setup/*

Device Profiles:
- iOS/macOS .mobileconfig (one-tap DoH install)
- Android Private DNS hostname
- Windows PowerShell setup + remove scripts
- Linux shell script (systemd-resolved / NetworkManager / resolv.conf)
- Router instructions

Infrastructure:
- Dockerfile (multi-stage, non-root)
- docker-compose.yml
- One-command installer: install.sh (Linux/macOS) + install.bat (Windows)
- GitHub Actions: CI (test + lint) + Release (5 platforms + Docker)
- .golangci.yml lint config"

# Push
git push -u origin main


## ── ۲. اول Tag بزن — Release خودکار می‌شه ───────────────────────────────────

git tag v1.0.0 -m "DNS-SHIELD v1.0.0 — Initial release"
git push origin v1.0.0

# GitHub Actions بعد از این:
# ✅ Linux amd64 + arm64 build می‌کنه
# ✅ macOS amd64 + arm64 build می‌کنه
# ✅ Windows amd64 build می‌کنه
# ✅ SHA-256 checksums می‌سازه
# ✅ GitHub Release می‌سازه با همه فایل‌ها
# ✅ Docker image به ghcr.io push می‌کنه


## ── ۳. بعد از Release — سرور را deploy کن ───────────────────────────────────

# روی سرور کانادایی (Ubuntu 22.04):
curl -fsSL https://raw.githubusercontent.com/24kode-team/DNS-SHIELD/main/install.sh | sudo bash

# یا Docker:
docker run -d \
  --name dns-shield \
  --restart unless-stopped \
  -p 53:53/udp \
  -p 53:53/tcp \
  -p 853:853/tcp \
  -p 8080:8080 \
  -e SHIELD_ADMIN_TOKEN=$(openssl rand -hex 20) \
  ghcr.io/24kode-team/dns-shield:latest


## ── ۴. DNS Record برای سرور ──────────────────────────────────────────────────

# در DNS panel دامنه‌ات:
# A   shield.nationcode.ca   →   YOUR_SERVER_IP
# یا
# A   dns.dns-shield.ca      →   YOUR_SERVER_IP


## ── ۵. NationCode page را با IP واقعی update کن ─────────────────────────────

# در /NATION-CODE/app/dns-shield/page.tsx این سه خط را عوض کن:
const DNS_SERVER = 'shield.nationcode.ca'    # ← دامنه واقعی
const DNS_IP     = 'YOUR_SERVER_IP'          # ← IP واقعی سرور
const DOH_URL    = 'https://shield.nationcode.ca/dns-query'
const SETUP_BASE = 'https://shield.nationcode.ca/setup'


## ── Verify بعد از deploy ─────────────────────────────────────────────────────

# DNS کار می‌کنه؟
dig @YOUR_SERVER_IP canada.ca

# Dashboard باز میشه؟
curl http://YOUR_SERVER_IP:8080/health

# iOS profile دانلود میشه؟
curl -I https://YOUR_DOMAIN/setup/ios
# باید برگردونه: Content-Type: application/x-apple-aspen-config
