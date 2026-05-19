package api

import "github.com/gofiber/fiber/v2"

// handleDashboard serves the built-in web dashboard.
// No external CDN — fully self-contained HTML/CSS/JS.
func (s *Server) handleDashboard(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(dashboardHTML)
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>DNS-SHIELD Dashboard</title>
<style>
  :root {
    --bg: #0a0e1a;
    --surface: #111827;
    --surface2: #1a2235;
    --border: #1e2d45;
    --accent: #00d4ff;
    --accent2: #0099cc;
    --green: #00ff88;
    --red: #ff4757;
    --yellow: #ffa502;
    --text: #e2e8f0;
    --muted: #64748b;
    --font: 'SF Mono', 'Fira Code', 'Consolas', monospace;
  }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    background: var(--bg);
    color: var(--text);
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    min-height: 100vh;
  }

  /* ── Header ── */
  header {
    background: var(--surface);
    border-bottom: 1px solid var(--border);
    padding: 0 2rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 60px;
    position: sticky;
    top: 0;
    z-index: 100;
  }
  .logo {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-size: 1.1rem;
    font-weight: 700;
    letter-spacing: -0.02em;
  }
  .logo-icon {
    width: 32px; height: 32px;
    background: linear-gradient(135deg, var(--accent), var(--accent2));
    border-radius: 8px;
    display: flex; align-items: center; justify-content: center;
    font-size: 18px;
  }
  .status-dot {
    width: 8px; height: 8px;
    border-radius: 50%;
    background: var(--green);
    box-shadow: 0 0 8px var(--green);
    animation: pulse 2s infinite;
  }
  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.4; }
  }
  .header-right {
    display: flex;
    align-items: center;
    gap: 1rem;
    font-size: 0.8rem;
    color: var(--muted);
  }
  .refresh-btn {
    background: var(--surface2);
    border: 1px solid var(--border);
    color: var(--text);
    padding: 0.35rem 0.9rem;
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.8rem;
    transition: all 0.15s;
  }
  .refresh-btn:hover { border-color: var(--accent); color: var(--accent); }

  /* ── Layout ── */
  main { padding: 2rem; max-width: 1200px; margin: 0 auto; }

  /* ── Stat cards ── */
  .cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
    margin-bottom: 2rem;
  }
  .card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 1.25rem 1.5rem;
    position: relative;
    overflow: hidden;
    transition: border-color 0.2s;
  }
  .card:hover { border-color: var(--accent); }
  .card::before {
    content: '';
    position: absolute;
    top: 0; left: 0; right: 0;
    height: 2px;
    background: var(--accent-color, var(--accent));
  }
  .card-label {
    font-size: 0.72rem;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: var(--muted);
    margin-bottom: 0.5rem;
  }
  .card-value {
    font-size: 2rem;
    font-weight: 700;
    font-family: var(--font);
    color: var(--card-color, var(--text));
    line-height: 1;
  }
  .card-sub {
    font-size: 0.75rem;
    color: var(--muted);
    margin-top: 0.4rem;
  }

  /* ── Two column grid ── */
  .grid2 {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    margin-bottom: 1.5rem;
  }
  @media (max-width: 768px) { .grid2 { grid-template-columns: 1fr; } }

  /* ── Panel ── */
  .panel {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 12px;
    overflow: hidden;
  }
  .panel-header {
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--border);
    font-size: 0.85rem;
    font-weight: 600;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .panel-body { padding: 1.25rem 1.5rem; }

  /* ── Category bars ── */
  .cat-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.75rem;
  }
  .cat-label {
    width: 90px;
    font-size: 0.78rem;
    color: var(--muted);
    text-transform: capitalize;
    flex-shrink: 0;
  }
  .cat-bar-wrap {
    flex: 1;
    background: var(--surface2);
    border-radius: 4px;
    height: 8px;
    overflow: hidden;
  }
  .cat-bar {
    height: 100%;
    border-radius: 4px;
    background: linear-gradient(90deg, var(--accent2), var(--accent));
    transition: width 0.6s ease;
  }
  .cat-count {
    font-family: var(--font);
    font-size: 0.78rem;
    color: var(--accent);
    width: 60px;
    text-align: right;
  }

  /* ── Latency bars ── */
  .lat-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.6rem;
    font-size: 0.8rem;
  }
  .lat-label { width: 70px; color: var(--muted); }
  .lat-bar-wrap { flex: 1; background: var(--surface2); border-radius: 4px; height: 6px; }
  .lat-bar { height: 100%; border-radius: 4px; background: var(--green); transition: width 0.6s; }
  .lat-count { width: 50px; text-align: right; color: var(--muted); font-family: var(--font); font-size: 0.75rem; }

  /* ── Block rate donut ── */
  .donut-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1.5rem;
    gap: 2rem;
  }
  .donut-svg { flex-shrink: 0; }
  .donut-legend { display: flex; flex-direction: column; gap: 0.5rem; }
  .legend-item { display: flex; align-items: center; gap: 0.5rem; font-size: 0.8rem; }
  .legend-dot { width: 10px; height: 10px; border-radius: 50%; }

  /* ── DoH Config box ── */
  .config-box {
    background: var(--surface2);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 1rem 1.25rem;
    margin-bottom: 1rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
  }
  .config-box-label { font-size: 0.75rem; color: var(--muted); margin-bottom: 0.25rem; }
  .config-box-value {
    font-family: var(--font);
    font-size: 0.85rem;
    color: var(--accent);
  }
  .copy-btn {
    background: transparent;
    border: 1px solid var(--border);
    color: var(--muted);
    padding: 0.3rem 0.7rem;
    border-radius: 5px;
    cursor: pointer;
    font-size: 0.75rem;
    flex-shrink: 0;
    transition: all 0.15s;
  }
  .copy-btn:hover { color: var(--accent); border-color: var(--accent); }

  /* ── OS Instructions ── */
  .tabs { display: flex; gap: 0.5rem; margin-bottom: 1rem; }
  .tab {
    padding: 0.4rem 1rem;
    border-radius: 6px;
    font-size: 0.8rem;
    cursor: pointer;
    border: 1px solid var(--border);
    background: transparent;
    color: var(--muted);
    transition: all 0.15s;
  }
  .tab.active { background: var(--accent); color: #000; border-color: var(--accent); font-weight: 600; }
  .instructions { display: none; }
  .instructions.active { display: block; }
  pre {
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 8px;
    padding: 1rem;
    font-family: var(--font);
    font-size: 0.8rem;
    color: var(--accent);
    overflow-x: auto;
    line-height: 1.6;
  }
  .comment { color: var(--muted); }

  /* ── Loading ── */
  .loading {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 120px;
    color: var(--muted);
    font-size: 0.85rem;
    gap: 0.5rem;
  }
  .spinner {
    width: 16px; height: 16px;
    border: 2px solid var(--border);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 0.7s linear infinite;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

  .error-msg { color: var(--red); font-size: 0.82rem; padding: 1rem; text-align: center; }
  .last-updated { font-size: 0.72rem; color: var(--muted); }
</style>
</head>
<body>

<header>
  <div class="logo">
    <div class="logo-icon">🛡️</div>
    DNS-SHIELD
  </div>
  <div class="header-right">
    <div class="status-dot" id="statusDot"></div>
    <span id="statusText">Connecting...</span>
    <button class="refresh-btn" onclick="loadAll()">↺ Refresh</button>
  </div>
</header>

<main>

  <!-- ── Stat Cards ── -->
  <div class="cards">
    <div class="card" style="--accent-color: var(--accent)">
      <div class="card-label">Total Queries</div>
      <div class="card-value" id="totalQueries">—</div>
      <div class="card-sub">since last restart</div>
    </div>
    <div class="card" style="--accent-color: var(--red); --card-color: var(--red)">
      <div class="card-label">Blocked</div>
      <div class="card-value" id="blockedQueries">—</div>
      <div class="card-sub" id="blockedPct">— % of total</div>
    </div>
    <div class="card" style="--accent-color: var(--green); --card-color: var(--green)">
      <div class="card-label">Allowed</div>
      <div class="card-value" id="allowedQueries">—</div>
      <div class="card-sub">clean traffic</div>
    </div>
    <div class="card" style="--accent-color: var(--yellow)">
      <div class="card-label">Blocklist Entries</div>
      <div class="card-value" id="totalEntries">—</div>
      <div class="card-sub">domains indexed</div>
    </div>
  </div>

  <!-- ── Categories + Latency ── -->
  <div class="grid2">
    <div class="panel">
      <div class="panel-header">
        <span>Blocks by Category</span>
        <span class="last-updated" id="lastUpdated"></span>
      </div>
      <div class="panel-body" id="categoryBars">
        <div class="loading"><div class="spinner"></div> Loading...</div>
      </div>
    </div>

    <div class="panel">
      <div class="panel-header">Response Latency</div>
      <div class="panel-body" id="latencyBars">
        <div class="loading"><div class="spinner"></div> Loading...</div>
      </div>
    </div>
  </div>

  <!-- ── Setup Instructions ── -->
  <div class="panel" style="margin-bottom:1.5rem">
    <div class="panel-header">Quick Setup — Point your DNS here</div>
    <div class="panel-body">

      <div class="config-box">
        <div>
          <div class="config-box-label">DNS Server (UDP/TCP)</div>
          <div class="config-box-value" id="dnsAddr">Loading...</div>
        </div>
        <button class="copy-btn" onclick="copyText(this, document.getElementById('dnsAddr').textContent)">Copy</button>
      </div>

      <div class="config-box">
        <div>
          <div class="config-box-label">DNS-over-HTTPS (DoH)</div>
          <div class="config-box-value" id="dohAddr">Loading...</div>
        </div>
        <button class="copy-btn" onclick="copyText(this, document.getElementById('dohAddr').textContent)">Copy</button>
      </div>

      <div class="config-box">
        <div>
          <div class="config-box-label">DNS-over-TLS (DoT)</div>
          <div class="config-box-value" id="dotAddr">Loading...</div>
        </div>
        <button class="copy-btn" onclick="copyText(this, document.getElementById('dotAddr').textContent)">Copy</button>
      </div>

      <div class="tabs">
        <button class="tab active" onclick="showTab('mac')">macOS</button>
        <button class="tab" onclick="showTab('windows')">Windows</button>
        <button class="tab" onclick="showTab('linux')">Linux</button>
        <button class="tab" onclick="showTab('router')">Router</button>
        <button class="tab" onclick="showTab('android')">Android</button>
        <button class="tab" onclick="showTab('ios')">iOS</button>
      </div>

      <div class="instructions active" id="tab-mac">
<pre><span class="comment"># macOS — System Settings → Wi-Fi → Details → DNS</span>
<span class="comment"># Or via terminal (replace en0 with your interface):</span>

sudo networksetup -setdnsservers Wi-Fi <span id="mac-dns">YOUR_SERVER_IP</span>

<span class="comment"># Verify:</span>
scutil --dns | grep nameserver</pre>
      </div>

      <div class="instructions" id="tab-windows">
<pre><span class="comment"># Windows — Settings → Network → DNS (Manual)</span>
<span class="comment"># Or via PowerShell (Admin):</span>

Set-DnsClientServerAddress -InterfaceAlias "Wi-Fi" -ServerAddresses ("<span id="win-dns">YOUR_SERVER_IP</span>")

<span class="comment"># Verify:</span>
nslookup canada.ca <span id="win-dns2">YOUR_SERVER_IP</span></pre>
      </div>

      <div class="instructions" id="tab-linux">
<pre><span class="comment"># Linux — systemd-resolved</span>

sudo systemctl stop systemd-resolved
echo "nameserver <span id="lin-dns">YOUR_SERVER_IP</span>" | sudo tee /etc/resolv.conf

<span class="comment"># Or NetworkManager:</span>
nmcli con mod "Your-Connection" ipv4.dns "<span id="lin-dns2">YOUR_SERVER_IP</span>"
nmcli con up "Your-Connection"</pre>
      </div>

      <div class="instructions" id="tab-router">
<pre><span class="comment"># Router — protects ALL devices on your network</span>
<span class="comment"># Log into your router admin panel (usually 192.168.1.1)</span>
<span class="comment"># Find: DHCP Settings → Primary DNS Server</span>

Primary DNS:   <span id="router-dns">YOUR_SERVER_IP</span>
Secondary DNS: 9.9.9.9   <span class="comment"># Quad9 fallback</span>

<span class="comment"># Save and restart router</span></pre>
      </div>

      <div class="instructions" id="tab-android">
<pre><span class="comment"># Android 9+ — Private DNS (DoT)</span>
<span class="comment"># Settings → Network → Private DNS → Custom hostname</span>

Hostname: <span id="android-dot">YOUR_SERVER_HOSTNAME</span>

<span class="comment"># Or: Settings → Wi-Fi → Modify → DNS</span>
DNS 1: <span id="android-dns">YOUR_SERVER_IP</span></pre>
      </div>

      <div class="instructions" id="tab-ios">
<pre><span class="comment"># iOS — Settings → Wi-Fi → (i) → Configure DNS</span>
<span class="comment"># Change to Manual, add:</span>

DNS Server: <span id="ios-dns">YOUR_SERVER_IP</span>

<span class="comment"># For DoH on iOS 14+, install a .mobileconfig profile</span>
<span class="comment"># (see docs for profile generator)</span></pre>
      </div>

    </div>
  </div>

  <!-- ── Blocklist Entries per Category ── -->
  <div class="panel">
    <div class="panel-header">Blocklist Feed Stats</div>
    <div class="panel-body" id="feedStats">
      <div class="loading"><div class="spinner"></div> Loading...</div>
    </div>
  </div>

</main>

<script>
  const HOST = window.location.hostname;
  const PORT = window.location.port;
  const BASE = window.location.origin;

  // ── Populate addresses ─────────────────────────────────────────────────────
  document.getElementById('dnsAddr').textContent = HOST + ':53';
  document.getElementById('dohAddr').textContent = BASE + '/dns-query';
  document.getElementById('dotAddr').textContent = HOST + ':853';

  ['mac-dns','win-dns','win-dns2','lin-dns','lin-dns2','router-dns','android-dns','ios-dns'].forEach(id => {
    const el = document.getElementById(id);
    if (el) el.textContent = HOST;
  });
  const dotEl = document.getElementById('android-dot');
  if (dotEl) dotEl.textContent = HOST;

  // ── Fetch ──────────────────────────────────────────────────────────────────
  async function fetchJSON(url) {
    const r = await fetch(url);
    if (!r.ok) throw new Error(r.status);
    return r.json();
  }

  function fmt(n) {
    if (n === undefined || n === null) return '—';
    if (n >= 1_000_000) return (n/1_000_000).toFixed(1) + 'M';
    if (n >= 1_000) return (n/1_000).toFixed(1) + 'K';
    return n.toString();
  }

  // ── Render metrics ─────────────────────────────────────────────────────────
  function renderMetrics(data) {
    document.getElementById('totalQueries').textContent = fmt(data.total_queries);
    document.getElementById('blockedQueries').textContent = fmt(data.blocked_queries);
    document.getElementById('allowedQueries').textContent = fmt(data.allowed_queries);

    const pct = data.total_queries > 0
      ? ((data.blocked_queries / data.total_queries) * 100).toFixed(1)
      : 0;
    document.getElementById('blockedPct').textContent = pct + '% of total';

    // Category bars
    const cats = data.blocked_by_category || {};
    const maxCat = Math.max(1, ...Object.values(cats));
    const catOrder = ['phishing','malware','scam','porn','gambling','predatory','deepfake','sextortion'];
    const catHtml = catOrder.map(cat => {
      const count = cats[cat] || 0;
      const pct = ((count / maxCat) * 100).toFixed(1);
      return `<div class="cat-row">
        <div class="cat-label">${cat}</div>
        <div class="cat-bar-wrap"><div class="cat-bar" style="width:${pct}%"></div></div>
        <div class="cat-count">${fmt(count)}</div>
      </div>`;
    }).join('');
    document.getElementById('categoryBars').innerHTML = catHtml || '<div class="loading">No blocks yet</div>';

    // Latency bars
    const lat = data.latency_ms || {};
    const latData = [
      { label: '<1ms',  val: lat.under_1ms  || 0 },
      { label: '<5ms',  val: lat.under_5ms  || 0 },
      { label: '<10ms', val: lat.under_10ms || 0 },
      { label: '<50ms', val: lat.under_50ms || 0 },
      { label: '50ms+', val: lat.over_50ms  || 0 },
    ];
    const maxLat = Math.max(1, ...latData.map(d => d.val));
    const latHtml = latData.map(d => {
      const pct = ((d.val / maxLat) * 100).toFixed(1);
      return `<div class="lat-row">
        <div class="lat-label">${d.label}</div>
        <div class="lat-bar-wrap"><div class="lat-bar" style="width:${pct}%"></div></div>
        <div class="lat-count">${fmt(d.val)}</div>
      </div>`;
    }).join('');
    document.getElementById('latencyBars').innerHTML = latHtml;

    document.getElementById('lastUpdated').textContent =
      'Updated ' + new Date().toLocaleTimeString();
  }

  // ── Render feed stats ──────────────────────────────────────────────────────
  function renderFeedStats(data) {
    const total = Object.values(data).reduce((a, b) => a + b, 0);
    document.getElementById('totalEntries').textContent = fmt(total);

    const sorted = Object.entries(data).sort((a,b) => b[1]-a[1]);
    const maxVal = Math.max(1, ...sorted.map(d => d[1]));
    const html = sorted.map(([cat, count]) => {
      const pct = ((count / maxVal) * 100).toFixed(1);
      return `<div class="cat-row">
        <div class="cat-label">${cat}</div>
        <div class="cat-bar-wrap"><div class="cat-bar" style="width:${pct}%"></div></div>
        <div class="cat-count">${fmt(count)} entries</div>
      </div>`;
    }).join('');
    document.getElementById('feedStats').innerHTML = html || '<div class="loading">No data</div>';
  }

  // ── Health ─────────────────────────────────────────────────────────────────
  function setStatus(ok) {
    const dot = document.getElementById('statusDot');
    const txt = document.getElementById('statusText');
    dot.style.background = ok ? 'var(--green)' : 'var(--red)';
    dot.style.boxShadow = ok ? '0 0 8px var(--green)' : '0 0 8px var(--red)';
    txt.textContent = ok ? 'Online' : 'Unreachable';
  }

  // ── Load all data ──────────────────────────────────────────────────────────
  async function loadAll() {
    try {
      const [metrics, feedStats] = await Promise.all([
        fetchJSON('/metrics'),
        fetchJSON('/admin/blocklist/stats').catch(() => ({})),
      ]);
      renderMetrics(metrics);
      renderFeedStats(feedStats);
      setStatus(true);
    } catch (e) {
      setStatus(false);
      document.getElementById('categoryBars').innerHTML =
        '<div class="error-msg">Could not reach DNS-SHIELD API</div>';
    }
  }

  // ── Tab switching ──────────────────────────────────────────────────────────
  function showTab(id) {
    document.querySelectorAll('.instructions').forEach(el => el.classList.remove('active'));
    document.querySelectorAll('.tab').forEach(el => el.classList.remove('active'));
    document.getElementById('tab-' + id).classList.add('active');
    event.target.classList.add('active');
  }

  // ── Copy helper ────────────────────────────────────────────────────────────
  function copyText(btn, text) {
    navigator.clipboard.writeText(text).then(() => {
      const orig = btn.textContent;
      btn.textContent = '✓ Copied';
      btn.style.color = 'var(--green)';
      setTimeout(() => { btn.textContent = orig; btn.style.color = ''; }, 1500);
    });
  }

  // ── Auto-refresh every 10s ─────────────────────────────────────────────────
  loadAll();
  setInterval(loadAll, 10000);
</script>
</body>
</html>`
