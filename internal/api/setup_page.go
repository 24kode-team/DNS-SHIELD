package api

import "fmt"

func setupPageHTML(host, ip, doh, baseURL string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>DNS-SHIELD — Setup</title>
<style>
  :root {
    --bg: #0a0e1a;
    --surface: #111827;
    --surface2: #1a2235;
    --border: #1e2d45;
    --accent: #00d4ff;
    --green: #00ff88;
    --red: #ff4757;
    --yellow: #ffa502;
    --text: #e2e8f0;
    --muted: #64748b;
  }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    background: var(--bg);
    color: var(--text);
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    min-height: 100vh;
  }

  header {
    background: var(--surface);
    border-bottom: 1px solid var(--border);
    padding: 1rem 2rem;
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  .logo { font-size: 1.2rem; font-weight: 700; }
  .back { margin-left: auto; font-size: 0.82rem; color: var(--muted); text-decoration: none; }
  .back:hover { color: var(--accent); }

  main { max-width: 760px; margin: 0 auto; padding: 2.5rem 1.5rem; }

  h1 { font-size: 1.6rem; font-weight: 700; margin-bottom: 0.5rem; }
  .subtitle { color: var(--muted); font-size: 0.9rem; margin-bottom: 2.5rem; }

  /* Platform cards */
  .platforms {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
    gap: 0.75rem;
    margin-bottom: 2.5rem;
  }
  .platform-btn {
    background: var(--surface);
    border: 2px solid var(--border);
    border-radius: 12px;
    padding: 1.25rem 0.75rem;
    cursor: pointer;
    text-align: center;
    transition: all 0.15s;
    color: var(--text);
  }
  .platform-btn:hover, .platform-btn.active {
    border-color: var(--accent);
    background: var(--surface2);
  }
  .platform-btn.active { border-color: var(--accent); }
  .platform-icon { font-size: 2rem; margin-bottom: 0.5rem; display: block; }
  .platform-name { font-size: 0.82rem; font-weight: 600; }

  /* Steps panel */
  .steps-panel {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 16px;
    overflow: hidden;
  }
  .steps-header {
    padding: 1.25rem 1.5rem;
    border-bottom: 1px solid var(--border);
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  .steps-header .icon { font-size: 1.5rem; }
  .steps-header .title { font-size: 1rem; font-weight: 700; }
  .steps-header .desc { font-size: 0.8rem; color: var(--muted); margin-top: 0.15rem; }

  .steps { padding: 1.5rem; }

  .step {
    display: flex;
    gap: 1rem;
    margin-bottom: 1.5rem;
    align-items: flex-start;
  }
  .step:last-child { margin-bottom: 0; }
  .step-num {
    width: 28px; height: 28px;
    border-radius: 50%;
    background: var(--accent);
    color: #000;
    font-weight: 700;
    font-size: 0.8rem;
    display: flex; align-items: center; justify-content: center;
    flex-shrink: 0;
    margin-top: 2px;
  }
  .step-content { flex: 1; }
  .step-title { font-size: 0.9rem; font-weight: 600; margin-bottom: 0.4rem; }
  .step-desc { font-size: 0.82rem; color: var(--muted); line-height: 1.5; }

  /* Download button */
  .dl-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--accent);
    color: #000;
    font-weight: 700;
    font-size: 0.9rem;
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    text-decoration: none;
    margin-top: 0.75rem;
    transition: opacity 0.15s;
  }
  .dl-btn:hover { opacity: 0.85; }
  .dl-btn-secondary {
    background: var(--surface2);
    border: 1px solid var(--border);
    color: var(--muted);
  }
  .dl-btn-secondary:hover { border-color: var(--accent); color: var(--accent); opacity: 1; }

  /* Revert section */
  .revert {
    margin-top: 1.5rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border);
  }
  .revert-title {
    font-size: 0.78rem;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--muted);
    margin-bottom: 0.75rem;
  }

  /* Copy field */
  .copy-field {
    display: flex;
    align-items: center;
    background: var(--surface2);
    border: 1px solid var(--border);
    border-radius: 8px;
    overflow: hidden;
    margin-top: 0.75rem;
  }
  .copy-field-value {
    flex: 1;
    padding: 0.65rem 1rem;
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 0.9rem;
    color: var(--accent);
    word-break: break-all;
  }
  .copy-field-btn {
    padding: 0.65rem 1rem;
    background: transparent;
    border: none;
    border-left: 1px solid var(--border);
    color: var(--muted);
    cursor: pointer;
    font-size: 0.8rem;
    white-space: nowrap;
    transition: color 0.15s;
  }
  .copy-field-btn:hover { color: var(--accent); }

  /* Notice */
  .notice {
    background: rgba(0, 212, 255, 0.06);
    border: 1px solid rgba(0, 212, 255, 0.2);
    border-radius: 8px;
    padding: 0.75rem 1rem;
    font-size: 0.8rem;
    color: var(--muted);
    margin-top: 1rem;
    line-height: 1.5;
  }
  .notice strong { color: var(--text); }

  .platform-section { display: none; }
  .platform-section.active { display: block; }
</style>
</head>
<body>

<header>
  <span style="font-size:1.4rem">🛡️</span>
  <span class="logo">DNS-SHIELD</span>
  <a class="back" href="/">← Dashboard</a>
</header>

<main>
  <h1>Connect Your Device</h1>
  <p class="subtitle">Choose your platform — download and run. No App Store. No account. Remove anytime.</p>

  <!-- Platform selector -->
  <div class="platforms">
    <div class="platform-btn active" onclick="showPlatform('ios', this)">
      <span class="platform-icon">📱</span>
      <div class="platform-name">iPhone / iPad</div>
    </div>
    <div class="platform-btn" onclick="showPlatform('macos', this)">
      <span class="platform-icon">🖥️</span>
      <div class="platform-name">Mac</div>
    </div>
    <div class="platform-btn" onclick="showPlatform('windows', this)">
      <span class="platform-icon">🪟</span>
      <div class="platform-name">Windows</div>
    </div>
    <div class="platform-btn" onclick="showPlatform('android', this)">
      <span class="platform-icon">🤖</span>
      <div class="platform-name">Android</div>
    </div>
    <div class="platform-btn" onclick="showPlatform('linux', this)">
      <span class="platform-icon">🐧</span>
      <div class="platform-name">Linux</div>
    </div>
    <div class="platform-btn" onclick="showPlatform('router', this)">
      <span class="platform-icon">📡</span>
      <div class="platform-name">Router</div>
    </div>
  </div>

  <div class="steps-panel">

    <!-- ── iOS ── -->
    <div class="platform-section active" id="section-ios">
      <div class="steps-header">
        <span class="icon">📱</span>
        <div>
          <div class="title">iPhone & iPad</div>
          <div class="desc">One-tap install via Apple Configuration Profile</div>
        </div>
      </div>
      <div class="steps">
        <div class="step">
          <div class="step-num">1</div>
          <div class="step-content">
            <div class="step-title">Download the profile</div>
            <div class="step-desc">Tap the button below on your iPhone or iPad.</div>
            <a class="dl-btn" href="/setup/ios">⬇ Download for iPhone / iPad</a>
          </div>
        </div>
        <div class="step">
          <div class="step-num">2</div>
          <div class="step-content">
            <div class="step-title">Install the profile</div>
            <div class="step-desc">
              iOS will show a prompt — tap <strong>Allow</strong>, then go to:<br>
              <strong>Settings → General → VPN &amp; Device Management → DNS-SHIELD → Install</strong>
            </div>
          </div>
        </div>
        <div class="step">
          <div class="step-num">3</div>
          <div class="step-content">
            <div class="step-title">Done ✓</div>
            <div class="step-desc">DNS-SHIELD is now active on Wi-Fi and cellular.</div>
          </div>
        </div>

        <div class="revert">
          <div class="revert-title">↩ Remove / Revert to Default</div>
          <div class="step-desc">
            Go to <strong>Settings → General → VPN &amp; Device Management → DNS-SHIELD → Remove Profile</strong>.<br>
            Your DNS will go back to automatic immediately.
          </div>
        </div>

        <div class="notice">
          <strong>Privacy:</strong> This profile only changes your DNS server. It cannot read your messages, access your camera, or track your location.
        </div>
      </div>
    </div>

    <!-- ── macOS ── -->
    <div class="platform-section" id="section-macos">
      <div class="steps-header">
        <span class="icon">🖥️</span>
        <div>
          <div class="title">Mac</div>
          <div class="desc">One-click install via Apple Configuration Profile</div>
        </div>
      </div>
      <div class="steps">
        <div class="step">
          <div class="step-num">1</div>
          <div class="step-content">
            <div class="step-title">Download the profile</div>
            <a class="dl-btn" href="/setup/macos">⬇ Download for Mac</a>
          </div>
        </div>
        <div class="step">
          <div class="step-num">2</div>
          <div class="step-content">
            <div class="step-title">Open the profile</div>
            <div class="step-desc">
              Double-click the downloaded file. macOS will open <strong>System Settings → Privacy &amp; Security → Profiles</strong>.
            </div>
          </div>
        </div>
        <div class="step">
          <div class="step-num">3</div>
          <div class="step-content">
            <div class="step-title">Click Install</div>
            <div class="step-desc">Enter your Mac password if prompted. Done ✓</div>
          </div>
        </div>

        <div class="revert">
          <div class="revert-title">↩ Remove / Revert to Default</div>
          <div class="step-desc">
            <strong>System Settings → Privacy &amp; Security → Profiles → DNS-SHIELD → Remove</strong>.<br>
            DNS goes back to automatic immediately.
          </div>
        </div>
      </div>
    </div>

    <!-- ── Windows ── -->
    <div class="platform-section" id="section-windows">
      <div class="steps-header">
        <span class="icon">🪟</span>
        <div>
          <div class="title">Windows</div>
          <div class="desc">PowerShell script — sets DNS on all network adapters</div>
        </div>
      </div>
      <div class="steps">
        <div class="step">
          <div class="step-num">1</div>
          <div class="step-content">
            <div class="step-title">Download the setup script</div>
            <a class="dl-btn" href="/setup/windows">⬇ Download dns-shield-setup.ps1</a>
          </div>
        </div>
        <div class="step">
          <div class="step-num">2</div>
          <div class="step-content">
            <div class="step-title">Run as Administrator</div>
            <div class="step-desc">
              Right-click the downloaded file → <strong>Run with PowerShell</strong>.<br>
              Click <strong>Yes</strong> on the UAC prompt. The script sets DNS and verifies automatically.
            </div>
          </div>
        </div>
        <div class="step">
          <div class="step-num">3</div>
          <div class="step-content">
            <div class="step-title">Done ✓</div>
            <div class="step-desc">All network adapters now use DNS-SHIELD.</div>
          </div>
        </div>

        <div class="revert">
          <div class="revert-title">↩ Remove / Revert to Default</div>
          <a class="dl-btn dl-btn-secondary" href="/setup/windows-remove">⬇ Download dns-shield-remove.ps1</a>
          <div class="step-desc" style="margin-top:0.5rem">
            Run it the same way (right-click → Run with PowerShell). DNS goes back to automatic.
          </div>
        </div>
      </div>
    </div>

    <!-- ── Android ── -->
    <div class="platform-section" id="section-android">
      <div class="steps-header">
        <span class="icon">🤖</span>
        <div>
          <div class="title">Android</div>
          <div class="desc">Private DNS (DoT) — built into Android 9 and later</div>
        </div>
      </div>
      <div class="steps">
        <div class="step">
          <div class="step-num">1</div>
          <div class="step-content">
            <div class="step-title">Copy this hostname</div>
            <div class="copy-field">
              <div class="copy-field-value" id="android-host">%s</div>
              <button class="copy-field-btn" onclick="copyField('android-host', this)">Copy</button>
            </div>
          </div>
        </div>
        <div class="step">
          <div class="step-num">2</div>
          <div class="step-content">
            <div class="step-title">Open Android Settings</div>
            <div class="step-desc">
              Go to <strong>Settings → Network &amp; Internet → Private DNS</strong>.<br>
              (On Samsung: <strong>Settings → Connections → More connection settings → Private DNS</strong>)
            </div>
          </div>
        </div>
        <div class="step">
          <div class="step-num">3</div>
          <div class="step-content">
            <div class="step-title">Select "Private DNS provider hostname"</div>
            <div class="step-desc">Paste the hostname and tap <strong>Save</strong>. Done ✓</div>
          </div>
        </div>

        <div class="revert">
          <div class="revert-title">↩ Remove / Revert to Default</div>
          <div class="step-desc">
            Go back to <strong>Settings → Network → Private DNS</strong> and select <strong>Automatic</strong>.
          </div>
        </div>
      </div>
    </div>

    <!-- ── Linux ── -->
    <div class="platform-section" id="section-linux">
      <div class="steps-header">
        <span class="icon">🐧</span>
        <div>
          <div class="title">Linux</div>
          <div class="desc">Shell script — supports systemd-resolved, NetworkManager, resolv.conf</div>
        </div>
      </div>
      <div class="steps">
        <div class="step">
          <div class="step-num">1</div>
          <div class="step-content">
            <div class="step-title">Download and run</div>
            <a class="dl-btn" href="/setup/linux">⬇ Download dns-shield-setup.sh</a>
            <div class="step-desc" style="margin-top:0.75rem">Then run:</div>
            <div class="copy-field" style="margin-top:0.4rem">
              <div class="copy-field-value">sudo bash dns-shield-setup.sh</div>
              <button class="copy-field-btn" onclick="copyText(this, 'sudo bash dns-shield-setup.sh')">Copy</button>
            </div>
          </div>
        </div>
        <div class="step">
          <div class="step-num">2</div>
          <div class="step-content">
            <div class="step-title">Done ✓</div>
            <div class="step-desc">Script auto-detects your DNS system and configures it.</div>
          </div>
        </div>

        <div class="revert">
          <div class="revert-title">↩ Remove / Revert to Default</div>
          <a class="dl-btn dl-btn-secondary" href="/setup/linux-remove">⬇ Download dns-shield-remove.sh</a>
          <div class="step-desc" style="margin-top:0.5rem">
            Run: <code>sudo bash dns-shield-remove.sh</code>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Router ── -->
    <div class="platform-section" id="section-router">
      <div class="steps-header">
        <span class="icon">📡</span>
        <div>
          <div class="title">Router</div>
          <div class="desc">Protects every device on your network — TV, consoles, everything</div>
        </div>
      </div>
      <div class="steps">
        <div class="step">
          <div class="step-num">1</div>
          <div class="step-content">
            <div class="step-title">Log into your router</div>
            <div class="step-desc">Usually <strong>192.168.1.1</strong> or <strong>192.168.0.1</strong> in your browser.</div>
          </div>
        </div>
        <div class="step">
          <div class="step-num">2</div>
          <div class="step-content">
            <div class="step-title">Find DNS settings</div>
            <div class="step-desc">Look for <strong>DHCP Settings</strong>, <strong>LAN Settings</strong>, or <strong>DNS Server</strong>.</div>
          </div>
        </div>
        <div class="step">
          <div class="step-num">3</div>
          <div class="step-content">
            <div class="step-title">Set DNS servers</div>
            <div class="step-desc">Primary DNS:</div>
            <div class="copy-field">
              <div class="copy-field-value" id="router-ip">%s</div>
              <button class="copy-field-btn" onclick="copyField('router-ip', this)">Copy</button>
            </div>
            <div class="step-desc" style="margin-top:0.5rem">Secondary DNS (fallback):</div>
            <div class="copy-field">
              <div class="copy-field-value">9.9.9.9</div>
              <button class="copy-field-btn" onclick="copyText(this, '9.9.9.9')">Copy</button>
            </div>
          </div>
        </div>
        <div class="step">
          <div class="step-num">4</div>
          <div class="step-content">
            <div class="step-title">Save and restart router</div>
            <div class="step-desc">All devices on the network are now protected. ✓</div>
          </div>
        </div>

        <div class="revert">
          <div class="revert-title">↩ Remove / Revert to Default</div>
          <div class="step-desc">
            Go back to router settings and set DNS to <strong>Automatic</strong> or delete the custom DNS entries.
          </div>
        </div>
      </div>
    </div>

  </div><!-- /steps-panel -->
</main>

<script>
  function showPlatform(id, btn) {
    document.querySelectorAll('.platform-section').forEach(s => s.classList.remove('active'));
    document.querySelectorAll('.platform-btn').forEach(b => b.classList.remove('active'));
    document.getElementById('section-' + id).classList.add('active');
    btn.classList.add('active');
  }

  function copyField(fieldId, btn) {
    const text = document.getElementById(fieldId).textContent.trim();
    copyText(btn, text);
  }

  function copyText(btn, text) {
    navigator.clipboard.writeText(text).then(() => {
      const orig = btn.textContent;
      btn.textContent = '✓ Copied';
      btn.style.color = 'var(--green)';
      setTimeout(() => { btn.textContent = orig; btn.style.color = ''; }, 1800);
    });
  }
</script>
</body>
</html>`, host, ip)
}
