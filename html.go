package main

// indexHTML is served at /
var indexHTML = []byte(`<!DOCTYPE html>
<html lang="fa" dir="rtl">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Napster Config Generator</title>
    <link href="https://fonts.googleapis.com/css2?family=Vazirmatn:wght@300;400;700&family=JetBrains+Mono&display=swap" rel="stylesheet">
    <style>
        :root {
            --primary: #3b82f6;
            --primary-hover: #2563eb;
            --bg-body: #0f172a;
            --bg-card: #1e293b;
            --bg-input: #0f172a;
            --text-main: #f8fafc;
            --text-dim: #94a3b8;
            --border: #334155;
            --radius: 12px;
            --success: #10b981;
            --error: #ef4444;
            --warning: #f59e0b;
        }
        .theme-emerald { --primary: #10b981; --primary-hover: #059669; }
        .theme-rose    { --primary: #f43f5e; --primary-hover: #e11d48; }
        .theme-amber   { --primary: #f59e0b; --primary-hover: #d97706; }
        .theme-violet  { --primary: #8b5cf6; --primary-hover: #7c3aed; }
        .light-mode {
            --bg-body: #f1f5f9; --bg-card: #ffffff;
            --bg-input: #f8fafc; --text-main: #1e293b;
            --text-dim: #64748b; --border: #e2e8f0;
        }
        *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: 'Vazirmatn', sans-serif; background: var(--bg-body); color: var(--text-main); line-height: 1.6; min-height: 100vh; transition: background .3s, color .3s; }
        ::-webkit-scrollbar { width: 5px; height: 5px; }
        ::-webkit-scrollbar-thumb { background: var(--border); border-radius: 10px; }
        .container { max-width: 1180px; margin: 0 auto; padding: 0 20px 60px; }
        header { display: flex; align-items: center; justify-content: space-between; padding: 18px 0; margin-bottom: 28px; border-bottom: 1px solid var(--border); position: sticky; top: 0; background: var(--bg-body); z-index: 50; }
        .logo { font-size: 22px; font-weight: 700; color: var(--primary); display: flex; align-items: center; gap: 8px; }
        .logo span { font-weight: 300; opacity: .7; }
        .hcontrols { display: flex; gap: 12px; align-items: center; }
        .color-picker { display: flex; gap: 6px; }
        .cdot { width: 16px; height: 16px; border-radius: 50%; cursor: pointer; border: 2px solid transparent; transition: transform .2s; }
        .cdot:hover { transform: scale(1.3); }
        .btn-icon { background: var(--bg-card); border: 1px solid var(--border); color: var(--text-main); border-radius: 8px; padding: 6px 10px; cursor: pointer; font-size: 15px; }
        .nav-tabs { display: flex; background: var(--bg-card); padding: 5px; border-radius: 14px; gap: 4px; margin-bottom: 28px; border: 1px solid var(--border); }
        .tab-btn { flex: 1; padding: 11px 8px; border: none; background: none; color: var(--text-dim); font-family: inherit; font-weight: 600; font-size: 14px; cursor: pointer; border-radius: 10px; transition: background .2s, color .2s; display: flex; align-items: center; justify-content: center; gap: 6px; }
        .tab-btn.active { background: var(--primary); color: #fff; }
        .tab-btn:not(.active):hover { background: var(--border); color: var(--text-main); }
        .card { background: var(--bg-card); border-radius: 16px; border: 1px solid var(--border); padding: 22px; margin-bottom: 20px; position: relative; overflow: hidden; }
        .card::before { content: ""; position: absolute; top: 0; right: 0; width: 3px; height: 100%; background: var(--primary); opacity: .6; }
        .card-title { font-weight: 700; font-size: 16px; margin-bottom: 18px; display: flex; align-items: center; gap: 8px; color: var(--primary); }
        .grid2 { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
        @media (max-width: 760px) { .grid2 { grid-template-columns: 1fr; } }
        .fld { margin-bottom: 14px; }
        .fld label { display: block; font-size: 12px; color: var(--text-dim); margin-bottom: 5px; font-weight: 500; }
        textarea, input[type=text], input[type=password], select { width: 100%; background: var(--bg-input); border: 1px solid var(--border); border-radius: 10px; padding: 10px 14px; color: var(--text-main); font-family: inherit; font-size: 13px; outline: none; transition: border-color .2s, box-shadow .2s; }
        textarea:focus, input:focus, select:focus { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(59,130,246,.12); }
        textarea { resize: vertical; min-height: 80px; }
        .row2 { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
        .btn { padding: 10px 20px; border-radius: 10px; border: none; font-family: inherit; font-weight: 700; font-size: 13px; cursor: pointer; display: inline-flex; align-items: center; gap: 7px; justify-content: center; transition: all .2s; }
        .btn:active { transform: scale(.97); }
        .btn-primary { background: var(--primary); color: #fff; }
        .btn-primary:hover { background: var(--primary-hover); transform: translateY(-1px); }
        .btn-outline { background: transparent; border: 1px solid var(--border); color: var(--text-dim); }
        .btn-outline:hover { border-color: var(--primary); color: var(--primary); }
        .btn-danger { background: rgba(239,68,68,.1); border: 1px solid rgba(239,68,68,.3); color: var(--error); }
        .btn-danger:hover { background: rgba(239,68,68,.2); }
        .btn-success { background: rgba(16,185,129,.15); border: 1px solid rgba(16,185,129,.3); color: var(--success); }
        .btn-success:hover { background: rgba(16,185,129,.25); }
        .btn-wide { width: 100%; padding: 14px; font-size: 15px; }
        .btn-row { display: flex; gap: 8px; flex-wrap: wrap; margin-top: 14px; }
        .sw-row { display: flex; align-items: center; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid var(--border); }
        .sw-row:last-child { border-bottom: none; }
        .sw-label { font-size: 13px; font-weight: 500; }
        .sw-desc { font-size: 11px; color: var(--text-dim); margin-top: 2px; }
        .switch { position: relative; width: 42px; height: 22px; flex-shrink: 0; }
        .switch input { opacity: 0; width: 0; height: 0; position: absolute; }
        .slider { position: absolute; inset: 0; background: var(--border); border-radius: 22px; cursor: pointer; transition: background .25s; }
        .slider::before { content: ''; position: absolute; width: 16px; height: 16px; left: 3px; top: 3px; background: #fff; border-radius: 50%; transition: transform .25s; }
        input:checked + .slider { background: var(--primary); }
        input:checked + .slider::before { transform: translateX(20px); }
        .sw-extra { margin-top: 10px; padding: 12px; background: var(--bg-input); border-radius: 8px; border: 1px solid var(--border); }
        .terminal { background: #020617; color: #4ade80; font-family: 'JetBrains Mono', monospace; padding: 16px; border-radius: 10px; font-size: 12px; direction: ltr; overflow: auto; border: 1px solid #1e293b; max-height: 360px; line-height: 1.65; white-space: pre; }
        .profile-item { background: var(--bg-input); border: 1px solid var(--border); border-radius: 12px; padding: 14px 16px; margin-bottom: 10px; display: flex; align-items: center; gap: 12px; }
        .profile-info { flex: 1; min-width: 0; }
        .profile-name { font-weight: 700; font-size: 14px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
        .profile-meta { font-size: 11px; color: var(--text-dim); margin-top: 3px; display: flex; gap: 10px; flex-wrap: wrap; }
        .pbadge { display: inline-block; padding: 2px 8px; border-radius: 6px; font-size: 10px; font-weight: 700; text-transform: uppercase; background: rgba(59,130,246,.2); color: var(--primary); }
        .pbadge.vmess  { background: rgba(99,102,241,.2); color: #818cf8; }
        .pbadge.vless  { background: rgba(16,185,129,.2); color: #4ade80; }
        .pbadge.trojan { background: rgba(245,158,11,.2); color: #fbbf24; }
        .profile-actions { display: flex; gap: 6px; flex-shrink: 0; }
        .status-bar { display: flex; align-items: center; gap: 8px; padding: 10px 14px; border-radius: 10px; font-size: 13px; margin-top: 12px; }
        .status-bar.ok   { background: rgba(16,185,129,.1);  border: 1px solid rgba(16,185,129,.25);  color: var(--success); }
        .status-bar.err  { background: rgba(239,68,68,.1);   border: 1px solid rgba(239,68,68,.25);   color: var(--error); }
        .status-bar.warn { background: rgba(245,158,11,.1);  border: 1px solid rgba(245,158,11,.25);  color: var(--warning); }
        .sdot { width: 8px; height: 8px; border-radius: 50%; background: currentColor; flex-shrink: 0; }
        .info-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; margin-top: 10px; }
        .info-cell { background: var(--bg-input); border: 1px solid var(--border); border-radius: 8px; padding: 10px 12px; }
        .info-cell-label { font-size: 10px; color: var(--text-dim); margin-bottom: 3px; text-transform: uppercase; letter-spacing: .5px; }
        .info-cell-val { font-size: 13px; font-weight: 600; font-family: 'JetBrains Mono', monospace; word-break: break-all; }
        #toast-box { position: fixed; bottom: 20px; left: 20px; display: flex; flex-direction: column; gap: 8px; z-index: 9999; }
        .toast { background: var(--bg-card); border-left: 4px solid var(--primary); padding: 12px 18px; border-radius: 10px; font-size: 13px; font-weight: 500; box-shadow: 0 8px 24px rgba(0,0,0,.4); animation: toastIn .3s ease; max-width: 280px; }
        .toast.err  { border-left-color: var(--error); }
        .toast.warn { border-left-color: var(--warning); }
        @keyframes toastIn { from { transform: translateX(-110%); opacity: 0; } to { transform: none; opacity: 1; } }
        .spinner { width: 18px; height: 18px; border: 2px solid rgba(255,255,255,.2); border-top-color: #fff; border-radius: 50%; animation: spin .8s linear infinite; }
        @keyframes spin { to { transform: rotate(360deg); } }
        .hidden { display: none !important; }
        .fade-in { animation: fadeIn .35s ease; }
        @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: none; } }
        .divider { height: 1px; background: var(--border); margin: 14px 0; }
        .import-pw-box { margin-top: 12px; padding: 14px; background: rgba(245,158,11,.07); border: 1px solid rgba(245,158,11,.3); border-radius: 10px; }
    </style>
</head>
<body>
<div class="container">

    <header>
        <div class="logo">⚡ Napster <span>Config</span></div>
        <div class="hcontrols">
            <div class="color-picker">
                <div class="cdot" style="background:#3b82f6" onclick="setTheme('')"></div>
                <div class="cdot" style="background:#10b981" onclick="setTheme('emerald')"></div>
                <div class="cdot" style="background:#f43f5e" onclick="setTheme('rose')"></div>
                <div class="cdot" style="background:#f59e0b" onclick="setTheme('amber')"></div>
                <div class="cdot" style="background:#8b5cf6" onclick="setTheme('violet')"></div>
            </div>
            <button class="btn-icon" onclick="toggleLight()" id="themeBtn">🌙</button>
        </div>
    </header>

    <nav class="nav-tabs">
        <button class="tab-btn active" id="tab-btn-0" onclick="openTab(0)">⚙️ تولید کانفیگ</button>
        <button class="tab-btn"        id="tab-btn-1" onclick="openTab(1)">📂 پروفایل‌ها</button>
        <button class="tab-btn"        id="tab-btn-2" onclick="openTab(2)">📥 ایمپورت NPVT</button>
    </nav>

    <!-- TAB 0 — Generate -->
    <div id="tab-0" class="fade-in">
        <div class="grid2">
            <div>
                <div class="card">
                    <div class="card-title">🔗 لینک V2Ray</div>
                    <div class="fld">
                        <label>لینک vmess / vless / trojan</label>
                        <textarea id="v2link" rows="4" placeholder="vmess://... یا vless://... یا trojan://..."></textarea>
                    </div>
                    <div style="display:flex;gap:8px">
                        <button class="btn btn-primary" style="flex:2" onclick="apiParse()">🔍 آنالیز لینک</button>
                        <button class="btn btn-outline" style="flex:1" onclick="apiServerInfo()">📡 تست سرور</button>
                    </div>
                    <div id="parse-result" class="hidden">
                        <div class="info-grid" id="parse-grid"></div>
                    </div>
                    <div id="server-result" class="hidden">
                        <div class="status-bar" id="server-bar"><div class="sdot"></div><span id="server-text"></span></div>
                        <div class="info-grid" id="server-grid"></div>
                    </div>
                </div>
                <div class="card">
                    <div class="card-title">🏷️ مشخصات پروفایل</div>
                    <div class="fld">
                        <label>نام نمایشی</label>
                        <input type="text" id="profileName" placeholder="مثلاً: سرور آلمان">
                    </div>
                </div>
            </div>

            <div>
                <div class="card">
                    <div class="card-title">🛡️ تنظیمات</div>

                    <div class="sw-row">
                        <div><div class="sw-label">رمز فایل NPVT</div><div class="sw-desc">فایل خروجی رمزگذاری می‌شود</div></div>
                        <label class="switch"><input type="checkbox" id="enablePassword" onchange="toggleExtra('pw-extra',this)"><span class="slider"></span></label>
                    </div>
                    <div id="pw-extra" class="sw-extra hidden">
                        <label style="font-size:12px;color:var(--text-dim);display:block;margin-bottom:5px">رمز عبور</label>
                        <input type="password" id="password" placeholder="Password...">
                    </div>

                    <div class="sw-row">
                        <div><div class="sw-label">قفل دیوایس</div><div class="sw-desc">فایل فقط روی این دستگاه باز می‌شود</div></div>
                        <label class="switch"><input type="checkbox" id="enableDeviceLock" onchange="toggleExtra('dl-extra',this)"><span class="slider"></span></label>
                    </div>
                    <div id="dl-extra" class="sw-extra hidden" style="display:none;flex;gap:8px;align-items:flex-end">
                        <label style="font-size:12px;color:var(--text-dim);display:block;margin-bottom:5px">Device ID</label>
                        <div style="display:flex;gap:8px">
                            <input type="text" id="deviceId" placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" style="flex:1">
                            <button class="btn btn-outline" style="flex-shrink:0" onclick="apiDeviceId()">🎲 تولید</button>
                        </div>
                    </div>

                    <div class="sw-row">
                        <div><div class="sw-label">فعال‌سازی IPv6</div></div>
                        <label class="switch"><input type="checkbox" id="enableIpv6"><span class="slider"></span></label>
                    </div>

                    <div class="divider"></div>

                    <div class="row2" style="margin-bottom:14px">
                        <div class="fld" style="margin:0">
                            <label>Proxy Mode</label>
                            <select id="proxyMode">
                                <option value="rule">Rule (توصیه‌شده)</option>
                                <option value="global">Global</option>
                                <option value="direct">Direct</option>
                            </select>
                        </div>
                        <div class="fld" style="margin:0">
                            <label>Log Level</label>
                            <select id="logLevel">
                                <option value="warning">Warning</option>
                                <option value="info">Info</option>
                                <option value="debug">Debug</option>
                                <option value="error">Error</option>
                                <option value="silent">Silent</option>
                            </select>
                        </div>
                    </div>

                    <div class="fld">
                        <label>User Agent</label>
                        <input type="text" id="userAgent" value="Napster/2.0">
                    </div>

                    <div class="sw-row">
                        <div><div class="sw-label">Bypass ایران</div><div class="sw-desc">دامنه‌ها و IP های ایران مستقیم</div></div>
                        <label class="switch"><input type="checkbox" id="enableBypass" checked onchange="toggleExtra('bypass-extra',this)"><span class="slider"></span></label>
                    </div>
                    <div id="bypass-extra" class="sw-extra">
                        <div class="row2">
                            <div class="fld" style="margin:0">
                                <label>دامنه‌ها (هر خط یکی)</label>
                                <textarea id="bypassDomains" rows="4" style="font-size:11px;font-family:'JetBrains Mono',monospace">ir
shaparak.ir
digikala.com
aparat.com
snapp.ir
divar.ir</textarea>
                            </div>
                            <div class="fld" style="margin:0">
                                <label>IP Ranges (هر خط یکی)</label>
                                <textarea id="bypassIPs" rows="4" style="font-size:11px;font-family:'JetBrains Mono',monospace">192.168.0.0/16
10.0.0.0/8
172.16.0.0/12</textarea>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <button class="btn btn-primary btn-wide" id="gen-btn" onclick="apiGenerate()" style="margin-bottom:24px">✨ تولید کانفیگ نهایی</button>

        <div id="gen-result" class="hidden fade-in">
            <div class="card">
                <div class="card-title">📄 کانفیگ Clash (YAML)</div>
                <div class="terminal" id="yaml-out"></div>
                <div class="btn-row">
                    <button class="btn btn-outline" onclick="copyEl('yaml-out')">📋 کپی کانفیگ</button>
                    <button class="btn btn-outline" id="btn-dl-yaml">⬇️ دانلود .yaml</button>
                    <button class="btn btn-primary"  id="btn-dl-npvt">📦 دانلود .npvt</button>
                </div>
            </div>
        </div>
    </div>

    <!-- TAB 1 — Profiles -->
    <div id="tab-1" class="hidden fade-in">
        <div class="card">
            <div class="card-title">📂 پروفایل‌های ذخیره‌شده</div>
            <div id="profiles-container"><div style="text-align:center;padding:40px;color:var(--text-dim)">در حال بارگذاری…</div></div>
        </div>
    </div>

    <!-- TAB 2 — Import -->
    <div id="tab-2" class="hidden fade-in">
        <div class="card">
            <div class="card-title">📥 ایمپورت فایل NPVT</div>
            <div class="fld">
                <label>محتوای فایل .npvt را پیست کنید (NPVT1 text یا Base64)</label>
                <textarea id="imp-content" rows="5" style="font-family:'JetBrains Mono',monospace;font-size:12px" placeholder="NPVT1 ... یا Base64..."></textarea>
            </div>
            <div id="imp-pw-box" class="import-pw-box hidden">
                <div style="font-size:13px;color:var(--warning);margin-bottom:8px;font-weight:600">🔐 فایل رمزگذاری شده — رمز عبور را وارد کنید:</div>
                <input type="password" id="imp-password" placeholder="رمز عبور فایل NPVT...">
                <button class="btn btn-primary btn-wide" style="margin-top:10px" onclick="apiImport(true)">🔓 رمزگشایی و ایمپورت</button>
            </div>
            <button class="btn btn-primary btn-wide" style="margin-top:12px" onclick="apiImport(false)">🔍 بررسی و استخراج</button>
        </div>
        <div id="imp-result" class="hidden fade-in">
            <div class="card">
                <div class="card-title">✅ اطلاعات استخراج‌شده</div>
                <div id="imp-meta" class="info-grid" style="margin-bottom:14px"></div>
                <div style="margin-bottom:6px;font-size:12px;color:var(--text-dim)">V2Ray Link:</div>
                <div style="display:flex;gap:8px;margin-bottom:14px">
                    <input type="text" id="imp-v2link" readonly style="font-family:'JetBrains Mono',monospace;font-size:12px;flex:1">
                    <button class="btn btn-outline" onclick="copyVal('imp-v2link')">📋</button>
                </div>
                <div style="margin-bottom:6px;font-size:12px;color:var(--text-dim)">Config YAML:</div>
                <div class="terminal" id="imp-yaml" style="max-height:260px"></div>
                <div class="btn-row">
                    <button class="btn btn-success" onclick="loadImportedToGen()">⚙️ بارگذاری در سازنده</button>
                    <button class="btn btn-outline" onclick="copyEl('imp-yaml')">📋 کپی کانفیگ</button>
                </div>
            </div>
        </div>
    </div>

</div>
<div id="toast-box"></div>

<script>
var currentNpvtB64 = "";
var currentProfileName = "config";
var importedData = null;
var lightMode = false;

function toggleLight() {
    lightMode = !lightMode;
    document.body.classList.toggle('light-mode', lightMode);
    document.getElementById('themeBtn').textContent = lightMode ? '☀️' : '🌙';
}
function setTheme(name) {
    document.body.classList.remove('theme-emerald','theme-rose','theme-amber','theme-violet');
    if (name) document.body.classList.add('theme-' + name);
}
function openTab(idx) {
    [0,1,2].forEach(function(i) {
        document.getElementById('tab-' + i).classList.toggle('hidden', i !== idx);
        document.getElementById('tab-btn-' + i).classList.toggle('active', i === idx);
    });
    if (idx === 1) apiLoadProfiles();
}
function toggleExtra(id, cb) {
    document.getElementById(id).classList.toggle('hidden', !cb.checked);
}
function showToast(msg, type) {
    var t = document.createElement('div');
    t.className = 'toast' + (type === 'err' ? ' err' : type === 'warn' ? ' warn' : '');
    t.textContent = msg;
    document.getElementById('toast-box').appendChild(t);
    setTimeout(function() { t.remove(); }, 3200);
}
function copyEl(id) {
    var text = document.getElementById(id).innerText || document.getElementById(id).textContent;
    navigator.clipboard.writeText(text).then(function() { showToast('کپی شد ✓'); });
}
function copyVal(id) {
    navigator.clipboard.writeText(document.getElementById(id).value).then(function() { showToast('کپی شد ✓'); });
}
function dlText(content, filename, mime) {
    var a = document.createElement('a');
    a.href = URL.createObjectURL(new Blob([content], {type: mime || 'text/plain'}));
    a.download = filename;
    a.click();
}
function infoCell(label, val) {
    return '<div class="info-cell"><div class="info-cell-label">' + label + '</div><div class="info-cell-val">' + escHtml(val) + '</div></div>';
}
function escHtml(str) {
    return String(str || '—').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}
function collectSettings() {
    return {
        enablePassword:   document.getElementById('enablePassword').checked,
        password:         document.getElementById('password').value,
        enableDeviceLock: document.getElementById('enableDeviceLock').checked,
        deviceId:         document.getElementById('deviceId').value,
        userAgent:        document.getElementById('userAgent').value || 'Napster/2.0',
        proxyMode:        document.getElementById('proxyMode').value,
        logLevel:         document.getElementById('logLevel').value,
        enableIpv6:       document.getElementById('enableIpv6').checked,
        enableBypass:     document.getElementById('enableBypass').checked,
        bypassDomains:    document.getElementById('bypassDomains').value.split('\n').map(function(s){return s.trim();}).filter(Boolean),
        bypassIPs:        document.getElementById('bypassIPs').value.split('\n').map(function(s){return s.trim();}).filter(Boolean)
    };
}
function fillSettings(s) {
    if (!s) return;
    function set(id, v) { var el = document.getElementById(id); if(el) el.value = v || ''; }
    function chk(id, v) { var el = document.getElementById(id); if(el) el.checked = !!v; }
    chk('enablePassword',   s.enablePassword);
    chk('enableDeviceLock', s.enableDeviceLock);
    chk('enableBypass',     s.enableBypass !== false);
    chk('enableIpv6',       s.enableIpv6);
    set('userAgent',  s.userAgent || 'Napster/2.0');
    set('proxyMode',  s.proxyMode  || 'rule');
    set('logLevel',   s.logLevel   || 'warning');
    set('deviceId',   s.deviceId   || '');
    if (s.bypassDomains && s.bypassDomains.length) document.getElementById('bypassDomains').value = s.bypassDomains.join('\n');
    if (s.bypassIPs     && s.bypassIPs.length)     document.getElementById('bypassIPs').value     = s.bypassIPs.join('\n');
    document.getElementById('pw-extra').classList.toggle('hidden', !s.enablePassword);
    document.getElementById('dl-extra').classList.toggle('hidden', !s.enableDeviceLock);
    document.getElementById('bypass-extra').classList.toggle('hidden', !s.enableBypass);
}

// ── API: Parse
async function apiParse() {
    var link = document.getElementById('v2link').value.trim();
    if (!link) { showToast('لینک وارد کنید', 'err'); return; }
    try {
        var r = await fetch('/api/parse?link=' + encodeURIComponent(link));
        var d = await r.json();
        if (d.error) { showToast(d.error, 'err'); return; }
        if (d.remarks) document.getElementById('profileName').value = d.remarks;
        document.getElementById('parse-grid').innerHTML =
            infoCell('پروتکل', d.protocol) + infoCell('آدرس', d.address) +
            infoCell('پورت',   d.port)     + infoCell('شبکه', d.network) +
            infoCell('TLS',    d.tls)      + infoCell('SNI',   d.sni);
        document.getElementById('parse-result').classList.remove('hidden');
        showToast('لینک آنالیز شد ✓');
    } catch(e) { showToast('خطا در آنالیز', 'err'); }
}

// ── API: Server Info
async function apiServerInfo() {
    var link = document.getElementById('v2link').value.trim();
    var addr = '';
    try {
        if (link.startsWith('vmess://')) {
            var json = JSON.parse(atob(link.replace('vmess://','')));
            addr = json.add;
        } else {
            var u = new URL(link);
            addr = u.hostname;
        }
    } catch(e) { showToast('ابتدا لینک را وارد کنید', 'err'); return; }
    var bar = document.getElementById('server-bar');
    var txt = document.getElementById('server-text');
    document.getElementById('server-result').classList.remove('hidden');
    bar.className = 'status-bar warn';
    txt.textContent = 'در حال بررسی…';
    document.getElementById('server-grid').innerHTML = '';
    try {
        var r = await fetch('/api/server-info?address=' + encodeURIComponent(addr));
        var d = await r.json();
        var online = d.status === 'online';
        bar.className = 'status-bar ' + (online ? 'ok' : 'err');
        txt.textContent = online ? ('آنلاین — ' + d.ping) : 'آفلاین';
        document.getElementById('server-grid').innerHTML = infoCell('کشور', d.country) + infoCell('ISP', d.org);
    } catch(e) { bar.className = 'status-bar err'; txt.textContent = 'خطا در اتصال'; }
}

// ── API: Device ID
async function apiDeviceId() {
    try {
        var r = await fetch('/api/device-id');
        var d = await r.json();
        document.getElementById('deviceId').value = d.deviceId;
        showToast('Device ID ساخته شد');
    } catch(e) { showToast('خطا', 'err'); }
}

// ── API: Generate
async function apiGenerate() {
    var v2link = document.getElementById('v2link').value.trim();
    if (!v2link) { showToast('لینک V2Ray را وارد کنید', 'err'); return; }
    var btn = document.getElementById('gen-btn');
    btn.disabled = true;
    btn.innerHTML = '<div class="spinner"></div> در حال تولید…';
    var profileName = document.getElementById('profileName').value || 'config';
    try {
        var r = await fetch('/api/generate', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ v2rayLink: v2link, profileName: profileName, settings: collectSettings() })
        });
        var d = await r.json();
        if (!d.success) { showToast(d.error || 'خطا در تولید', 'err'); return; }
        currentNpvtB64    = d.npvtB64 || '';
        currentProfileName = profileName;
        document.getElementById('yaml-out').textContent = d.config;
        document.getElementById('gen-result').classList.remove('hidden');
        // Download YAML
        document.getElementById('btn-dl-yaml').onclick = function() {
            dlText(d.config, profileName + '.yaml', 'text/yaml');
        };
        // Download NPVT — npvtB64 is base64 of raw NPVT1 text, decode back to text then save
        document.getElementById('btn-dl-npvt').onclick = function() {
            if (!currentNpvtB64) { showToast('فایل NPVT موجود نیست', 'err'); return; }
            try {
                var npvtText = atob(currentNpvtB64);
                dlText(npvtText, profileName + '.npvt', 'application/octet-stream');
            } catch(e) { showToast('خطا در دانلود NPVT', 'err'); }
        };
        window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
        showToast('کانفیگ با موفقیت ساخته شد ✨');
    } catch(e) { showToast('خطای سیستمی', 'err'); }
    btn.disabled = false;
    btn.innerHTML = '✨ تولید کانفیگ نهایی';
}

// ── API: Profiles
async function apiLoadProfiles() {
    var c = document.getElementById('profiles-container');
    c.innerHTML = '<div style="text-align:center;padding:40px;color:var(--text-dim)">در حال بارگذاری…</div>';
    try {
        var r = await fetch('/api/profiles');
        var d = await r.json();
        if (!d.profiles || d.profiles.length === 0) {
            c.innerHTML = '<div style="text-align:center;padding:40px;color:var(--text-dim)">هیچ پروفایلی ذخیره نشده است.</div>';
            return;
        }
        c.innerHTML = '';
        d.profiles.slice().reverse().forEach(function(p) {
            var proto = (p.parsedProxy && p.parsedProxy.protocol) ? p.parsedProxy.protocol : '?';
            var addr  = p.parsedProxy ? (p.parsedProxy.address + ':' + p.parsedProxy.port) : '—';
            var div = document.createElement('div');
            div.className = 'profile-item';
            div.innerHTML =
                '<div class="profile-info">' +
                  '<div class="profile-name">' + escHtml(p.name) + ' <span class="pbadge ' + proto + '">' + proto + '</span></div>' +
                  '<div class="profile-meta"><span>🕒 ' + escHtml(p.createdAt) + '</span><span>🌐 ' + escHtml(addr) + '</span></div>' +
                '</div>' +
                '<div class="profile-actions">' +
                  '<button class="btn btn-outline" style="padding:7px 12px;font-size:12px" data-pload>📂 Load</button>' +
                  '<button class="btn btn-danger"  style="padding:7px 12px;font-size:12px" data-did="' + escHtml(p.id) + '">🗑️</button>' +
                '</div>';
            div.querySelector('[data-pload]').addEventListener('click', function() { loadProfile(p); });
            div.querySelector('[data-did]').addEventListener('click', function() { delProfile(p.id); });
            c.appendChild(div);
        });
    } catch(e) { c.innerHTML = '<div style="color:var(--error);padding:20px">خطا در بارگذاری</div>'; }
}
function loadProfile(p) {
    document.getElementById('v2link').value = p.v2rayLink || '';
    document.getElementById('profileName').value = p.name || '';
    fillSettings(p.settings);
    openTab(0);
    showToast('پروفایل بارگذاری شد');
}
async function delProfile(id) {
    if (!confirm('حذف شود؟')) return;
    await fetch('/api/profiles?id=' + id, {method: 'DELETE'});
    showToast('حذف شد');
    apiLoadProfiles();
}

// ── API: Import NPVT
async function apiImport(withPassword) {
    var content  = document.getElementById('imp-content').value.trim();
    var password = withPassword ? document.getElementById('imp-password').value : '';
    if (!content) { showToast('محتوای فایل را وارد کنید', 'err'); return; }
    try {
        var r = await fetch('/api/import-npvt', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({ content: content, password: password })
        });
        var d = await r.json();
        if (d.error === 'NEEDS_PASSWORD') {
            document.getElementById('imp-pw-box').classList.remove('hidden');
            showToast('فایل رمزگذاری شده — رمز وارد کنید', 'warn');
            return;
        }
        if (d.error) { showToast(d.error, 'err'); return; }
        importedData = d;
        document.getElementById('imp-pw-box').classList.add('hidden');
        var meta = d.meta || {};
        document.getElementById('imp-meta').innerHTML =
            infoCell('نسخه',      meta.version)   +
            infoCell('تاریخ',     meta.createdAt) +
            infoCell('رمزشده',    meta.encrypted ? 'بله ✓' : 'خیر') +
            infoCell('User-Agent', meta.userAgent);
        document.getElementById('imp-v2link').value = d.v2link || '';
        document.getElementById('imp-yaml').textContent = d.config || '';
        document.getElementById('imp-result').classList.remove('hidden');
        showToast('فایل با موفقیت باز شد ✅');
    } catch(e) { showToast('خطا در ایمپورت', 'err'); }
}
function loadImportedToGen() {
    if (!importedData) return;
    document.getElementById('v2link').value = importedData.v2link || '';
    fillSettings(importedData.settings);
    openTab(0);
    showToast('در سازنده بارگذاری شد');
}
</script>
</body>
</html>`)
