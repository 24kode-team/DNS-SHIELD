# DNS-SHIELD 🛡️

**Privacy-first DNS filtering server — built for Canada**

Blocks phishing, malware, scams, adult content, gambling, predatory sites, harmful AI/deepfake domains, and sextortion — with **zero tracking**, **no data sales**, and **minimal logging**.

---

## Install — One Command

### Linux / macOS
```bash
curl -fsSL https://raw.githubusercontent.com/dns-shield/shield/main/install.sh | sudo bash
```

### Windows (Run as Administrator)
```powershell
# Download and run:
# https://github.com/dns-shield/shield/releases/latest/download/install.bat
```

### Docker
```bash
docker run -d \
  --name dns-shield \
  -p 53:53/udp -p 53:53/tcp \
  -p 853:853/tcp \
  -p 8080:8080 \
  -e SHIELD_ADMIN_TOKEN=your-secret-token \
  ghcr.io/dns-shield/shield:latest
```

After install, open **http://YOUR_SERVER_IP:8080** to see the dashboard.

---

## What Gets Blocked

| Category | Sources |
|---|---|
| Phishing | PhishTank, OpenPhish, PhishingArmy |
| Malware | URLhaus, Hagezi Pro |
| Scam | ScamBlocker |
| Adult / Porn | StevenBlack Adult |
| Gambling | NiceHash Gambling List |
| Predatory | CleanBrowsing Family |
| Deepfake / Harmful AI | Community list |
| Sextortion | Community list |

---

## Privacy Guarantees

- **No per-query logging** — we never record which domain was queried
- **No IP logging** — client IPs are not stored
- **Aggregate stats only** — counters track totals, not individual requests
- **Canada-hosted** — stays within Canadian jurisdiction (PIPEDA compliant)
- **Open source** — fully auditable

---

## Protocol Support

| Protocol | Address | Use case |
|---|---|---|
| DNS (UDP/TCP) | `:53` | All devices, routers |
| DNS-over-TLS | `:853` | Android Private DNS, iOS |
| DNS-over-HTTPS | `:8080/dns-query` | Browsers, modern OS |

---

## Quick Setup After Install

### macOS
```bash
sudo networksetup -setdnsservers Wi-Fi YOUR_SERVER_IP
```

### Windows (PowerShell Admin)
```powershell
Set-DnsClientServerAddress -InterfaceAlias "Wi-Fi" -ServerAddresses ("YOUR_SERVER_IP")
```

### Linux
```bash
echo "nameserver YOUR_SERVER_IP" | sudo tee /etc/resolv.conf
```

### Router (protects ALL devices)
```
DHCP Primary DNS: YOUR_SERVER_IP
DHCP Secondary DNS: 9.9.9.9
```

### Android (Private DNS / DoT)
```
Settings → Network → Private DNS → YOUR_SERVER_HOSTNAME
```

---

## Dashboard

Visit `http://YOUR_SERVER_IP:8080` after install.

Shows:
- Total / Blocked / Allowed query counts
- Blocks by category (phishing, malware, etc.)
- Response latency breakdown
- Per-feed blocklist entry counts
- One-click setup instructions per OS

---

## API Reference

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/health` | Public | Service health |
| GET | `/metrics` | Public | Aggregate stats (no PII) |
| GET | `/dns-query` | Public | DoH (RFC 8484) |
| POST | `/dns-query` | Public | DoH (RFC 8484) |
| GET | `/admin/blocklist/stats` | Bearer | Feed domain counts |
| POST | `/admin/allowlist` | Bearer | Add allowlist domain |
| DELETE | `/admin/allowlist` | Bearer | Remove allowlist domain |

### Add domain to allowlist
```bash
curl -X POST http://localhost:8080/admin/allowlist \
  -H "Authorization: Bearer $SHIELD_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"domain":"trusted.example.ca"}'
```

---

## Uninstall

```bash
# Linux / macOS
curl -fsSL https://raw.githubusercontent.com/dns-shield/shield/main/install.sh | sudo bash -s uninstall

# Windows (Admin PowerShell)
sc stop dns-shield && sc delete dns-shield
```

---

## Build from Source

```bash
git clone https://github.com/dns-shield/shield
cd shield
make deps
make build
sudo make run
```

---

## Project Structure

```
DNS-SHIELD/
├── cmd/shield/main.go           Entry point
├── internal/
│   ├── config/                  YAML config loader
│   ├── blocklist/               Feed fetcher + domain set manager
│   ├── filter/                  Allow/block decision engine
│   ├── resolver/                DNS server (UDP, TCP, DoT)
│   ├── metrics/                 Zero-PII aggregate counters
│   └── api/
│       ├── server.go            REST API + middleware
│       ├── doh.go               DNS-over-HTTPS handler (RFC 8484)
│       └── dashboard.go         Built-in web dashboard
├── configs/shield.yaml          Configuration + feed list
├── deploy/docker-compose.yml
├── install.sh                   Linux/macOS one-line installer
├── install.bat                  Windows installer
└── Makefile
```

---

## Roadmap

- [x] DNS filtering engine (UDP/TCP)
- [x] DNS-over-TLS (DoT)
- [x] DNS-over-HTTPS (DoH, RFC 8484)
- [x] REST API + admin endpoints
- [x] Built-in web dashboard
- [x] One-command installer (Linux/macOS/Windows)
- [ ] Kubernetes Helm chart
- [ ] Anycast setup (Canadian PoPs)
- [ ] iOS .mobileconfig profile generator
- [ ] CCCS / CRTC compliance documentation

---

## License

MIT
