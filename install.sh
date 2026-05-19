#!/usr/bin/env bash
# DNS-SHIELD Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/dns-shield/shield/main/install.sh | bash

set -euo pipefail

# ── Colours ───────────────────────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
BLUE='\033[0;34m'; BOLD='\033[1m'; RESET='\033[0m'

info()    { echo -e "${BLUE}[DNS-SHIELD]${RESET} $*"; }
success() { echo -e "${GREEN}[DNS-SHIELD]${RESET} $*"; }
warn()    { echo -e "${YELLOW}[DNS-SHIELD]${RESET} $*"; }
error()   { echo -e "${RED}[DNS-SHIELD]${RESET} $*" >&2; exit 1; }

REPO="https://github.com/dns-shield/shield"
INSTALL_DIR="/opt/dns-shield"
BIN_DIR="/usr/local/bin"
CONFIG_DIR="/etc/dns-shield"
DATA_DIR="/var/lib/dns-shield"
LOG_DIR="/var/log/dns-shield"
SERVICE_USER="dns-shield"

# ── Detect OS / Arch ──────────────────────────────────────────────────────────
detect_platform() {
  OS="$(uname -s)"
  ARCH="$(uname -m)"

  case "$OS" in
    Linux)  PLATFORM="linux" ;;
    Darwin) PLATFORM="darwin" ;;
    *)      error "Unsupported OS: $OS. Use Docker instead." ;;
  esac

  case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) error "Unsupported architecture: $ARCH" ;;
  esac

  info "Detected platform: ${PLATFORM}/${ARCH}"
}

# ── Check dependencies ─────────────────────────────────────────────────────────
check_deps() {
  local missing=()
  for cmd in curl tar; do
    command -v "$cmd" &>/dev/null || missing+=("$cmd")
  done
  [[ ${#missing[@]} -eq 0 ]] || error "Missing required tools: ${missing[*]}"
}

# ── Get latest release version ─────────────────────────────────────────────────
get_latest_version() {
  VERSION=$(curl -fsSL "https://api.github.com/repos/dns-shield/shield/releases/latest" \
    | grep '"tag_name"' | sed 's/.*"tag_name": *"\(.*\)".*/\1/')
  [[ -n "$VERSION" ]] || VERSION="v1.0.0"
  info "Installing DNS-SHIELD $VERSION"
}

# ── Download binary ────────────────────────────────────────────────────────────
download_binary() {
  local url="${REPO}/releases/download/${VERSION}/dns-shield_${PLATFORM}_${ARCH}.tar.gz"
  local tmp="$(mktemp -d)"

  info "Downloading from $url"
  curl -fsSL "$url" -o "$tmp/dns-shield.tar.gz" \
    || error "Download failed. Check your internet connection."

  tar -xzf "$tmp/dns-shield.tar.gz" -C "$tmp"
  install -m 755 "$tmp/dns-shield" "$BIN_DIR/dns-shield"
  rm -rf "$tmp"

  success "Binary installed at $BIN_DIR/dns-shield"
}

# ── Create system user ─────────────────────────────────────────────────────────
create_user() {
  if ! id "$SERVICE_USER" &>/dev/null; then
    info "Creating system user: $SERVICE_USER"
    if [[ "$PLATFORM" == "linux" ]]; then
      useradd --system --no-create-home --shell /bin/false "$SERVICE_USER"
    else
      dscl . -create "/Users/$SERVICE_USER" 2>/dev/null || true
    fi
  fi
}

# ── Create directories ─────────────────────────────────────────────────────────
create_dirs() {
  for dir in "$CONFIG_DIR/tls" "$DATA_DIR/blocklists" "$LOG_DIR"; do
    mkdir -p "$dir"
  done
  chown -R "$SERVICE_USER" "$DATA_DIR" "$LOG_DIR" 2>/dev/null || true
  info "Created directories"
}

# ── Install config ─────────────────────────────────────────────────────────────
install_config() {
  if [[ -f "$CONFIG_DIR/shield.yaml" ]]; then
    warn "Config already exists at $CONFIG_DIR/shield.yaml — skipping"
    return
  fi

  cat > "$CONFIG_DIR/shield.yaml" << 'YAML'
resolver:
  listen_addr: "0.0.0.0:53"
  dot_addr: "0.0.0.0:853"
  doh_path: "/dns-query"
  tls_cert: "/etc/dns-shield/tls/cert.pem"
  tls_key: "/etc/dns-shield/tls/key.pem"
  read_timeout: 5s
  upstreams:
    - "9.9.9.9:53"
    - "149.112.112.112:53"
    - "1.1.1.1:53"

filter:
  block_page: "0.0.0.0"
  categories:
    - phishing
    - malware
    - scam
    - porn
    - gambling
    - predatory
    - deepfake
    - sextortion
  allowlist:
    - "canada.ca"
    - "gc.ca"

api:
  listen_addr: "0.0.0.0:8080"

blocklists:
  data_dir: "/var/lib/dns-shield/blocklists"
  refresh_every: 24h
  feeds:
    - name: "PhishingArmy"
      url: "https://phishing.army/download/phishing_army_blocklist_extended.txt"
      category: phishing
      format: domains
      enabled: true
    - name: "URLhaus"
      url: "https://urlhaus.abuse.ch/downloads/hostfile/"
      category: malware
      format: hosts
      enabled: true
    - name: "Hagezi Pro"
      url: "https://cdn.jsdelivr.net/gh/hagezi/dns-blocklists@latest/domains/pro.txt"
      category: malware
      format: domains
      enabled: true
    - name: "ScamBlocker"
      url: "https://raw.githubusercontent.com/durablenapkin/scamblocklist/master/hosts.txt"
      category: scam
      format: hosts
      enabled: true
    - name: "StevenBlack Adult"
      url: "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/porn/hosts"
      category: porn
      format: hosts
      enabled: true
    - name: "Gambling Block"
      url: "https://raw.githubusercontent.com/nicehash/nicehash-antivirus-lists/master/gambling.txt"
      category: gambling
      format: domains
      enabled: true
    - name: "CleanBrowsing Family"
      url: "https://download.cleanbrowsing.org/domains/blocked-family.txt"
      category: predatory
      format: domains
      enabled: true
YAML

  success "Config installed at $CONFIG_DIR/shield.yaml"
}

# ── Generate admin token ───────────────────────────────────────────────────────
generate_token() {
  if [[ -f "$CONFIG_DIR/.env" ]]; then
    return
  fi
  local token
  token=$(LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 40 || true)
  echo "SHIELD_ADMIN_TOKEN=${token}" > "$CONFIG_DIR/.env"
  chmod 600 "$CONFIG_DIR/.env"
  success "Admin token saved to $CONFIG_DIR/.env"
}

# ── Install systemd service (Linux) ───────────────────────────────────────────
install_systemd() {
  [[ "$PLATFORM" != "linux" ]] && return
  command -v systemctl &>/dev/null || return

  info "Installing systemd service"
  cat > /etc/systemd/system/dns-shield.service << SERVICE
[Unit]
Description=DNS-SHIELD — Privacy-first DNS filtering for Canada
Documentation=https://github.com/dns-shield/shield
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=${SERVICE_USER}
EnvironmentFile=${CONFIG_DIR}/.env
ExecStart=${BIN_DIR}/dns-shield
Restart=always
RestartSec=5

# Security hardening
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=${DATA_DIR} ${LOG_DIR}
PrivateTmp=yes
PrivateDevices=yes

# Allow binding privileged ports
AmbientCapabilities=CAP_NET_BIND_SERVICE
CapabilityBoundingSet=CAP_NET_BIND_SERVICE

StandardOutput=journal
StandardError=journal
SyslogIdentifier=dns-shield

[Install]
WantedBy=multi-user.target
SERVICE

  systemctl daemon-reload
  systemctl enable dns-shield
  systemctl start dns-shield
  success "Service started: systemctl status dns-shield"
}

# ── Install launchd plist (macOS) ──────────────────────────────────────────────
install_launchd() {
  [[ "$PLATFORM" != "darwin" ]] && return

  info "Installing launchd service"
  cat > /Library/LaunchDaemons/com.dns-shield.plist << PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.dns-shield</string>
  <key>ProgramArguments</key>
  <array>
    <string>${BIN_DIR}/dns-shield</string>
  </array>
  <key>EnvironmentVariables</key>
  <dict>
    <key>SHIELD_ADMIN_TOKEN</key>
    <string>$(grep SHIELD_ADMIN_TOKEN "$CONFIG_DIR/.env" | cut -d= -f2)</string>
  </dict>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
  <key>StandardOutPath</key>
  <string>${LOG_DIR}/shield.log</string>
  <key>StandardErrorPath</key>
  <string>${LOG_DIR}/shield.error.log</string>
</dict>
</plist>
PLIST

  launchctl load /Library/LaunchDaemons/com.dns-shield.plist
  success "Service loaded via launchd"
}

# ── Wait for service health ────────────────────────────────────────────────────
wait_healthy() {
  info "Waiting for service to become healthy..."
  local retries=12
  while (( retries-- > 0 )); do
    if curl -fsSL http://127.0.0.1:8080/health &>/dev/null; then
      success "DNS-SHIELD is healthy!"
      return
    fi
    sleep 2
  done
  warn "Health check timed out — check logs: journalctl -u dns-shield -f"
}

# ── Print summary ──────────────────────────────────────────────────────────────
print_summary() {
  local token
  token=$(grep SHIELD_ADMIN_TOKEN "$CONFIG_DIR/.env" | cut -d= -f2)

  echo ""
  echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
  echo -e "${GREEN}${BOLD}  DNS-SHIELD installed successfully! 🛡️${RESET}"
  echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
  echo ""
  echo -e "  ${BOLD}DNS Server:${RESET}    127.0.0.1:53 (UDP/TCP)"
  echo -e "  ${BOLD}DoT:${RESET}           127.0.0.1:853"
  echo -e "  ${BOLD}DoH:${RESET}           http://127.0.0.1:8080/dns-query"
  echo -e "  ${BOLD}Dashboard:${RESET}     http://127.0.0.1:8080"
  echo -e "  ${BOLD}Admin Token:${RESET}   ${token}"
  echo ""
  echo -e "  ${BOLD}Config:${RESET}        $CONFIG_DIR/shield.yaml"
  echo -e "  ${BOLD}Logs:${RESET}          $LOG_DIR"
  echo ""
  echo -e "  ${BOLD}Quick test:${RESET}"
  echo -e "    dig @127.0.0.1 canada.ca       ${GREEN}# should resolve${RESET}"
  echo -e "    curl http://127.0.0.1:8080/metrics"
  echo ""
  echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
}

# ── Uninstall ──────────────────────────────────────────────────────────────────
uninstall() {
  info "Uninstalling DNS-SHIELD..."

  if [[ "$PLATFORM" == "linux" ]]; then
    systemctl stop dns-shield 2>/dev/null || true
    systemctl disable dns-shield 2>/dev/null || true
    rm -f /etc/systemd/system/dns-shield.service
    systemctl daemon-reload
  elif [[ "$PLATFORM" == "darwin" ]]; then
    launchctl unload /Library/LaunchDaemons/com.dns-shield.plist 2>/dev/null || true
    rm -f /Library/LaunchDaemons/com.dns-shield.plist
  fi

  rm -f "$BIN_DIR/dns-shield"
  rm -rf "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR"

  success "DNS-SHIELD uninstalled"
  exit 0
}

# ── Main ───────────────────────────────────────────────────────────────────────
main() {
  [[ "${EUID:-$(id -u)}" -eq 0 ]] || error "Please run as root: sudo bash install.sh"

  [[ "${1:-}" == "uninstall" ]] && uninstall

  echo ""
  echo -e "${BOLD}  DNS-SHIELD Installer${RESET}"
  echo -e "  Privacy-first DNS filtering — Canada"
  echo ""

  detect_platform
  check_deps
  get_latest_version
  download_binary
  create_user
  create_dirs
  install_config
  generate_token
  install_systemd
  install_launchd
  wait_healthy
  print_summary
}

main "$@"
