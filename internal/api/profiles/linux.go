package profiles

import "fmt"

// LinuxScript generates a shell script that sets DNS on Linux.
// Supports systemd-resolved, NetworkManager, and direct /etc/resolv.conf.
// User runs: sudo bash dns-shield-setup.sh
// To revert: sudo bash dns-shield-remove.sh
func LinuxScript(serverIP, fallbackIP string) string {
	return fmt.Sprintf(`#!/usr/bin/env bash
# DNS-SHIELD Setup — Linux
# Run: sudo bash dns-shield-setup.sh
# Revert: sudo bash dns-shield-remove.sh

set -euo pipefail

DNS_PRIMARY="%s"
DNS_FALLBACK="%s"
BACKUP="/etc/dns-shield-backup"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RESET='\033[0m'

[[ "${EUID:-$(id -u)}" -eq 0 ]] || { echo -e "${RED}Please run as root: sudo bash $0${RESET}"; exit 1; }

echo ""
echo -e "  ${GREEN}DNS-SHIELD — Setting up safe DNS...${RESET}"
echo ""

mkdir -p "$BACKUP"

# ── Method 1: systemd-resolved ──────────────────────────────────────────────
if systemctl is-active --quiet systemd-resolved 2>/dev/null; then
    echo -e "  ${YELLOW}Detected: systemd-resolved${RESET}"
    cp /etc/systemd/resolved.conf "$BACKUP/resolved.conf.bak" 2>/dev/null || true

    cat > /etc/systemd/resolved.conf.d/dns-shield.conf << EOF
[Resolve]
DNS=${DNS_PRIMARY} ${DNS_FALLBACK}
FallbackDNS=9.9.9.9
DNSStubListener=yes
EOF

    systemctl restart systemd-resolved
    echo -e "  ${GREEN}OK: systemd-resolved configured${RESET}"

# ── Method 2: NetworkManager ─────────────────────────────────────────────────
elif command -v nmcli &>/dev/null; then
    echo -e "  ${YELLOW}Detected: NetworkManager${RESET}"
    CONN=$(nmcli -t -f NAME,DEVICE,STATE con show --active | grep -v lo | head -1 | cut -d: -f1)
    if [[ -n "$CONN" ]]; then
        nmcli con mod "$CONN" ipv4.dns "${DNS_PRIMARY},${DNS_FALLBACK}"
        nmcli con mod "$CONN" ipv4.ignore-auto-dns yes
        nmcli con up "$CONN"
        echo -e "  ${GREEN}OK: NetworkManager connection '$CONN' configured${RESET}"
    fi

# ── Method 3: Direct resolv.conf ─────────────────────────────────────────────
else
    echo -e "  ${YELLOW}Fallback: writing /etc/resolv.conf directly${RESET}"
    cp /etc/resolv.conf "$BACKUP/resolv.conf.bak" 2>/dev/null || true

    # Make writable (some distros symlink it)
    chattr -i /etc/resolv.conf 2>/dev/null || true

    cat > /etc/resolv.conf << EOF
# DNS-SHIELD — Safe Canadian DNS
nameserver ${DNS_PRIMARY}
nameserver ${DNS_FALLBACK}
nameserver 9.9.9.9
EOF

    chattr +i /etc/resolv.conf 2>/dev/null || true
    echo -e "  ${GREEN}OK: /etc/resolv.conf updated${RESET}"
fi

# ── Verify ───────────────────────────────────────────────────────────────────
echo ""
echo -e "  Verifying..."
if dig +short +timeout=3 canada.ca @${DNS_PRIMARY} &>/dev/null; then
    echo -e "  ${GREEN}DNS-SHIELD is working!${RESET}"
else
    echo -e "  ${YELLOW}Could not verify — check server connectivity${RESET}"
fi

echo ""
echo -e "  ${GREEN}Done. DNS-SHIELD active.${RESET}"
echo -e "  To remove: sudo bash dns-shield-remove.sh"
echo ""
`, serverIP, fallbackIP)
}

// LinuxScriptRemove generates the revert script for Linux.
func LinuxScriptRemove() string {
	return `#!/usr/bin/env bash
# DNS-SHIELD Remove — Revert DNS to system defaults

set -euo pipefail
BACKUP="/etc/dns-shield-backup"

GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RESET='\033[0m'
[[ "${EUID:-$(id -u)}" -eq 0 ]] || { echo "Please run as root: sudo bash $0"; exit 1; }

echo ""
echo -e "  ${YELLOW}DNS-SHIELD — Reverting DNS...${RESET}"
echo ""

# systemd-resolved
if [[ -f /etc/systemd/resolved.conf.d/dns-shield.conf ]]; then
    rm -f /etc/systemd/resolved.conf.d/dns-shield.conf
    systemctl restart systemd-resolved
    echo -e "  ${GREEN}OK: systemd-resolved restored${RESET}"
fi

# resolv.conf backup
if [[ -f "$BACKUP/resolv.conf.bak" ]]; then
    chattr -i /etc/resolv.conf 2>/dev/null || true
    cp "$BACKUP/resolv.conf.bak" /etc/resolv.conf
    echo -e "  ${GREEN}OK: /etc/resolv.conf restored${RESET}"
fi

# NetworkManager — set back to auto
if command -v nmcli &>/dev/null; then
    CONN=$(nmcli -t -f NAME,DEVICE,STATE con show --active | grep -v lo | head -1 | cut -d: -f1)
    if [[ -n "$CONN" ]]; then
        nmcli con mod "$CONN" ipv4.dns ""
        nmcli con mod "$CONN" ipv4.ignore-auto-dns no
        nmcli con up "$CONN" 2>/dev/null || true
        echo -e "  ${GREEN}OK: NetworkManager DNS set to automatic${RESET}"
    fi
fi

echo ""
echo -e "  ${GREEN}DNS-SHIELD removed. DNS is now automatic.${RESET}"
echo ""
`
}

// AndroidInstructions returns the Private DNS hostname for Android 9+.
// No script needed — user pastes this into Settings → Private DNS.
func AndroidPrivateDNSHostname(serverHost string) string {
	return serverHost
}
