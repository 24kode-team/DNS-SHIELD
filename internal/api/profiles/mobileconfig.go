package profiles

import (
	"fmt"
	"time"
)

// MobileConfig generates an Apple .mobileconfig profile for iOS and macOS.
// When user opens this file on iPhone/iPad/Mac → Settings asks to install → DNS set automatically.
func MobileConfig(_, serverIP, dohURL string) string {
	uuid1 := generateUUID()
	uuid2 := generateUUID()
	uuid3 := generateUUID()
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>

  <key>PayloadDisplayName</key>
  <string>DNS-SHIELD — Safe Canadian DNS</string>

  <key>PayloadDescription</key>
  <string>Blocks phishing, malware, scams, adult content, gambling and sextortion.
No tracking. No data sold. Canada-hosted. Remove anytime.</string>

  <key>PayloadIdentifier</key>
  <string>ca.dns-shield.dns.profile</string>

  <key>PayloadUUID</key>
  <string>%s</string>

  <key>PayloadType</key>
  <string>Configuration</string>

  <key>PayloadVersion</key>
  <integer>1</integer>

  <key>PayloadOrganization</key>
  <string>DNS-SHIELD Canada</string>

  <key>PayloadRemovalDisallowed</key>
  <false/>

  <key>ConsentText</key>
  <dict>
    <key>en</key>
    <string>This profile configures your device to use DNS-SHIELD,
a privacy-respecting DNS filter hosted in Canada.

What it blocks: phishing, malware, scams, adult content, gambling, predatory sites, deepfakes, sextortion.

What it does NOT do: track you, log your queries, or sell your data.

Remove anytime: Settings → General → VPN &amp; Device Management → DNS-SHIELD → Remove.</string>
  </dict>

  <key>PayloadContent</key>
  <array>
    <dict>
      <key>PayloadType</key>
      <string>com.apple.dnsSettings.managed</string>

      <key>PayloadIdentifier</key>
      <string>ca.dns-shield.dns.settings</string>

      <key>PayloadUUID</key>
      <string>%s</string>

      <key>PayloadVersion</key>
      <integer>1</integer>

      <key>PayloadDisplayName</key>
      <string>DNS-SHIELD DNS Settings</string>

      <key>DNSSettings</key>
      <dict>
        <key>DNSProtocol</key>
        <string>HTTPS</string>

        <key>ServerURL</key>
        <string>%s</string>

        <key>ServerAddresses</key>
        <array>
          <string>%s</string>
        </array>
      </dict>

      <!-- Apply to ALL networks: Wi-Fi + Cellular -->
      <key>OnDemandRules</key>
      <array>
        <dict>
          <key>Action</key>
          <string>Connect</string>
        </dict>
      </array>

    </dict>

    <dict>
      <key>PayloadType</key>
      <string>com.apple.configurationprofile.identification</string>
      <key>PayloadIdentifier</key>
      <string>ca.dns-shield.dns.id</string>
      <key>PayloadUUID</key>
      <string>%s</string>
      <key>PayloadVersion</key>
      <integer>1</integer>
      <key>PayloadDisplayName</key>
      <string>DNS-SHIELD Profile Info</string>
      <key>InstallDate</key>
      <string>%s</string>
    </dict>

  </array>
</dict>
</plist>`, uuid1, uuid2, dohURL, serverIP, uuid3, now)
}

// MobileConfigDoT generates a DoT (DNS-over-TLS) variant for iOS/macOS.
func MobileConfigDoT(serverHost, serverIP string) string {
	uuid1 := generateUUID()
	uuid2 := generateUUID()

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>PayloadDisplayName</key>
  <string>DNS-SHIELD — Safe Canadian DNS (DoT)</string>

  <key>PayloadDescription</key>
  <string>Blocks phishing, malware, scams, adult content, gambling and sextortion.
No tracking. No data sold. Canada-hosted. Remove anytime.</string>

  <key>PayloadIdentifier</key>
  <string>ca.dns-shield.dns.dot.profile</string>

  <key>PayloadUUID</key>
  <string>%s</string>

  <key>PayloadType</key>
  <string>Configuration</string>

  <key>PayloadVersion</key>
  <integer>1</integer>

  <key>PayloadOrganization</key>
  <string>DNS-SHIELD Canada</string>

  <key>PayloadRemovalDisallowed</key>
  <false/>

  <key>PayloadContent</key>
  <array>
    <dict>
      <key>PayloadType</key>
      <string>com.apple.dnsSettings.managed</string>

      <key>PayloadIdentifier</key>
      <string>ca.dns-shield.dns.dot.settings</string>

      <key>PayloadUUID</key>
      <string>%s</string>

      <key>PayloadVersion</key>
      <integer>1</integer>

      <key>PayloadDisplayName</key>
      <string>DNS-SHIELD DoT Settings</string>

      <key>DNSSettings</key>
      <dict>
        <key>DNSProtocol</key>
        <string>TLS</string>

        <key>ServerName</key>
        <string>%s</string>

        <key>ServerAddresses</key>
        <array>
          <string>%s</string>
        </array>
      </dict>

      <key>OnDemandRules</key>
      <array>
        <dict>
          <key>Action</key>
          <string>Connect</string>
        </dict>
      </array>
    </dict>
  </array>
</dict>
</plist>`, uuid1, uuid2, serverHost, serverIP)
}

// generateUUID makes a random UUID v4 without external dependencies.
func generateUUID() string {
	b := make([]byte, 16)
	t := time.Now().UnixNano()
	for i := range b {
		b[i] = byte(t >> (uint(i) * 4))
		t = t*6364136223846793005 + 1442695040888963407
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
