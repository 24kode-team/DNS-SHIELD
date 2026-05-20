package profiles

import "fmt"

// WindowsPowerShell generates a .ps1 script that sets DNS on ALL active adapters.
func WindowsPowerShell(serverIP, fallbackIP string) string {
	return fmt.Sprintf("# DNS-SHIELD Setup Script for Windows\n"+
		"# Run as Administrator: Right-click -> Run with PowerShell\n"+
		"# To revert: run dns-shield-remove.ps1\n\n"+
		"$ErrorActionPreference = \"Stop\"\n"+
		"$dnsShield  = \"%s\"\n"+
		"$fallback   = \"%s\"\n\n"+
		"Write-Host \"\"\n"+
		"Write-Host \"  DNS-SHIELD - Setting up safe DNS...\" -ForegroundColor Cyan\n"+
		"Write-Host \"\"\n\n"+
		"$adapters = Get-NetAdapter | Where-Object { $_.Status -eq \"Up\" }\n\n"+
		"if ($adapters.Count -eq 0) {\n"+
		"    Write-Host \"  [ERROR] No active network adapters found.\" -ForegroundColor Red\n"+
		"    exit 1\n"+
		"}\n\n"+
		"foreach ($adapter in $adapters) {\n"+
		"    Write-Host \"  Configuring: $($adapter.Name)\" -ForegroundColor Yellow\n"+
		"    Set-DnsClientServerAddress -InterfaceIndex $adapter.InterfaceIndex -ServerAddresses ($dnsShield, $fallback)\n"+
		"    Write-Host \"  OK: DNS set to $dnsShield (fallback: $fallback)\" -ForegroundColor Green\n"+
		"}\n\n"+
		"Write-Host \"\"\n"+
		"Write-Host \"  Flushing DNS cache...\" -ForegroundColor Yellow\n"+
		"Clear-DnsClientCache\n\n"+
		"Write-Host \"\"\n"+
		"Write-Host \"  Verification:\" -ForegroundColor Cyan\n"+
		"$result = Resolve-DnsName -Name \"canada.ca\" -Server $dnsShield -ErrorAction SilentlyContinue\n"+
		"if ($result) {\n"+
		"    Write-Host \"  [OK] DNS-SHIELD is working!\" -ForegroundColor Green\n"+
		"} else {\n"+
		"    Write-Host \"  [WARN] Could not verify - check server connectivity.\" -ForegroundColor Yellow\n"+
		"}\n\n"+
		"Write-Host \"\"\n"+
		"Write-Host \"  DNS-SHIELD active on all adapters.\" -ForegroundColor Green\n"+
		"Write-Host \"  To remove: run dns-shield-remove.ps1\" -ForegroundColor Gray\n"+
		"Write-Host \"\"\n"+
		"pause\n",
		serverIP, fallbackIP)
}

// WindowsPowerShellRemove generates the uninstall/revert script.
func WindowsPowerShellRemove() string {
	return "# DNS-SHIELD Remove - Revert DNS to automatic (DHCP)\n" +
		"# Run as Administrator\n\n" +
		"$ErrorActionPreference = \"Stop\"\n\n" +
		"Write-Host \"\"\n" +
		"Write-Host \"  DNS-SHIELD - Reverting DNS to automatic...\" -ForegroundColor Cyan\n" +
		"Write-Host \"\"\n\n" +
		"$adapters = Get-NetAdapter | Where-Object { $_.Status -eq \"Up\" }\n\n" +
		"foreach ($adapter in $adapters) {\n" +
		"    Write-Host \"  Resetting: $($adapter.Name)\" -ForegroundColor Yellow\n" +
		"    Set-DnsClientServerAddress -InterfaceIndex $adapter.InterfaceIndex -ResetServerAddresses\n" +
		"    Write-Host \"  OK: DNS set back to automatic (DHCP)\" -ForegroundColor Green\n" +
		"}\n\n" +
		"Clear-DnsClientCache\n\n" +
		"Write-Host \"\"\n" +
		"Write-Host \"  DNS-SHIELD removed. Your DNS is now automatic.\" -ForegroundColor Green\n" +
		"Write-Host \"\"\n" +
		"pause\n"
}

// WindowsReg generates a .reg file that sets DNS for all network adapters.
func WindowsReg(serverIP string) string {
	return fmt.Sprintf("Windows Registry Editor Version 5.00\n\n"+
		"; DNS-SHIELD - Safe Canadian DNS\n"+
		"[HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters]\n"+
		"\"NameServer\"=\"%s\"\n"+
		"\"SearchList\"=\"\"\n",
		serverIP)
}
