package main

// indexHTML is served at /
var indexHTML = []byte(`<!DOCTYPE html>
<html lang="fa" dir="rtl">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Napster Config Generator</title>
<style>
@import url('https://fonts.googleapis.com/css2?family=Vazirmatn:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;600&display=swap');
:root{--bg:#0a0e1a;--s:#111827;--s2:#1a2236;--br:#1e2d4a;--ac:#3b82f6;--ac2:#6366f1;--gr:#10b981;--rd:#ef4444;--yw:#f59e0b;--tx:#e2e8f0;--tx2:#94a3b8}
body.light{--bg:#f1f5f9;--s:#fff;--s2:#f8fafc;--br:#e2e8f0;--tx:#0f172a;--tx2:#64748b}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--tx);font-family:'Vazirmatn',sans-serif;min-height:100vh}
header{padding:14px 28px;border-bottom:1px solid var(--br);display:flex;align-items:center;gap:12px;background:rgba(17,24,39,.9);backdrop-filter:blur(10px);position:sticky;top:0;z-index:100}
.logo{width:36px;height:36px;background:linear-gradient(135deg,var(--ac),var(--ac2));border-radius:9px;display:flex;align-items:center;justify-content:center;font-size:18px;flex-shrink:0}
header h1{font-size:16px;font-weight:700}
header small{font-size:11px;color:var(--tx2);display:block}
.hright{margin-right:auto;display:flex;gap:8px;align-items:center}
.badge{background:rgba(59,130,246,.15);border:1px solid rgba(59,130,246,.3);color:var(--ac);padding:3px 10px;border-radius:20px;font-size:11px;font-weight:600}
nav{display:flex;border-bottom:1px solid var(--br);background:var(--s2);padding:0 28px}
.ntab{padding:11px 18px;font-size:13px;font-weight:500;cursor:pointer;border:none;border-bottom:3px solid transparent;color:var(--tx2);background:none;font-family:inherit;transition:.2s;outline:none}
.ntab.on{color:var(--ac);border-bottom-color:var(--ac);background:rgba(59,130,246,.05)}
.ntab:hover{color:var(--tx)}
.pg{display:none;padding:24px 28px;max-width:1100px;margin:0 auto}
.pg.on{display:block}
.g2{display:grid;grid-template-columns:1fr 1fr;gap:18px}
@media(max-width:760px){.g2{grid-template-columns:1fr}}
.card{background:var(--s);border:1px solid var(--br);border-radius:13px;overflow:hidden;margin-bottom:16px}
.ch{padding:13px 18px;border-bottom:1px solid var(--br);display:flex;align-items:center;gap:8px;font-weight:600;font-size:13px;background:var(--s2)}
.cb{padding:18px;display:flex;flex-direction:column;gap:13px}
.fld{display:flex;flex-direction:column;gap:5px}
.fld label{font-size:12px;color:var(--tx2);font-weight:500}
input[type=text],input[type=password],textarea{background:var(--bg);border:1px solid var(--br);color:var(--tx);padding:8px 11px;border-radius:8px;font-family:inherit;font-size:13px;width:100%;direction:ltr;transition:.2s;outline:none}
input:focus,textarea:focus{border-color:var(--ac);box-shadow:0 0 0 3px rgba(59,130,246,.12)}
textarea{resize:vertical;min-height:70px}
.rbtn{display:flex;gap:7px}
.rbtn input{flex:1}
.tr{display:flex;align-items:center;justify-content:space-between;padding:9px 11px;background:var(--bg);border:1px solid var(--br);border-radius:8px;gap:10px}
.tlb{font-size:13px;font-weight:500}
.tds{font-size:11px;color:var(--tx2);margin-top:2px}
.tog{position:relative;width:40px;height:22px;flex-shrink:0;cursor:pointer}
.tog input{opacity:0;width:0;height:0;position:absolute}
.tk{position:absolute;inset:0;background:var(--br);border-radius:22px;transition:.25s;pointer-events:none}
.tk::before{content:'';position:absolute;width:16px;height:16px;left:3px;top:3px;background:#fff;border-radius:50%;transition:.25s}
.tog input:checked + .tk{background:var(--ac)}
.tog input:checked + .tk::before{transform:translateX(18px)}
.btn{padding:8px 16px;border-radius:8px;border:none;cursor:pointer;font-family:inherit;font-size:13px;font-weight:600;transition:.15s;display:inline-flex;align-items:center;gap:6px;justify-content:center}
.btn:active{transform:scale(.97)}
.bmain{background:linear-gradient(135deg,var(--ac),var(--ac2));color:#fff;width:100%;padding:12px;font-size:14px;box-shadow:0 4px 14px rgba(59,130,246,.3)}
.bmain:hover{box-shadow:0 6px 18px rgba(59,130,246,.45);filter:brightness(1.08)}
.bsm{background:var(--s2);border:1px solid var(--br);color:var(--tx2);padding:7px 11px;font-size:12px;flex-shrink:0}
.bsm:hover{border-color:var(--ac);color:var(--ac)}
.bgr{background:rgba(16,185,129,.15);border:1px solid rgba(16,185,129,.3);color:var(--gr)}
.bgr:hover{background:rgba(16,185,129,.25)}
.bpu{background:rgba(99,102,241,.15);border:1px solid rgba(99,102,241,.3);color:var(--ac2)}
.bpu:hover{background:rgba(99,102,241,.25)}
.brd{background:rgba(239,68,68,.1);border:1px solid rgba(239,68,68,.3);color:var(--rd);padding:7px 11px;font-size:12px}
.brd:hover{background:rgba(239,68,68,.2)}
.out{width:100%;height:420px;background:var(--bg);border:1px solid var(--br);border-radius:10px;padding:13px;font-family:'JetBrains Mono',monospace;font-size:12px;color:#86efac;line-height:1.7;resize:none;direction:ltr}
.imp-ta{width:100%;min-height:120px;background:var(--bg);border:1px solid var(--br);border-radius:8px;padding:11px;font-family:'JetBrains Mono',monospace;font-size:12px;color:var(--tx);resize:vertical;direction:ltr}
.stbar{display:flex;align-items:center;gap:7px;padding:8px 11px;border-radius:8px;font-size:13px;font-weight:500}
.stbar.ok{background:rgba(16,185,129,.1);border:1px solid rgba(16,185,129,.3);color:var(--gr)}
.stbar.er{background:rgba(239,68,68,.1);border:1px solid rgba(239,68,68,.3);color:var(--rd)}
.stbar.wa{background:rgba(245,158,11,.1);border:1px solid rgba(245,158,11,.3);color:var(--yw)}
.hide{display:none!important}
.dot{width:7px;height:7px;border-radius:50%;background:currentColor;flex-shrink:0}
.pbox{background:var(--bg);border:1px solid var(--br);border-radius:8px;padding:11px;font-family:'JetBrains Mono',monospace;font-size:12px;color:var(--tx2)}
.pbdg{display:inline-block;background:rgba(59,130,246,.2);color:var(--ac);border-radius:5px;padding:2px 8px;font-size:10px;margin-bottom:6px;font-weight:700}
.sibox{display:grid;grid-template-columns:1fr 1fr;gap:7px;margin-top:4px}
.sii{background:var(--bg);border:1px solid var(--br);border-radius:7px;padding:9px 11px}
.sil{color:var(--tx2);font-size:10px;margin-bottom:2px}
.siv{font-weight:600;font-family:'JetBrains Mono',monospace;font-size:12px}
.pc{background:var(--bg);border:1px solid var(--br);border-radius:10px;padding:13px;display:flex;flex-direction:column;gap:8px;margin-bottom:10px}
.pn{font-weight:600;font-size:14px}
.pm{font-size:11px;color:var(--tx2);display:flex;gap:10px;flex-wrap:wrap;align-items:center}
.pa{display:flex;gap:7px;flex-wrap:wrap;margin-top:4px}
.pbg{display:inline-block;background:rgba(59,130,246,.15);color:var(--ac);border-radius:5px;padding:2px 7px;font-size:10px;font-weight:700}
.ibox{background:var(--bg);border:1px solid var(--br);border-radius:8px;padding:13px;font-size:13px;line-height:1.9}
::-webkit-scrollbar{width:5px;height:5px}
::-webkit-scrollbar-track{background:var(--bg)}
::-webkit-scrollbar-thumb{background:var(--br);border-radius:3px}
</style>
</head>
<body>

<header>
  <div class="logo">&#9889;</div>
  <div>
    <h1>Napster Config Generator</h1>
    <small>ساخت کانفیگ حرفه‌ای Napster VPN</small>
  </div>
  <div class="hright">
    <div class="badge">v3.1</div>
    <button class="btn bsm" id="thBtn" onclick="toggleTheme()">🌙 تاریک</button>
  </div>
</header>

<nav>
  <button class="ntab on" id="tab-gen" onclick="showTab('gen')">⚡ ساخت کانفیگ</button>
  <button class="ntab" id="tab-prof" onclick="showTab('prof')">📋 پروفایل‌ها</button>
  <button class="ntab" id="tab-imp" onclick="showTab('imp')">📂 Import .npvt</button>
</nav>

<!-- ════ TAB: GENERATE ════ -->
<div class="pg on" id="pg-gen">
<div class="g2">

<!-- ستون چپ -->
<div>

<div class="card">
  <div class="ch">🔗 لینک V2Ray</div>
  <div class="cb">
    <div class="fld">
      <label>لینک vmess / vless / trojan</label>
      <div class="rbtn">
        <input type="text" id="lnk" placeholder="vmess://... یا vless://... یا trojan://...">
        <button class="btn bsm" onclick="doParse()">📡 پارس</button>
      </div>
    </div>
    <div id="pbox" class="pbox hide"></div>
    <div id="siWrap" class="hide">
      <div style="font-size:10px;font-weight:700;color:var(--ac);text-transform:uppercase;letter-spacing:1.2px;margin-top:4px">اطلاعات سرور</div>
      <div class="sibox" id="siGrid"></div>
    </div>
    <div class="fld">
      <label>نام پروفایل (اختیاری)</label>
      <input type="text" id="pname" placeholder="مثلا: سرور آلمان">
    </div>
  </div>
</div>

<div class="card">
  <div class="ch">📱 Device ID Lock</div>
  <div class="cb">
    <div class="tr">
      <div>
        <div class="tlb">فعال کردن Device Lock</div>
        <div class="tds">کانفیگ فقط روی این Device ID کار کند</div>
      </div>
      <label class="tog">
        <input type="checkbox" id="ckLock" onchange="togEl('lockWrap', this.checked); togEl('lockOff', !this.checked)">
        <span class="tk"></span>
      </label>
    </div>
    <div id="lockWrap" class="hide">
      <div class="fld">
        <label>Device ID — UUID را از اپلیکیشن نپستر کپی کن</label>
        <div class="rbtn">
          <input type="text" id="devId" placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx">
          <button class="btn bsm" onclick="newDevId()">🔄 جدید</button>
        </div>
      </div>
    </div>
    <div id="lockOff" style="font-size:12px;color:var(--tx2)">
      ℹ️ بدون Device Lock — کانفیگ روی هر دستگاهی کار می‌کند
    </div>
  </div>
</div>

<div class="card">
  <div class="ch">🔒 رمزگذاری</div>
  <div class="cb">
    <div class="tr">
      <div>
        <div class="tlb">رمزگذاری فایل .npvt با پسورد</div>
        <div class="tds">AES-256 — بدون رمز باز نمی‌شود</div>
      </div>
      <label class="tog">
        <input type="checkbox" id="ckPwd" onchange="togEl('pwdWrap', this.checked)">
        <span class="tk"></span>
      </label>
    </div>
    <div id="pwdWrap" class="hide">
      <div class="fld">
        <label>رمز عبور (همان رمزی که خودت تعیین می‌کنی)</label>
        <input type="password" id="cfgPwd" placeholder="رمز قوی وارد کنید">
      </div>
    </div>
  </div>
</div>

<div class="card">
  <div class="ch">🇮🇷 Bypass سایت‌های ایرانی</div>
  <div class="cb">
    <div class="tr">
      <div>
        <div class="tlb">فعال کردن Bypass</div>
        <div class="tds">سایت‌های ایرانی بدون فیلترشکن باز می‌شوند</div>
      </div>
      <label class="tog">
        <input type="checkbox" id="ckBy" checked>
        <span class="tk"></span>
      </label>
    </div>
    <div class="fld">
      <label>دامنه‌های bypass (هر خط یک دامنه)</label>
      <textarea id="byDom" rows="4">ir
shaparak.ir
digikala.com
aparat.com
snapp.ir
divar.ir</textarea>
    </div>
    <div class="fld">
      <label>IP های bypass (هر خط یک CIDR)</label>
      <textarea id="byIP" rows="3">192.168.0.0/16
10.0.0.0/8
172.16.0.0/12</textarea>
    </div>
  </div>
</div>

<button class="btn bmain" onclick="doGen()">⚡ ساخت کانفیگ + خروجی .npvt</button>

</div>

<!-- ستون راست: خروجی -->
<div>
<div class="card">
  <div class="ch">📄 کانفیگ خروجی</div>
  <div class="cb">
    <div id="genSt" class="stbar hide"></div>
    <textarea class="out" id="outTa" readonly></textarea>
    <div style="display:flex;gap:7px;flex-wrap:wrap">
      <button class="btn bgr" onclick="copyOut()">📋 کپی کانفیگ</button>
      <button class="btn bpu" onclick="dlNpvt()">💾 دانلود .npvt</button>
      <button class="btn bsm" onclick="clearAll()">🗑️ پاک</button>
    </div>
  </div>
</div>
</div>

</div><!-- end g2 -->
</div><!-- end pg-gen -->

<!-- ════ TAB: PROFILES ════ -->
<div class="pg" id="pg-prof">
  <div style="max-width:900px">
    <div class="card">
      <div class="ch">📋 پروفایل‌های ذخیره‌شده</div>
      <div class="cb">
        <div id="profList">
          <div style="color:var(--tx2);text-align:center;padding:50px">کلیک بارگذاری را بزنید</div>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- ════ TAB: IMPORT ════ -->
<div class="pg" id="pg-imp">
  <div style="max-width:900px">
    <div class="card">
      <div class="ch">📂 Import کانفیگ .npvt</div>
      <div class="cb">
        <div class="fld">
          <label>محتوای فایل .npvt را paste کنید</label>
          <div style="font-size:11px;color:var(--tx2);margin-bottom:4px">
            ⚡ هر دو فرمت پشتیبانی می‌شود: <b>base64</b> (ساخته‌شده توسط این سرور) و <b>NPVT1 ...</b> (قفل‌شده توسط اپ نپستر)
          </div>
          <textarea class="imp-ta" id="impTa" rows="5" placeholder="NPVT1 ...,...,... یا eyJ2ZXJzaW9u..."></textarea>
        </div>
        <div id="impPwdWrap" class="hide">
          <div class="fld">
            <label>🔒 این فایل رمز دارد — رمز را وارد کنید</label>
            <input type="password" id="impPwd" placeholder="رمز عبور">
          </div>
        </div>
        <button class="btn bmain" onclick="doImport()">🔍 آنالیز و رمزگشایی</button>
        <div id="impSt" class="stbar hide"></div>

        <div id="impResult" class="hide">
          <div style="font-size:10px;font-weight:700;color:var(--ac);text-transform:uppercase;letter-spacing:1.2px;margin-top:12px">اطلاعات فایل</div>
          <div class="ibox" id="impMeta"></div>

          <div style="font-size:10px;font-weight:700;color:var(--ac);text-transform:uppercase;letter-spacing:1.2px;margin-top:12px">لینک V2Ray (قابل کپی و استفاده مجدد)</div>
          <div class="rbtn" style="margin-bottom:8px">
            <input type="text" id="impLink" readonly>
            <button class="btn bgr" onclick="copyImpLink()">📋 کپی</button>
          </div>

          <div style="font-size:10px;font-weight:700;color:var(--ac);text-transform:uppercase;letter-spacing:1.2px;margin-top:12px">محتوای کانفیگ</div>
          <textarea class="out" id="impOut" readonly style="height:300px"></textarea>
          <div style="display:flex;gap:7px;margin-top:8px">
            <button class="btn bgr" onclick="copyImp()">📋 کپی</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<script>
var lastCfg = '';
var lastNpvt = '';
var lastImpRaw = '';
var isDark = true;
var cachedProfs = [];

function toggleTheme() {
  isDark = !isDark;
  document.body.classList.toggle('light', !isDark);
  document.getElementById('thBtn').textContent = isDark ? '🌙 تاریک' : '☀️ روشن';
}

function showTab(name) {
  var names = ['gen', 'prof', 'imp'];
  for (var i = 0; i < names.length; i++) {
    var n = names[i];
    var tabEl = document.getElementById('tab-' + n);
    var pgEl  = document.getElementById('pg-' + n);
    if (tabEl) tabEl.classList.toggle('on', n === name);
    if (pgEl)  pgEl.classList.toggle('on', n === name);
  }
  if (name === 'prof') loadProfs();
}

function togEl(id, show) {
  var el = document.getElementById(id);
  if (!el) return;
  if (show) { el.classList.remove('hide'); } else { el.classList.add('hide'); }
}

function showSt(elId, type, msg) {
  var el = document.getElementById(elId);
  if (!el) return;
  el.className = 'stbar ' + type;
  el.innerHTML = '<div class="dot"></div>' + msg;
}
function hideSt(elId) {
  var el = document.getElementById(elId);
  if (el) el.className = 'stbar hide';
}

function gv(id) { var e = document.getElementById(id); return e ? e.value : ''; }
function sv(id, v) { var e = document.getElementById(id); if (e) e.value = v; }
function ck(id) { var e = document.getElementById(id); return e ? e.checked : false; }
function sck(id, v) { var e = document.getElementById(id); if (e) e.checked = !!v; }
function lines(id) { return gv(id).split('\n').map(function(s){return s.trim();}).filter(Boolean); }
function dl(content, filename) {
  var a = document.createElement('a');
  a.href = URL.createObjectURL(new Blob([content], {type:'application/octet-stream'}));
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
}

function doParse() {
  var link = gv('lnk').trim();
  if (!link) { showSt('genSt', 'wa', 'لینک V2Ray وارد کنید'); return; }
  fetch('/api/parse?link=' + encodeURIComponent(link))
    .then(function(r){ return r.json(); })
    .then(function(d){
      var box = document.getElementById('pbox');
      box.classList.remove('hide');
      if (d.error) { box.innerHTML = '<span style="color:var(--rd)">❌ ' + d.error + '</span>'; return; }
      var id = d.uuid || d.password || '-';
      var html = '<span class="pbdg">' + d.protocol.toUpperCase() + '</span>';
      if (d.remarks) html += ' <b>' + d.remarks + '</b>';
      html += '<br>🖥️ <b>' + d.address + ':' + d.port + '</b>';
      html += ' &nbsp;🌐 ' + (d.network || 'tcp');
      html += ' &nbsp;🔒 ' + (d.tls || 'none') + '<br>';
      html += '🔑 ' + id.substring(0, 36) + (id.length > 36 ? '...' : '');
      if (d.flow) html += '<br>⚡ flow: ' + d.flow;
      if (d.sni)  html += '<br>🌍 SNI: ' + d.sni;
      if (d.fingerprint) html += '<br>fp: ' + d.fingerprint;
      if (d.publicKey) html += '<br>🔑 Reality: ' + d.publicKey.substring(0,20) + '...';
      box.innerHTML = html;
      if (d.remarks && !gv('pname')) sv('pname', d.remarks);
      loadServerInfo(d.address);
    })
    .catch(function(e){
      var box = document.getElementById('pbox');
      box.classList.remove('hide');
      box.innerHTML = '<span style="color:var(--rd)">❌ ' + e.message + '</span>';
    });
}

function loadServerInfo(addr) {
  if (!addr) return;
  togEl('siWrap', true);
  document.getElementById('siGrid').innerHTML = '<div style="color:var(--tx2);font-size:12px;grid-column:1/-1">در حال بررسی...</div>';
  fetch('/api/server-info?address=' + encodeURIComponent(addr))
    .then(function(r){ return r.json(); })
    .then(function(d){
      var ok = d.status === 'online';
      document.getElementById('siGrid').innerHTML =
        '<div class="sii"><div class="sil">وضعیت</div><div class="siv" style="color:'+(ok?'var(--gr)':'var(--rd)')+'">'+(ok?'🟢 آنلاین':'🔴 آفلاین')+'</div></div>'+
        '<div class="sii"><div class="sil">پینگ</div><div class="siv">'+d.ping+'</div></div>'+
        '<div class="sii"><div class="sil">کشور</div><div class="siv">'+(d.countryCode||'')+' '+(d.country||'N/A')+'</div></div>'+
        '<div class="sii"><div class="sil">ISP/Org</div><div class="siv" style="font-size:11px">'+(d.org||'N/A')+'</div></div>';
    })
    .catch(function(){
      document.getElementById('siGrid').innerHTML = '<div style="color:var(--tx2);font-size:12px">خطا در دریافت اطلاعات</div>';
    });
}

function newDevId() {
  fetch('/api/device-id')
    .then(function(r){ return r.json(); })
    .then(function(d){ sv('devId', d.deviceId); })
    .catch(function(){});
}

function doGen() {
  var link = gv('lnk').trim();
  if (!link) { showSt('genSt', 'er', 'لینک V2Ray را وارد کنید'); return; }

  var settings = {
    enableDeviceLock: ck('ckLock'),
    deviceId:         gv('devId'),
    userAgent:        'Napster/2.0',
    enablePassword:   ck('ckPwd'),
    password:         gv('cfgPwd'),
    enableBypass:     ck('ckBy'),
    bypassDomains:    lines('byDom'),
    bypassIPs:        lines('byIP'),
    proxyMode:        'rule',
    logLevel:         'warning',
    enableIpv6:       false
  };

  showSt('genSt', 'wa', 'در حال ساخت...');

  fetch('/api/generate', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({v2rayLink: link, profileName: gv('pname'), settings: settings})
  })
  .then(function(r){ return r.json(); })
  .then(function(d){
    if (!d.success) { showSt('genSt', 'er', d.error); return; }
    lastCfg  = d.config;
    lastNpvt = d.npvtB64 || '';
    document.getElementById('outTa').value = d.config;
    var msg = '✅ کانفیگ ساخته شد';
    if (lastNpvt) msg += ' — فایل .npvt آماده است';
    showSt('genSt', 'ok', msg);
  })
  .catch(function(e){ showSt('genSt', 'er', 'خطا: ' + e.message); });
}

function copyOut() {
  if (!lastCfg) { showSt('genSt', 'wa', 'ابتدا کانفیگ بسازید'); return; }
  navigator.clipboard.writeText(lastCfg).then(function(){ showSt('genSt', 'ok', 'کپی شد ✓'); });
}

function dlNpvt() {
  if (!lastNpvt) { showSt('genSt', 'wa', 'ابتدا کانفیگ بسازید'); return; }
  var name = (gv('pname') || 'config').replace(/[^a-zA-Z0-9\u0600-\u06FF_-]/g, '_');
  dl(lastNpvt, name + '.npvt');
}

function clearAll() {
  lastCfg = ''; lastNpvt = '';
  sv('outTa', '');
  hideSt('genSt');
  document.getElementById('pbox').classList.add('hide');
  sv('lnk', '');
  togEl('siWrap', false);
}

/* ═══ PROFILES ═══ */
function loadProfs() {
  fetch('/api/profiles')
    .then(function(r){ return r.json(); })
    .then(function(d){
      cachedProfs = d.profiles || [];
      var list = document.getElementById('profList');
      if (!cachedProfs.length) {
        list.innerHTML = '<div style="color:var(--tx2);text-align:center;padding:50px">هنوز پروفایلی ذخیره نشده</div>';
        return;
      }
      var html = '';
      var rev = cachedProfs.slice().reverse();
      for (var i = 0; i < rev.length; i++) {
        var p = rev[i];
        var proto = (p.parsedProxy && p.parsedProxy.protocol) ? p.parsedProxy.protocol.toUpperCase() : '?';
        var srv = (p.parsedProxy && p.parsedProxy.address) ? (p.parsedProxy.address + ':' + p.parsedProxy.port) : '?';
        html += '<div class="pc">';
        html += '<div class="pn">' + (p.name || 'بدون نام') + '</div>';
        html += '<div class="pm">';
        html += '<span>📅 ' + p.createdAt + '</span>';
        html += '<span class="pbg">' + proto + '</span>';
        html += '<span>🖥️ ' + srv + '</span>';
        if (p.settings && p.settings.enableDeviceLock) html += '<span>🔒 Device Locked</span>';
        if (p.settings && p.settings.enablePassword)   html += '<span>🔐 رمزگذاری شده</span>';
        html += '</div>';
        html += '<div class="pa">';
        html += '<button class="btn bsm" onclick="viewProf(\'' + p.id + '\')">👁️ مشاهده</button>';
        html += '<button class="btn bsm" onclick="loadProf(\'' + p.id + '\')">✏️ بارگذاری</button>';
        html += '<button class="btn brd" onclick="delProf(\'' + p.id + '\')">🗑️</button>';
        html += '</div></div>';
      }
      list.innerHTML = html;
    })
    .catch(function(){
      document.getElementById('profList').innerHTML = '<div style="color:var(--rd)">خطا در بارگذاری</div>';
    });
}

function viewProf(id) {
  var p = findProf(id);
  if (!p) return;
  showTab('gen');
  sv('outTa', p.config);
  lastCfg = p.config;
  showSt('genSt', 'ok', 'پروفایل: ' + (p.name || ''));
}

function loadProf(id) {
  var p = findProf(id);
  if (!p) return;
  showTab('gen');
  sv('lnk', p.v2rayLink);
  sv('pname', p.name);
  var s = p.settings || {};
  sck('ckLock', s.enableDeviceLock);
  togEl('lockWrap', s.enableDeviceLock);
  togEl('lockOff', !s.enableDeviceLock);
  sv('devId', s.deviceId || '');
  sck('ckPwd', s.enablePassword);
  togEl('pwdWrap', s.enablePassword);
  sck('ckBy', s.enableBypass !== false);
  sv('byDom', (s.bypassDomains || []).join('\n'));
  sv('byIP', (s.bypassIPs || []).join('\n'));
  doParse();
  showSt('genSt', 'ok', 'پروفایل «' + (p.name || '') + '» بارگذاری شد');
}

function delProf(id) {
  if (!confirm('حذف شود؟')) return;
  fetch('/api/profiles?id=' + id, {method: 'DELETE'})
    .then(function(){ cachedProfs = []; loadProfs(); })
    .catch(function(){});
}

function findProf(id) {
  for (var i = 0; i < cachedProfs.length; i++) {
    if (cachedProfs[i].id === id) return cachedProfs[i];
  }
  return null;
}

/* ═══ IMPORT NPVT ═══ */
function doImport() {
  var content = gv('impTa').trim();
  var pwd = gv('impPwd');
  if (!content) { showSt('impSt', 'wa', 'محتوای فایل را وارد کنید'); return; }

  fetch('/api/import-npvt', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({content: content, password: pwd})
  })
  .then(function(r){ return r.json(); })
  .then(function(d){
    if (d.error === 'NEEDS_PASSWORD') {
      togEl('impPwdWrap', true);
      showSt('impSt', 'wa', '🔐 این فایل رمز دارد — رمز وارد کنید');
      togEl('impResult', false);
      return;
    }
    if (d.error) {
      showSt('impSt', 'er', '❌ ' + d.error);
      togEl('impResult', false);
      return;
    }

    lastImpRaw = content;
    showSt('impSt', 'ok', '✅ فایل با موفقیت خوانده شد');
    togEl('impResult', true);

    sv('impLink', d.v2link || '');
    document.getElementById('impOut').value = d.config || '';

    var m = d.meta || {};
    var pr = d.proxy || {};
    var html = '';
    if (pr.protocol) html += '🔌 پروتکل: <b>' + pr.protocol.toUpperCase() + '</b><br>';
    if (pr.address)  html += '🖥️ سرور: <b>' + pr.address + ':' + pr.port + '</b><br>';
    if (pr.remarks)  html += '📝 نام: <b>' + pr.remarks + '</b><br>';
    if (pr.tls)      html += '🔒 TLS: <b>' + pr.tls + '</b><br>';
    if (pr.network)  html += '🌐 Network: <b>' + pr.network + '</b><br>';
    if (m.deviceLock) html += '🔒 Device Lock: <b>' + m.deviceLock + '</b><br>';
    if (m.userAgent)  html += '📱 User-Agent: <b>' + m.userAgent + '</b><br>';
    if (m.createdAt)  html += '📅 ساخته شده: <b>' + m.createdAt + '</b><br>';
    if (!html) html = '<span style="color:var(--tx2)">اطلاعات اضافی یافت نشد</span>';
    document.getElementById('impMeta').innerHTML = html;
  })
  .catch(function(e){ showSt('impSt', 'er', 'خطا: ' + e.message); });
}

function copyImpLink() {
  var v = gv('impLink');
  if (!v) return;
  navigator.clipboard.writeText(v).then(function(){ showSt('impSt', 'ok', 'لینک کپی شد ✓'); });
}

function copyImp() {
  var v = document.getElementById('impOut').value;
  if (!v) return;
  navigator.clipboard.writeText(v).then(function(){ showSt('impSt', 'ok', 'کپی شد ✓'); });
}
</script>
</body>
</html>`)
