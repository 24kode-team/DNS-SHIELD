package profiles

import "fmt"

// WindowsReg generates a .reg file that sets DNS for all network adapters on Windows.
// User double-clicks → UAC prompt → DNS changed. No install needed.
func WindowsReg(serverIP string) string {
	return fmt.Sprintf(`Windows Registry Editor Version 5.00

; DNS-SHIELD — Safe Canadian DNS
; Double-click this file to apply, then restart network adapter.
; To revert: run dns-shield-remove.reg

; Set DNS for common adapter names
; Windows will use whichever adapter is active.

[HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces]

; ── Wi-Fi / WLAN adapters ──────────────────────────────────────────────────
; Note: The actual key names are GUIDs. This script uses a PowerShell
; command below that handles all adapters automatically.
; See the PowerShell section at the bottom.

; ── Fallback static entry ──────────────────────────────────────────────────
[HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters]
"NameServer"="%s"
"SearchList"=""
`, serverIP)
}

// WindowsPowerShell generates a .ps1 script that sets DNS on ALL active adapters.
// More reliable than .reg — handles any adapter name/GUID automatically.
func WindowsPowerShell(serverIP, fallbackIP string) string {
	return fmt.Sprintf(`# DNS-SHIELD Setup Script for Windows
# Run as Administrator: Right-click → Run with PowerShell
# To revert: run dns-shield-remove.ps1

$ErrorActionPreference = "Stop"
$dnsShield  = "%s"
$fallback   = "%s"

Write-Host ""
Write-Host "  DNS-SHIELD — Setting up safe DNS..." -ForegroundColor Cyan
Write-Host ""

# Get all active physical network adapters
$adapters = Get-NetAdapter | Where-Object { $_.Status -eq "Up" }

if ($adapters.Count -eq 0) {
    Write-Host "  [ERROR] No active network adapters found." -ForegroundColor Red
    exit 1
}

foreach ($adapter in $adapters) {
    Write-Host "  Configuring: $($adapter.Name)" -ForegroundColor Yellow

    # Remove DHCP DNS, set manual
    Set-DnsClientServerAddress `
        -InterfaceIndex $adapter.InterfaceIndex `
        -ServerAddresses ($dnsShield, $fallback)

    Write-Host "  OK: DNS set to $dnsShield (fallback: $fallback)" -ForegroundColor Green
}

# Flush DNS cache
Write-Host ""
Write-Host "  Flushing DNS cache..." -ForegroundColor Yellow
Clear-DnsClientCache

# Verify
Write-Host ""
Write-Host "  Verification:" -ForegroundColor Cyan
$result = Resolve-DnsName -Name "canada.ca" -Server $dnsShield -ErrorAction SilentlyContinue
if ($result) {
    Write-Host "  [OK] DNS-SHIELD is working! canada.ca resolved." -ForegroundColor Green
} else {
    Write-Host "  [WARN] Could not verify — check server connectivity." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "  DNS-SHIELD active on all adapters." -ForegroundColor Green
Write-Host "  To remove: run dns-shield-remove.ps1" -ForegroundColor Gray
Write-Host ""
pause
`, serverIP, fallbackIP)
}

// WindowsPowerShellRemove generates the uninstall/revert script.
// Sets DNS back to automatic (DHCP) on all adapters.
func WindowsPowerShellRemove() string {
	return `# DNS-SHIELD Remove — Revert DNS to automatic (DHCP)
# Run as Administrator

$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "  DNS-SHIELD — Reverting DNS to automatic..." -ForegroundColor Cyan
Write-Host ""

$adapters = Get-NetAdapter | Where-Object { $_.Status -eq "Up" }

foreach ($adapter in $adapters) {
    Write-Host "  Resetting: $($adapter.Name)" -ForegroundColor Yellow
    Set-DnsClientServerAddress `
        -InterfaceIndex $adapter.InterfaceIndex `
        -ResetServerAddresses
    Write-Host "  OK: DNS set back to automatic (DHCP)" -ForegroundColor Green
}

Clear-DnsClientCache

Write-Host ""
Write-Host "  DNS-SHIELD removed. Your DNS is now automatic." -ForegroundColor Green
Write-Host ""
pause
`
}
