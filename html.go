package main

// indexHTML is served at /
// Kept in separate file to avoid backtick conflicts with Go raw strings
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
.pg{display:none;padding:24px 28px;max-width:1500px;margin:0 auto}
.pg.on{display:block}
.g2{display:grid;grid-template-columns:1fr 1fr;gap:18px}
.g3{display:grid;grid-template-columns:1fr 1fr 1fr;gap:12px}
@media(max-width:860px){.g2,.g3{grid-template-columns:1fr}}
.card{background:var(--s);border:1px solid var(--br);border-radius:13px;overflow:hidden;margin-bottom:16px}
.ch{padding:13px 18px;border-bottom:1px solid var(--br);display:flex;align-items:center;gap:8px;font-weight:600;font-size:13px;background:var(--s2)}
.cb{padding:18px;display:flex;flex-direction:column;gap:13px}
.sec{font-size:10px;font-weight:700;color:var(--ac);text-transform:uppercase;letter-spacing:1.2px;display:flex;align-items:center;gap:8px;margin-top:4px}
.sec::after{content:'';flex:1;height:1px;background:var(--br)}
.fld{display:flex;flex-direction:column;gap:5px}
.fld label{font-size:12px;color:var(--tx2);font-weight:500}
input[type=text],input[type=number],input[type=password],select,textarea{background:var(--bg);border:1px solid var(--br);color:var(--tx);padding:8px 11px;border-radius:8px;font-family:inherit;font-size:13px;width:100%;direction:ltr;transition:.2s;outline:none}
input:focus,select:focus,textarea:focus{border-color:var(--ac);box-shadow:0 0 0 3px rgba(59,130,246,.12)}
textarea{resize:vertical;min-height:70px}
.rbtn{display:flex;gap:7px}
.rbtn input,.rbtn select{flex:1}
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
.ibox b{color:var(--tx)}
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
    <small>&#1587;&#1575;&#1582;&#1578; &#1705;&#1575;&#1606;&#1601;&#1740;&#1711; &#1581;&#1585;&#1601;&#1607;&#8204;&#1575;&#1740; Napster VPN</small>
  </div>
  <div class="hright">
    <div class="badge">v3.0</div>
    <button class="btn bsm" id="thBtn" onclick="toggleTheme()">&#127761; &#1578;&#1575;&#1585;&#1740;&#1705;</button>
  </div>
</header>

<nav>
  <button class="ntab on" id="tab-gen" onclick="showTab('gen')">&#9889; &#1587;&#1575;&#1582;&#1578; &#1705;&#1575;&#1606;&#1601;&#1740;&#1711;</button>
  <button class="ntab" id="tab-prof" onclick="showTab('prof')">&#128203; &#1662;&#1585;&#1608;&#1601;&#1575;&#1740;&#1604;&#8204;&#1607;&#1575;</button>
  <button class="ntab" id="tab-imp" onclick="showTab('imp')">&#128194; Import .npvt</button>
</nav>

<!-- ════════════════ TAB: GENERATE ════════════════ -->
<div class="pg on" id="pg-gen">
<div class="g2">

<!-- LEFT COLUMN -->
<div>

<div class="card">
  <div class="ch">&#128279; &#1604;&#1740;&#1606;&#1705; V2Ray</div>
  <div class="cb">
    <div class="fld">
      <label>&#1604;&#1740;&#1606;&#1705; vmess / vless / trojan</label>
      <div class="rbtn">
        <input type="text" id="lnk" placeholder="vmess://... &#1740;&#1575; vless://... &#1740;&#1575; trojan://...">
        <button class="btn bsm" onclick="doParse()">&#128225; &#1662;&#1575;&#1585;&#1587;</button>
      </div>
    </div>
    <div id="pbox" class="pbox hide"></div>
    <div id="siWrap" class="hide">
      <div class="sec">&#1575;&#1591;&#1604;&#1575;&#1593;&#1575;&#1578; &#1587;&#1585;&#1608;&#1585;</div>
      <div class="sibox" id="siGrid"></div>
    </div>
    <div class="fld">
      <label>&#1606;&#1575;&#1605; &#1662;&#1585;&#1608;&#1601;&#1575;&#1740;&#1604; (&#1575;&#1582;&#1578;&#1740;&#1575;&#1585;&#1740;)</label>
      <input type="text" id="pname" placeholder="&#1605;&#1579;&#1604;&#1575;: &#1587;&#1585;&#1608;&#1585; &#1570;&#1604;&#1605;&#1575;&#1606;">
    </div>
  </div>
</div>

<div class="card">
  <div class="ch">&#128241; Device ID Lock</div>
  <div class="cb">
    <div class="tr">
      <div>
        <div class="tlb">&#1601;&#1593;&#1575;&#1604; &#1705;&#1585;&#1583;&#1606; Device Lock</div>
        <div class="tds">&#1705;&#1575;&#1606;&#1601;&#1740;&#1711; &#1601;&#1602;&#1591; &#1585;&#1608;&#1740; &#1575;&#1740;&#1606; Device ID &#1705;&#1575;&#1585; &#1705;&#1606;&#1583;</div>
      </div>
      <label class="tog">
        <input type="checkbox" id="ckLock" onchange="togEl('lockWrap', this.checked); togEl('lockOff', !this.checked)">
        <span class="tk"></span>
      </label>
    </div>
    <div id="lockWrap" class="hide">
      <div class="fld">
        <label>Device ID &#8212; UUID &#1585;&#1575; &#1575;&#1586; &#1575;&#1662;&#1604;&#1740;&#1705;&#1740;&#1588;&#1606; &#1606;&#1662;&#1587;&#1578;&#1585; &#1705;&#1662;&#1740; &#1705;&#1606;</label>
        <div class="rbtn">
          <input type="text" id="devId" placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx">
          <button class="btn bsm" onclick="newDevId()">&#128260; &#1580;&#1583;&#1740;&#1583;</button>
        </div>
      </div>
      <div class="fld">
        <label>User Agent (&#1575;&#1582;&#1578;&#1740;&#1575;&#1585;&#1740;)</label>
        <input type="text" id="ua" value="Napster/2.0">
      </div>
    </div>
    <div id="lockOff" style="font-size:12px;color:var(--tx2)">
      &#8505;&#65039; &#1576;&#1583;&#1608;&#1606; Device Lock &#8212; &#1705;&#1575;&#1606;&#1601;&#1740;&#1711; &#1585;&#1608;&#1740; &#1607;&#1585; &#1583;&#1587;&#1578;&#1711;&#1575;&#1607;&#1740; &#1705;&#1575;&#1585; &#1605;&#1740;&#8204;&#1705;&#1606;&#1583;
    </div>
  </div>
</div>

<div class="card">
  <div class="ch">&#128274; &#1575;&#1605;&#1606;&#1740;&#1578;</div>
  <div class="cb">
    <div class="sec">&#1585;&#1605;&#1586;&#1711;&#1584;&#1575;&#1585;&#1740; &#1582;&#1585;&#1608;&#1580;&#1740;</div>
    <div class="tr">
      <div>
        <div class="tlb">&#1585;&#1605;&#1586;&#1711;&#1584;&#1575;&#1585;&#1740; &#1601;&#1575;&#1740;&#1604; .npvt &#1576;&#1575; &#1662;&#1587;&#1608;&#1585;&#1583;</div>
        <div class="tds">AES-256 &#8212; &#1576;&#1583;&#1608;&#1606; &#1585;&#1605;&#1586; &#1576;&#1575;&#1586; &#1606;&#1605;&#1740;&#8204;&#1588;&#1608;&#1583;</div>
      </div>
      <label class="tog">
        <input type="checkbox" id="ckPwd" onchange="togEl('pwdWrap', this.checked)">
        <span class="tk"></span>
      </label>
    </div>
    <div id="pwdWrap" class="hide">
      <div class="fld">
        <label>&#1585;&#1605;&#1586; &#1593;&#1576;&#1608;&#1585; (&#1607;&#1605;&#1575;&#1606; &#1585;&#1605;&#1586;&#1740; &#1705;&#1607; &#1582;&#1608;&#1583;&#1578; &#1578;&#1593;&#1740;&#1740;&#1606; &#1605;&#1740;&#8204;&#1705;&#1606;&#1740;)</label>
        <input type="password" id="cfgPwd" placeholder="&#1585;&#1605;&#1586; &#1602;&#1608;&#1740; &#1608;&#1575;&#1585;&#1583; &#1705;&#1606;&#1740;&#1583;">
      </div>
    </div>

    <div class="sec">TLS &amp; Fingerprint</div>
    <div class="tr">
      <div><div class="tlb">TLS 1.3</div></div>
      <label class="tog">
        <input type="checkbox" id="ckTls" checked>
        <span class="tk"></span>
      </label>
    </div>
    <div class="tr">
      <div>
        <div class="tlb">uTLS Browser Fingerprint</div>
        <div class="tds">&#1580;&#1593;&#1604; fingerprint &#1605;&#1585;&#1608;&#1585;&#1711;&#1585; &#8212; DPI bypass</div>
      </div>
      <label class="tog">
        <input type="checkbox" id="ckUtls" checked onchange="togEl('utlsWrap', this.checked)">
        <span class="tk"></span>
      </label>
    </div>
    <div id="utlsWrap">
      <div class="fld">
        <label>&#1606;&#1608;&#1593; Fingerprint</label>
        <select id="utlsFp">
          <option value="chrome">Chrome (&#1578;&#1608;&#1589;&#1740;&#1607;)</option>
          <option value="firefox">Firefox</option>
          <option value="safari">Safari</option>
          <option value="ios">iOS Safari</option>
          <option value="android">Android Chrome</option>
          <option value="edge">Microsoft Edge</option>
          <option value="random">Random</option>
        </select>
      </div>
    </div>

    <div class="sec">Mux</div>
    <div class="tr">
      <div>
        <div class="tlb">&#1601;&#1593;&#1575;&#1604;&#8204;&#1587;&#1575;&#1586;&#1740; Mux</div>
        <div class="tds">&#1575;&#1583;&#1594;&#1575;&#1605; &#1670;&#1606;&#1583; &#1575;&#1578;&#1589;&#1575;&#1604; &#1583;&#1585; &#1740;&#1705; &#1578;&#1575;&#1606;&#1604;</div>
      </div>
      <label class="tog">
        <input type="checkbox" id="ckMux" onchange="togEl('muxWrap', this.checked)">
        <span class="tk"></span>
      </label>
    </div>
    <div id="muxWrap" class="hide g2">
      <div class="fld">
        <label>Protocol</label>
        <select id="muxProt">
          <option value="smux">smux</option>
          <option value="yamux">yamux</option>
          <option value="h2mux">h2mux</option>
        </select>
      </div>
      <div class="fld">
        <label>Concurrency</label>
        <input type="number" id="muxN" value="8" min="1" max="64">
      </div>
    </div>
  </div>
</div>

<div class="card">
  <div class="ch">&#127760; DNS</div>
  <div class="cb">
    <div class="g2">
      <div class="fld">
        <label>DNS Mode</label>
        <select id="dnsMode">
          <option value="fake-ip">Fake-IP (&#1578;&#1608;&#1589;&#1740;&#1607;)</option>
          <option value="redir-host">Redir-Host</option>
        </select>
      </div>
      <div class="fld">
        <label>&#1587;&#1585;&#1608;&#1585; DNS</label>
        <select id="dnsSrv">
          <option value="1.1.1.1">Cloudflare 1.1.1.1</option>
          <option value="8.8.8.8">Google 8.8.8.8</option>
          <option value="9.9.9.9">Quad9</option>
        </select>
      </div>
    </div>
    <div class="fld">
      <label>DNS Fallback</label>
      <select id="dnsFb">
        <option value="8.8.8.8">Google 8.8.8.8</option>
        <option value="1.1.1.1">Cloudflare</option>
        <option value="208.67.220.220">OpenDNS</option>
      </select>
    </div>
    <div class="tr">
      <div><div class="tlb">DNS over TLS</div></div>
      <label class="tog">
        <input type="checkbox" id="ckDoT" onchange="togEl('dotWrap', this.checked)">
        <span class="tk"></span>
      </label>
    </div>
    <div id="dotWrap" class="hide">
      <div class="fld">
        <label>DoT Server</label>
        <select id="dotSrv">
          <option value="cloudflare-dns.com">cloudflare-dns.com</option>
          <option value="dns.google">dns.google</option>
          <option value="dns.quad9.net">dns.quad9.net</option>
        </select>
      </div>
    </div>
  </div>
</div>

<div class="card">
  <div class="ch">&#128506;&#65039; &#1605;&#1587;&#1740;&#1585;&#1740;&#1575;&#1576;&#1740;</div>
  <div class="cb">
    <div class="g2">
      <div class="fld">
        <label>Proxy Mode</label>
        <select id="pMode">
          <option value="rule">Rule (&#1607;&#1608;&#1588;&#1605;&#1606;&#1583;)</option>
          <option value="global">Global</option>
          <option value="direct">Direct</option>
        </select>
      </div>
      <div class="fld">
        <label>Log Level</label>
        <select id="logLv">
          <option value="warning">Warning</option>
          <option value="info">Info</option>
          <option value="debug">Debug</option>
          <option value="silent">Silent</option>
        </select>
      </div>
    </div>

    <div class="tr">
      <div>
        <div class="tlb">Bypass &#8212; &#1587;&#1575;&#1740;&#1578;&#8204;&#1607;&#1575;&#1740; &#1575;&#1740;&#1585;&#1575;&#1606;&#1740; &#1605;&#1587;&#1578;&#1602;&#1740;&#1605;</div>
        <div class="tds">GEOIP &#1575;&#1740;&#1585;&#1575;&#1606; &#1607;&#1605;&#1740;&#1588;&#1607; &#1601;&#1593;&#1575;&#1604; &#1575;&#1587;&#1578;</div>
      </div>
      <label class="tog">
        <input type="checkbox" id="ckBy" checked onchange="togEl('byWrap', this.checked)">
        <span class="tk"></span>
      </label>
    </div>
    <div id="byWrap">
      <div class="fld">
        <label>&#1583;&#1575;&#1605;&#1606;&#1607;&#8204;&#1607;&#1575;&#1740; Bypass &#8212; &#1607;&#1585; &#1582;&#1591; &#1740;&#1705; &#1583;&#1575;&#1605;&#1606;&#1607;</label>
        <textarea id="byDom" rows="5">ir
shaparak.ir
digikala.com
aparat.com
snapp.ir
divar.ir</textarea>
      </div>
      <div class="fld">
        <label>IP/CIDR &#8204;&#1607;&#1575;&#1740; Bypass &#8212; &#1607;&#1585; &#1582;&#1591; &#1740;&#1705;</label>
        <textarea id="byIP" rows="3">192.168.0.0/16
10.0.0.0/8
172.16.0.0/12</textarea>
      </div>
    </div>

    <div class="tr">
      <div><div class="tlb">Block &#1583;&#1575;&#1605;&#1606;&#1607;&#8204;&#1607;&#1575;</div></div>
      <label class="tog">
        <input type="checkbox" id="ckBk" onchange="togEl('bkWrap', this.checked)">
        <span class="tk"></span>
      </label>
    </div>
    <div id="bkWrap" class="hide">
      <div class="fld">
        <label>&#1604;&#1740;&#1587;&#1578; Block</label>
        <textarea id="bkList" rows="3" placeholder="ads.example.com"></textarea>
      </div>
    </div>
  </div>
</div>

<div class="card">
  <div class="ch">&#9881;&#65039; &#1662;&#1740;&#1588;&#1585;&#1601;&#1578;&#1607;</div>
  <div class="cb">
    <div class="g3">
      <div class="fld"><label>MTU</label><input type="number" id="mtu" value="1500" min="576" max="9000"></div>
      <div class="fld"><label>IPv6</label>
        <div class="tr" style="margin-top:4px">
          <span style="font-size:12px">&#1601;&#1593;&#1575;&#1604;</span>
          <label class="tog"><input type="checkbox" id="ckV6"><span class="tk"></span></label>
        </div>
      </div>
      <div class="fld"><label>Sniffing</label>
        <div class="tr" style="margin-top:4px">
          <span style="font-size:12px">&#1601;&#1593;&#1575;&#1604;</span>
          <label class="tog"><input type="checkbox" id="ckSn" checked><span class="tk"></span></label>
        </div>
      </div>
    </div>
  </div>
</div>

<button class="btn bmain" onclick="doGen()">&#9889; &#1587;&#1575;&#1582;&#1578; &#1705;&#1575;&#1606;&#1601;&#1740;&#1711; + &#1582;&#1585;&#1608;&#1580;&#1740; .npvt</button>

</div><!-- end left col -->

<!-- RIGHT COLUMN: Output -->
<div>
<div class="card" style="position:sticky;top:72px">
  <div class="ch">&#128196; &#1582;&#1585;&#1608;&#1580;&#1740;</div>
  <div class="cb">
    <div id="genSt" class="stbar hide"></div>
    <textarea class="out" id="outTa" readonly placeholder="# &#1705;&#1575;&#1606;&#1601;&#1740;&#1711; &#1575;&#1740;&#1606;&#1580;&#1575; &#1606;&#1605;&#1575;&#1740;&#1588; &#1583;&#1575;&#1583;&#1607; &#1605;&#1740;&#8204;&#1588;&#1608;&#1583;..."></textarea>
    <div style="display:flex;gap:7px;flex-wrap:wrap">
      <button class="btn bgr" onclick="copyOut()">&#128203; &#1705;&#1662;&#1740; &#1705;&#1575;&#1606;&#1601;&#1740;&#1711;</button>
      <button class="btn bmain" style="width:auto;flex:1" onclick="dlNpvt()">&#128190; &#1583;&#1575;&#1606;&#1604;&#1608;&#1583; .npvt</button>
      <button class="btn brd" onclick="clearAll()">&#128465;&#65039;</button>
    </div>
  </div>
</div>
</div>

</div><!-- end grid -->
</div><!-- end pg-gen -->

<!-- ════════════════ TAB: PROFILES ════════════════ -->
<div class="pg" id="pg-prof">
  <div style="max-width:900px">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:18px">
      <h2 style="font-size:17px;font-weight:700">&#128203; &#1662;&#1585;&#1608;&#1601;&#1575;&#1740;&#1604;&#8204;&#1607;&#1575;&#1740; &#1584;&#1582;&#1740;&#1585;&#1607;&#8204;&#1588;&#1583;&#1607;</h2>
      <button class="btn bsm" onclick="loadProfs()">&#128260; &#1576;&#1575;&#1585;&#1711;&#1584;&#1575;&#1585;&#1740;</button>
    </div>
    <div id="profList">
      <div style="color:var(--tx2);text-align:center;padding:50px">&#1603;&#1604;&#1610;&#1603; &#1576;&#1575;&#1585;&#1711;&#1584;&#1575;&#1585;&#1740; &#1585;&#1575; &#1576;&#1586;&#1606;&#1610;&#1583;</div>
    </div>
  </div>
</div>

<!-- ════════════════ TAB: IMPORT ════════════════ -->
<div class="pg" id="pg-imp">
  <div style="max-width:900px">
    <div class="card">
      <div class="ch">&#128194; Import &#1705;&#1575;&#1606;&#1601;&#1740;&#1711; .npvt</div>
      <div class="cb">
        <div class="fld">
          <label>&#1605;&#1581;&#1578;&#1608;&#1575;&#1740; &#1601;&#1575;&#1740;&#1604; .npvt &#1585;&#1575; paste &#1603;&#1606;&#1610;&#1583; (base64 &#1601;&#1575;&#1610;&#1604;)</label>
          <textarea class="imp-ta" id="impTa" rows="5" placeholder="&#1605;&#1581;&#1578;&#1608;&#1575;&#1740; &#1601;&#1575;&#1610;&#1604; .npvt &#1585;&#1575; &#1575;&#1610;&#1606;&#1580;&#1575; paste &#1603;&#1606;&#1610;&#1583;..."></textarea>
        </div>
        <div id="impPwdWrap" class="hide">
          <div class="fld">
            <label>&#128274; &#1575;&#1610;&#1606; &#1601;&#1575;&#1610;&#1604; &#1585;&#1605;&#1586; &#1583;&#1575;&#1585;&#1583; &#8212; &#1585;&#1605;&#1586; &#1608;&#1575;&#1585;&#1583; &#1603;&#1606;&#1610;&#1583;</label>
            <input type="password" id="impPwd" placeholder="&#1585;&#1605;&#1586; &#1593;&#1576;&#1608;&#1585;">
          </div>
        </div>
        <button class="btn bmain" onclick="doImport()">&#128269; &#1570;&#1606;&#1575;&#1604;&#1610;&#1586; &#1608; &#1585;&#1605;&#1586;&#1711;&#1588;&#1575;&#1610;&#1610;</button>
        <div id="impSt" class="stbar hide"></div>

        <div id="impResult" class="hide">
          <div class="sec">&#1575;&#1591;&#1604;&#1575;&#1593;&#1575;&#1578; &#1601;&#1575;&#1610;&#1604;</div>
          <div class="ibox" id="impMeta"></div>

          <div class="sec" style="margin-top:12px">&#1604;&#1610;&#1606;&#1603; V2Ray (&#1602;&#1575;&#1576;&#1604; &#1603;&#1662;&#1610; &#1608; &#1575;&#1587;&#1578;&#1601;&#1575;&#1583;&#1607; &#1605;&#1580;&#1583;&#1583;)</div>
          <div class="rbtn" style="margin-bottom:8px">
            <input type="text" id="impLink" readonly>
            <button class="btn bgr" onclick="copyImpLink()">&#128203; &#1603;&#1662;&#1610;</button>
          </div>

          <div class="sec">&#1605;&#1581;&#1578;&#1608;&#1575;&#1610; &#1603;&#1575;&#1606;&#1601;&#1610;&#1711;</div>
          <textarea class="out" id="impOut" readonly style="height:300px"></textarea>
          <div style="display:flex;gap:7px;margin-top:8px">
            <button class="btn bgr" onclick="copyImp()">&#128203; &#1603;&#1662;&#1610;</button>
            <button class="btn bpu" onclick="dlImp()">&#128190; &#1583;&#1575;&#1606;&#1604;&#1608;&#1583; .npvt</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<script>
/* ═══════════════════════════════════════════════════════
   STATE
═══════════════════════════════════════════════════════ */
var lastCfg = '';
var lastNpvt = '';
var lastImpRaw = '';
var isDark = true;
var cachedProfs = [];

/* ═══════════════════════════════════════════════════════
   THEME
═══════════════════════════════════════════════════════ */
function toggleTheme() {
  isDark = !isDark;
  document.body.classList.toggle('light', !isDark);
  document.getElementById('thBtn').textContent = isDark ? '\uD83C\uDF19 \u062A\u0627\u0631\u06CC\u06A9' : '\u2600\uFE0F \u0631\u0648\u0634\u0646';
}

/* ═══════════════════════════════════════════════════════
   TABS
═══════════════════════════════════════════════════════ */
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

/* ═══════════════════════════════════════════════════════
   TOGGLE HELPER
═══════════════════════════════════════════════════════ */
function togEl(id, show) {
  var el = document.getElementById(id);
  if (!el) return;
  if (show) {
    el.classList.remove('hide');
  } else {
    el.classList.add('hide');
  }
}

/* ═══════════════════════════════════════════════════════
   STATUS
═══════════════════════════════════════════════════════ */
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

/* ═══════════════════════════════════════════════════════
   HELPERS
═══════════════════════════════════════════════════════ */
function gv(id) {
  var e = document.getElementById(id);
  return e ? e.value : '';
}
function sv(id, v) {
  var e = document.getElementById(id);
  if (e) e.value = v;
}
function ck(id) {
  var e = document.getElementById(id);
  return e ? e.checked : false;
}
function sck(id, v) {
  var e = document.getElementById(id);
  if (e) e.checked = !!v;
}
function lines(id) {
  return gv(id).split('\n').map(function(s) { return s.trim(); }).filter(Boolean);
}
function dl(content, filename) {
  var a = document.createElement('a');
  a.href = URL.createObjectURL(new Blob([content], {type: 'application/octet-stream'}));
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
}

/* ═══════════════════════════════════════════════════════
   PARSE LINK
═══════════════════════════════════════════════════════ */
function doParse() {
  var link = document.getElementById('lnk').value.trim();
  if (!link) { showSt('genSt', 'wa', '\u0644\u06CC\u0646\u06A9 V2Ray \u0648\u0627\u0631\u062F \u06A9\u0646\u06CC\u062F'); return; }

  fetch('/api/parse?link=' + encodeURIComponent(link))
    .then(function(r) { return r.json(); })
    .then(function(d) {
      var box = document.getElementById('pbox');
      box.classList.remove('hide');
      if (d.error) {
        box.innerHTML = '<span style="color:var(--rd)">\u274C ' + d.error + '</span>';
        return;
      }
      var id = d.uuid || d.password || '-';
      var html = '<span class="pbdg">' + d.protocol.toUpperCase() + '</span>';
      if (d.remarks) html += ' <b>' + d.remarks + '</b>';
      html += '<br>\uD83D\uDDA5\uFE0F <b>' + d.address + ':' + d.port + '</b>';
      html += ' &nbsp;\uD83C\uDF10 ' + (d.network || 'tcp');
      html += ' &nbsp;\uD83D\uDD12 ' + (d.tls || 'none') + '<br>';
      html += '\uD83D\uDD11 ' + id.substring(0, 36) + (id.length > 36 ? '...' : '');
      if (d.flow) html += '<br>\u26A1 flow: ' + d.flow;
      if (d.sni)  html += '<br>\uD83C\uDF0D SNI: ' + d.sni;
      if (d.fingerprint) html += '<br>fp: ' + d.fingerprint;
      if (d.publicKey) html += '<br>\uD83D\uDD11 Reality: ' + d.publicKey.substring(0, 20) + '...';
      box.innerHTML = html;
      if (d.remarks && !document.getElementById('pname').value)
        document.getElementById('pname').value = d.remarks;
      loadServerInfo(d.address);
    })
    .catch(function(e) {
      var box = document.getElementById('pbox');
      box.classList.remove('hide');
      box.innerHTML = '<span style="color:var(--rd)">\u274C ' + e.message + '</span>';
    });
}

function loadServerInfo(addr) {
  if (!addr) return;
  togEl('siWrap', true);
  document.getElementById('siGrid').innerHTML = '<div style="color:var(--tx2);font-size:12px;grid-column:1/-1">\u062F\u0631 \u062D\u0627\u0644 \u0628\u0631\u0631\u0633\u06CC...</div>';
  fetch('/api/server-info?address=' + encodeURIComponent(addr))
    .then(function(r) { return r.json(); })
    .then(function(d) {
      var ok = d.status === 'online';
      document.getElementById('siGrid').innerHTML =
        '<div class="sii"><div class="sil">\u0648\u0636\u0639\u06CC\u062A</div><div class="siv" style="color:' + (ok ? 'var(--gr)' : 'var(--rd)') + '">' + (ok ? '\uD83D\uDFE2 \u0622\u0646\u0644\u0627\u06CC\u0646' : '\uD83D\uDD34 \u0622\u0641\u0644\u0627\u06CC\u0646') + '</div></div>' +
        '<div class="sii"><div class="sil">\u067E\u06CC\u0646\u06AF</div><div class="siv">' + d.ping + '</div></div>' +
        '<div class="sii"><div class="sil">\u06A9\u0634\u0648\u0631</div><div class="siv">' + (d.countryCode || '') + ' ' + (d.country || 'N/A') + '</div></div>' +
        '<div class="sii"><div class="sil">ISP/Org</div><div class="siv" style="font-size:11px">' + (d.org || 'N/A') + '</div></div>';
    })
    .catch(function() {
      document.getElementById('siGrid').innerHTML = '<div style="color:var(--tx2);font-size:12px">\u062E\u0637\u0627 \u062F\u0631 \u062F\u0631\u06CC\u0627\u0641\u062A \u0627\u0637\u0644\u0627\u0639\u0627\u062A</div>';
    });
}

/* ═══════════════════════════════════════════════════════
   NEW DEVICE ID
═══════════════════════════════════════════════════════ */
function newDevId() {
  fetch('/api/device-id')
    .then(function(r) { return r.json(); })
    .then(function(d) { document.getElementById('devId').value = d.deviceId; })
    .catch(function() {});
}

/* ═══════════════════════════════════════════════════════
   GENERATE
═══════════════════════════════════════════════════════ */
function doGen() {
  var link = document.getElementById('lnk').value.trim();
  if (!link) { showSt('genSt', 'er', '\u0644\u06CC\u0646\u06A9 V2Ray \u0631\u0627 \u0648\u0627\u0631\u062F \u06A9\u0646\u06CC\u062F'); return; }

  var settings = {
    enableDeviceLock: ck('ckLock'),
    deviceId:         gv('devId'),
    userAgent:        gv('ua'),
    enablePassword:   ck('ckPwd'),
    password:         gv('cfgPwd'),
    enableTls13:      ck('ckTls'),
    enableUtls:       ck('ckUtls'),
    utlsFingerprint:  gv('utlsFp'),
    enableMux:        ck('ckMux'),
    muxConcurrency:   parseInt(gv('muxN')) || 8,
    muxProtocol:      gv('muxProt'),
    dnsMode:          gv('dnsMode'),
    dnsServer:        gv('dnsSrv'),
    fallbackDNS:      gv('dnsFb'),
    enableDnsoTls:    ck('ckDoT'),
    dnsotlsServer:    gv('dotSrv'),
    proxyMode:        gv('pMode'),
    logLevel:         gv('logLv'),
    enableBypass:     ck('ckBy'),
    bypassDomains:    lines('byDom'),
    bypassIPs:        lines('byIP'),
    enableBlock:      ck('ckBk'),
    blockList:        lines('bkList'),
    mtu:              parseInt(gv('mtu')) || 1500,
    enableIpv6:       ck('ckV6'),
    enableSniffing:   ck('ckSn')
  };

  showSt('genSt', 'wa', '\u062F\u0631 \u062D\u0627\u0644 \u0633\u0627\u062E\u062A...');

  fetch('/api/generate', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({v2rayLink: link, profileName: gv('pname'), settings: settings})
  })
  .then(function(r) { return r.json(); })
  .then(function(d) {
    if (!d.success) { showSt('genSt', 'er', d.error); return; }
    lastCfg  = d.config;
    lastNpvt = d.npvtB64 || '';
    document.getElementById('outTa').value = d.config;
    var msg = '\u2705 \u06A9\u0627\u0646\u0641\u06CC\u06AF \u0633\u0627\u062E\u062A\u0647 \u0634\u062F';
    if (lastNpvt) msg += ' \u2014 \u0641\u0627\u06CC\u0644 .npvt \u0622\u0645\u0627\u062F\u0647 \u0627\u0633\u062A';
    showSt('genSt', 'ok', msg);
  })
  .catch(function(e) { showSt('genSt', 'er', '\u062E\u0637\u0627: ' + e.message); });
}

function copyOut() {
  if (!lastCfg) { showSt('genSt', 'wa', '\u0627\u0628\u062A\u062F\u0627 \u06A9\u0627\u0646\u0641\u06CC\u06AF \u0628\u0633\u0627\u0632\u06CC\u062F'); return; }
  navigator.clipboard.writeText(lastCfg).then(function() { showSt('genSt', 'ok', '\u06A9\u067E\u06CC \u0634\u062F \u2713'); });
}

function dlNpvt() {
  if (!lastNpvt) { showSt('genSt', 'wa', '\u0627\u0628\u062A\u062F\u0627 \u06A9\u0627\u0646\u0641\u06CC\u06AF \u0628\u0633\u0627\u0632\u06CC\u062F'); return; }
  var name = (gv('pname') || 'config').replace(/[^a-zA-Z0-9\u0600-\u06FF_-]/g, '_');
  dl(lastNpvt, name + '.npvt');
}

function clearAll() {
  lastCfg = ''; lastNpvt = '';
  document.getElementById('outTa').value = '';
  hideSt('genSt');
  document.getElementById('pbox').classList.add('hide');
  document.getElementById('lnk').value = '';
  togEl('siWrap', false);
}

/* ═══════════════════════════════════════════════════════
   PROFILES
═══════════════════════════════════════════════════════ */
function loadProfs() {
  fetch('/api/profiles')
    .then(function(r) { return r.json(); })
    .then(function(d) {
      cachedProfs = d.profiles || [];
      var list = document.getElementById('profList');
      if (!cachedProfs.length) {
        list.innerHTML = '<div style="color:var(--tx2);text-align:center;padding:50px">\u0647\u0646\u0648\u0632 \u067E\u0631\u0648\u0641\u0627\u06CC\u0644\u06CC \u0630\u062E\u06CC\u0631\u0647 \u0646\u0634\u062F\u0647</div>';
        return;
      }
      var html = '';
      var rev = cachedProfs.slice().reverse();
      for (var i = 0; i < rev.length; i++) {
        var p = rev[i];
        var proto = (p.parsedProxy && p.parsedProxy.protocol) ? p.parsedProxy.protocol.toUpperCase() : '?';
        var srv = (p.parsedProxy && p.parsedProxy.address) ? (p.parsedProxy.address + ':' + p.parsedProxy.port) : '?';
        html += '<div class="pc">';
        html += '<div class="pn">' + (p.name || '\u0628\u062F\u0648\u0646 \u0646\u0627\u0645') + '</div>';
        html += '<div class="pm">';
        html += '<span>\uD83D\uDCC5 ' + p.createdAt + '</span>';
        html += '<span class="pbg">' + proto + '</span>';
        html += '<span>\uD83D\uDDA5\uFE0F ' + srv + '</span>';
        if (p.settings && p.settings.enableDeviceLock) html += '<span>\uD83D\uDD12 Device Locked</span>';
        if (p.settings && p.settings.enablePassword)   html += '<span>\uD83D\uDD10 \u0631\u0645\u0632\u06AF\u0630\u0627\u0631\u06CC \u0634\u062F\u0647</span>';
        html += '</div>';
        html += '<div class="pa">';
        html += '<button class="btn bsm" onclick="viewProf(\'' + p.id + '\')">\uD83D\uDC41\uFE0F \u0645\u0634\u0627\u0647\u062F\u0647</button>';
        html += '<button class="btn bsm" onclick="loadProf(\'' + p.id + '\')">\u270F\uFE0F \u0628\u0627\u0631\u06AF\u0630\u0627\u0631\u06CC</button>';
        html += '<button class="btn brd" onclick="delProf(\'' + p.id + '\')">\uD83D\uDDD1\uFE0F</button>';
        html += '</div></div>';
      }
      list.innerHTML = html;
    })
    .catch(function() {
      document.getElementById('profList').innerHTML = '<div style="color:var(--rd)">\u062E\u0637\u0627 \u062F\u0631 \u0628\u0627\u0631\u06AF\u0630\u0627\u0631\u06CC</div>';
    });
}

function viewProf(id) {
  var p = findProf(id);
  if (!p) return;
  showTab('gen');
  document.getElementById('outTa').value = p.config;
  lastCfg = p.config;
  showSt('genSt', 'ok', '\u067E\u0631\u0648\u0641\u0627\u06CC\u0644: ' + (p.name || ''));
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
  sv('ua', s.userAgent || '');
  sck('ckPwd', s.enablePassword);
  togEl('pwdWrap', s.enablePassword);
  sck('ckTls', s.enableTls13 !== false);
  sck('ckUtls', s.enableUtls !== false);
  togEl('utlsWrap', s.enableUtls !== false);
  sv('utlsFp', s.utlsFingerprint || 'chrome');
  sck('ckMux', s.enableMux);
  togEl('muxWrap', s.enableMux);
  sv('muxProt', s.muxProtocol || 'smux');
  sv('muxN', s.muxConcurrency || 8);
  sv('dnsMode', s.dnsMode || 'fake-ip');
  sv('dnsSrv', s.dnsServer || '1.1.1.1');
  sv('dnsFb', s.fallbackDNS || '8.8.8.8');
  sck('ckDoT', s.enableDnsoTls);
  togEl('dotWrap', s.enableDnsoTls);
  sv('dotSrv', s.dnsotlsServer || 'cloudflare-dns.com');
  sv('pMode', s.proxyMode || 'rule');
  sv('logLv', s.logLevel || 'warning');
  sck('ckBy', s.enableBypass !== false);
  togEl('byWrap', s.enableBypass !== false);
  sv('byDom', (s.bypassDomains || []).join('\n'));
  sv('byIP', (s.bypassIPs || []).join('\n'));
  sck('ckBk', s.enableBlock);
  togEl('bkWrap', s.enableBlock);
  sv('bkList', (s.blockList || []).join('\n'));
  sv('mtu', s.mtu || 1500);
  sck('ckV6', s.enableIpv6);
  sck('ckSn', s.enableSniffing !== false);
  doParse();
  showSt('genSt', 'ok', '\u067E\u0631\u0648\u0641\u0627\u06CC\u0644 \u00AB' + (p.name || '') + '\u00BB \u0628\u0627\u0631\u06AF\u0630\u0627\u0631\u06CC \u0634\u062F');
}

function delProf(id) {
  if (!confirm('\u062D\u0630\u0641 \u0634\u0648\u062F\u061F')) return;
  fetch('/api/profiles?id=' + id, {method: 'DELETE'})
    .then(function() { cachedProfs = []; loadProfs(); })
    .catch(function() {});
}

function findProf(id) {
  for (var i = 0; i < cachedProfs.length; i++) {
    if (cachedProfs[i].id === id) return cachedProfs[i];
  }
  return null;
}

/* ═══════════════════════════════════════════════════════
   IMPORT NPVT
═══════════════════════════════════════════════════════ */
function doImport() {
  var content = document.getElementById('impTa').value.trim();
  var pwd = document.getElementById('impPwd').value;
  if (!content) { showSt('impSt', 'wa', '\u0645\u062D\u062A\u0648\u0627\u06CC \u0641\u0627\u06CC\u0644 \u0631\u0627 \u0648\u0627\u0631\u062F \u06A9\u0646\u06CC\u062F'); return; }

  fetch('/api/import-npvt', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({content: content, password: pwd})
  })
  .then(function(r) { return r.json(); })
  .then(function(d) {
    if (d.error === 'NEEDS_PASSWORD') {
      togEl('impPwdWrap', true);
      showSt('impSt', 'wa', '\uD83D\uDD10 \u0627\u06CC\u0646 \u0641\u0627\u06CC\u0644 \u0631\u0645\u0632 \u062F\u0627\u0631\u062F \u2014 \u0631\u0645\u0632 \u0648\u0627\u0631\u062F \u06A9\u0646\u06CC\u062F');
      togEl('impResult', false);
      return;
    }
    if (d.error) {
      showSt('impSt', 'er', '\u274C ' + d.error);
      togEl('impResult', false);
      return;
    }

    lastImpRaw = content;
    showSt('impSt', 'ok', '\u2705 \u0641\u0627\u06CC\u0644 \u0628\u0627 \u0645\u0648\u0641\u0642\u06CC\u062A \u062E\u0648\u0627\u0646\u062F\u0647 \u0634\u062F');
    togEl('impResult', true);

    // Fill link
    sv('impLink', d.v2link || '');

    // Fill config text
    document.getElementById('impOut').value = d.config || '';

    // Meta info
    var m = d.meta || {};
    var pr = d.proxy || {};
    var html = '';
    if (pr.protocol) html += '\uD83D\uDD0C \u067E\u0631\u0648\u062A\u06A9\u0644: <b>' + pr.protocol.toUpperCase() + '</b><br>';
    if (pr.address)  html += '\uD83D\uDDA5\uFE0F \u0633\u0631\u0648\u0631: <b>' + pr.address + ':' + pr.port + '</b><br>';
    if (pr.remarks)  html += '\uD83D\uDCDD \u0646\u0627\u0645: <b>' + pr.remarks + '</b><br>';
    if (pr.tls)      html += '\uD83D\uDD12 TLS: <b>' + pr.tls + '</b><br>';
    if (pr.network)  html += '\uD83C\uDF10 Network: <b>' + pr.network + '</b><br>';
    if (m.deviceLock) html += '\uD83D\uDD12 Device Lock: <b>' + m.deviceLock + '</b><br>';
    if (m.userAgent)  html += '\uD83D\uDCF1 User-Agent: <b>' + m.userAgent + '</b><br>';
    if (m.createdAt)  html += '\uD83D\uDCC5 \u0633\u0627\u062E\u062A\u0647 \u0634\u062F\u0647: <b>' + m.createdAt + '</b><br>';
    if (!html) html = '<span style="color:var(--tx2)">\u0627\u0637\u0644\u0627\u0639\u0627\u062A \u0627\u0636\u0627\u0641\u06CC\u06CC \u06CC\u0627\u0641\u062A \u0646\u0634\u062F</span>';
    document.getElementById('impMeta').innerHTML = html;
  })
  .catch(function(e) { showSt('impSt', 'er', '\u062E\u0637\u0627: ' + e.message); });
}

function copyImpLink() {
  var v = document.getElementById('impLink').value;
  if (!v) return;
  navigator.clipboard.writeText(v).then(function() { showSt('impSt', 'ok', '\u0644\u06CC\u0646\u06A9 \u06A9\u067E\u06CC \u0634\u062F \u2713'); });
}

function copyImp() {
  var v = document.getElementById('impOut').value;
  if (!v) return;
  navigator.clipboard.writeText(v).then(function() { showSt('impSt', 'ok', '\u06A9\u067E\u06CC \u0634\u062F \u2713'); });
}

function dlImp() {
  dl(lastImpRaw, 'imported.npvt');
}
</script>
</body>
</html>`)
