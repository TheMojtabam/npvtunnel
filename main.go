package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ─── Data Structures ────────────────────────────────────────────────────────

type V2RayConfig struct {
	// Common
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	UUID     string `json:"uuid"`
	Remarks  string `json:"remarks"`

	// VMess specific
	AlterId  int    `json:"alterId"`
	Security string `json:"security"`
	Network  string `json:"network"`
	Path     string `json:"path"`
	Host     string `json:"host"`
	TLS      string `json:"tls"`
	SNI      string `json:"sni"`
	ALPN     string `json:"alpn"`

	// VLess specific
	Flow        string `json:"flow"`
	Fingerprint string `json:"fingerprint"`
	PublicKey   string `json:"publicKey"`
	ShortID     string `json:"shortId"`
	SpiderX     string `json:"spiderX"`

	// Trojan specific
	Password string `json:"password"`
}

type NapsterConfig struct {
	// Identity
	DeviceID    string `json:"deviceId"`
	DeviceName  string `json:"deviceName"`
	UserAgent   string `json:"userAgent"`

	// Proxy
	V2RayLink   string `json:"v2rayLink"`
	ParsedProxy V2RayConfig `json:"parsedProxy"`

	// Security
	EnableTLS13     bool   `json:"enableTls13"`
	EnableUTLS      bool   `json:"enableUtls"`
	UTLSFingerprint string `json:"utlsFingerprint"`
	EnableMux       bool   `json:"enableMux"`
	MuxConcurrency  int    `json:"muxConcurrency"`
	MuxProtocol     string `json:"muxProtocol"`

	// DNS
	DNSMode       string `json:"dnsMode"`
	DNSServer     string `json:"dnsServer"`
	FallbackDNS   string `json:"fallbackDNS"`
	EnableDNSoTLS bool   `json:"enableDnsoTls"`
	DNSoTLSServer string `json:"dnsotlsServer"`

	// Routing
	EnableBypass    bool     `json:"enableBypass"`
	BypassList      []string `json:"bypassList"`
	EnableBlock     bool     `json:"enableBlock"`
	BlockList       []string `json:"blockList"`
	ProxyMode       string   `json:"proxyMode"`

	// Advanced
	MTU             int    `json:"mtu"`
	UDPTimeout      int    `json:"udpTimeout"`
	EnableIPv6      bool   `json:"enableIpv6"`
	EnableSniffing  bool   `json:"enableSniffing"`
	SniffingDomains bool   `json:"sniffingDomains"`
	LogLevel        string `json:"logLevel"`

	// Output
	OutputFormat string `json:"outputFormat"`
}

type GenerateResponse struct {
	Success bool   `json:"success"`
	Config  string `json:"config"`
	Error   string `json:"error,omitempty"`
}

// ─── V2Ray Link Parser ───────────────────────────────────────────────────────

func parseVMessLink(link string) (V2RayConfig, error) {
	b64 := strings.TrimPrefix(link, "vmess://")
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(b64)
		if err != nil {
			return V2RayConfig{}, fmt.Errorf("invalid vmess base64")
		}
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(decoded, &raw); err != nil {
		return V2RayConfig{}, fmt.Errorf("invalid vmess JSON")
	}

	cfg := V2RayConfig{Protocol: "vmess"}
	if v, ok := raw["add"].(string); ok { cfg.Address = v }
	if v, ok := raw["port"]; ok { cfg.Port = toInt(v) }
	if v, ok := raw["id"].(string); ok { cfg.UUID = v }
	if v, ok := raw["aid"]; ok { cfg.AlterId = toInt(v) }
	if v, ok := raw["scy"].(string); ok { cfg.Security = v }
	if v, ok := raw["net"].(string); ok { cfg.Network = v }
	if v, ok := raw["path"].(string); ok { cfg.Path = v }
	if v, ok := raw["host"].(string); ok { cfg.Host = v }
	if v, ok := raw["tls"].(string); ok { cfg.TLS = v }
	if v, ok := raw["sni"].(string); ok { cfg.SNI = v }
	if v, ok := raw["alpn"].(string); ok { cfg.ALPN = v }
	if v, ok := raw["ps"].(string); ok { cfg.Remarks = v }
	if cfg.Security == "" { cfg.Security = "auto" }

	return cfg, nil
}

func parseVLessLink(link string) (V2RayConfig, error) {
	u, err := url.Parse(link)
	if err != nil { return V2RayConfig{}, err }

	cfg := V2RayConfig{Protocol: "vless"}
	cfg.UUID = u.User.Username()
	cfg.Address = u.Hostname()
	cfg.Port = toIntStr(u.Port())
	cfg.Remarks, _ = url.QueryUnescape(u.Fragment)

	q := u.Query()
	cfg.Network = q.Get("type")
	cfg.TLS = q.Get("security")
	cfg.SNI = q.Get("sni")
	cfg.ALPN = q.Get("alpn")
	cfg.Path = q.Get("path")
	cfg.Host = q.Get("host")
	cfg.Flow = q.Get("flow")
	cfg.Fingerprint = q.Get("fp")
	cfg.PublicKey = q.Get("pbk")
	cfg.ShortID = q.Get("sid")
	cfg.SpiderX = q.Get("spx")

	return cfg, nil
}

func parseTrojanLink(link string) (V2RayConfig, error) {
	u, err := url.Parse(link)
	if err != nil { return V2RayConfig{}, err }

	cfg := V2RayConfig{Protocol: "trojan"}
	cfg.Password = u.User.Username()
	cfg.Address = u.Hostname()
	cfg.Port = toIntStr(u.Port())
	cfg.Remarks, _ = url.QueryUnescape(u.Fragment)

	q := u.Query()
	cfg.Network = q.Get("type")
	cfg.TLS = q.Get("security")
	cfg.SNI = q.Get("sni")
	cfg.ALPN = q.Get("alpn")
	cfg.Path = q.Get("path")
	cfg.Host = q.Get("host")
	cfg.Fingerprint = q.Get("fp")

	return cfg, nil
}

func parseV2RayLink(link string) (V2RayConfig, error) {
	link = strings.TrimSpace(link)
	switch {
	case strings.HasPrefix(link, "vmess://"):
		return parseVMessLink(link)
	case strings.HasPrefix(link, "vless://"):
		return parseVLessLink(link)
	case strings.HasPrefix(link, "trojan://"):
		return parseTrojanLink(link)
	default:
		return V2RayConfig{}, fmt.Errorf("unsupported protocol")
	}
}

// ─── Config Generators ──────────────────────────────────────────────────────

func generateNapsterConfig(cfg NapsterConfig) string {
	p := cfg.ParsedProxy
	var sb strings.Builder

	sb.WriteString("# ╔══════════════════════════════════════╗\n")
	sb.WriteString("# ║     Napster Config Generator         ║\n")
	sb.WriteString(fmt.Sprintf("# ║     Generated: %-21s║\n", time.Now().Format("2006-01-02 15:04")))
	sb.WriteString("# ╚══════════════════════════════════════╝\n\n")

	// Mixed port
	sb.WriteString("mixed-port: 7890\n")
	sb.WriteString("socks-port: 7891\n")
	sb.WriteString("port: 7892\n")
	sb.WriteString("redir-port: 7893\n\n")

	// Modes
	sb.WriteString(fmt.Sprintf("mode: %s\n", cfg.ProxyMode))
	sb.WriteString("allow-lan: false\n")
	sb.WriteString(fmt.Sprintf("log-level: %s\n", cfg.LogLevel))
	sb.WriteString(fmt.Sprintf("ipv6: %v\n\n", cfg.EnableIPv6))

	// External controller
	sb.WriteString("external-controller: 127.0.0.1:9090\n")
	sb.WriteString("secret: \"\"\n\n")

	// DNS
	sb.WriteString("dns:\n")
	sb.WriteString("  enable: true\n")
	sb.WriteString(fmt.Sprintf("  enhanced-mode: %s\n", cfg.DNSMode))
	sb.WriteString("  listen: 0.0.0.0:53\n")
	sb.WriteString(fmt.Sprintf("  use-hosts: true\n"))
	if cfg.EnableDNSoTLS {
		sb.WriteString(fmt.Sprintf("  default-nameserver:\n    - %s\n", cfg.DNSServer))
		sb.WriteString(fmt.Sprintf("  nameserver:\n    - tls://%s\n", cfg.DNSoTLSServer))
	} else {
		sb.WriteString(fmt.Sprintf("  nameserver:\n    - %s\n", cfg.DNSServer))
	}
	sb.WriteString(fmt.Sprintf("  fallback:\n    - %s\n", cfg.FallbackDNS))
	sb.WriteString("  fallback-filter:\n    geoip: true\n    geoip-code: IR\n\n")

	// Sniffing
	if cfg.EnableSniffing {
		sb.WriteString("sniffer:\n  enable: true\n")
		if cfg.SniffingDomains {
			sb.WriteString("  sniff:\n    HTTP:\n      ports: [80]\n    TLS:\n      ports: [443]\n\n")
		}
	}

	// Proxy
	sb.WriteString("proxies:\n")
	proxyName := p.Remarks
	if proxyName == "" { proxyName = "proxy-1" }

	switch p.Protocol {
	case "vmess":
		sb.WriteString(fmt.Sprintf("  - name: \"%s\"\n", proxyName))
		sb.WriteString("    type: vmess\n")
		sb.WriteString(fmt.Sprintf("    server: %s\n", p.Address))
		sb.WriteString(fmt.Sprintf("    port: %d\n", p.Port))
		sb.WriteString(fmt.Sprintf("    uuid: %s\n", p.UUID))
		sb.WriteString(fmt.Sprintf("    alterId: %d\n", p.AlterId))
		sb.WriteString(fmt.Sprintf("    cipher: %s\n", p.Security))
		if p.TLS == "tls" {
			sb.WriteString("    tls: true\n")
			if p.SNI != "" { sb.WriteString(fmt.Sprintf("    servername: %s\n", p.SNI)) }
			if cfg.EnableTLS13 { sb.WriteString("    skip-cert-verify: false\n") }
		}
		if p.Network != "" && p.Network != "tcp" {
			sb.WriteString(fmt.Sprintf("    network: %s\n", p.Network))
			if p.Network == "ws" {
				sb.WriteString(fmt.Sprintf("    ws-opts:\n      path: %s\n", p.Path))
				if p.Host != "" { sb.WriteString(fmt.Sprintf("      headers:\n        Host: %s\n", p.Host)) }
			}
			if p.Network == "grpc" {
				sb.WriteString(fmt.Sprintf("    grpc-opts:\n      grpc-service-name: %s\n", p.Path))
			}
		}
		if cfg.EnableUTLS {
			sb.WriteString(fmt.Sprintf("    client-fingerprint: %s\n", cfg.UTLSFingerprint))
		}

	case "vless":
		sb.WriteString(fmt.Sprintf("  - name: \"%s\"\n", proxyName))
		sb.WriteString("    type: vless\n")
		sb.WriteString(fmt.Sprintf("    server: %s\n", p.Address))
		sb.WriteString(fmt.Sprintf("    port: %d\n", p.Port))
		sb.WriteString(fmt.Sprintf("    uuid: %s\n", p.UUID))
		sb.WriteString("    udp: true\n")
		if p.Flow != "" { sb.WriteString(fmt.Sprintf("    flow: %s\n", p.Flow)) }
		if p.TLS == "tls" || p.TLS == "reality" {
			sb.WriteString("    tls: true\n")
			if p.SNI != "" { sb.WriteString(fmt.Sprintf("    servername: %s\n", p.SNI)) }
		}
		if p.TLS == "reality" {
			sb.WriteString("    reality-opts:\n")
			sb.WriteString(fmt.Sprintf("      public-key: %s\n", p.PublicKey))
			if p.ShortID != "" { sb.WriteString(fmt.Sprintf("      short-id: %s\n", p.ShortID)) }
		}
		if p.Fingerprint != "" || cfg.EnableUTLS {
			fp := p.Fingerprint
			if fp == "" { fp = cfg.UTLSFingerprint }
			sb.WriteString(fmt.Sprintf("    client-fingerprint: %s\n", fp))
		}
		if p.Network != "" && p.Network != "tcp" {
			sb.WriteString(fmt.Sprintf("    network: %s\n", p.Network))
			if p.Network == "ws" {
				sb.WriteString(fmt.Sprintf("    ws-opts:\n      path: %s\n", p.Path))
				if p.Host != "" { sb.WriteString(fmt.Sprintf("      headers:\n        Host: %s\n", p.Host)) }
			}
			if p.Network == "grpc" {
				sb.WriteString(fmt.Sprintf("    grpc-opts:\n      grpc-service-name: %s\n", p.Path))
			}
		}

	case "trojan":
		sb.WriteString(fmt.Sprintf("  - name: \"%s\"\n", proxyName))
		sb.WriteString("    type: trojan\n")
		sb.WriteString(fmt.Sprintf("    server: %s\n", p.Address))
		sb.WriteString(fmt.Sprintf("    port: %d\n", p.Port))
		sb.WriteString(fmt.Sprintf("    password: %s\n", p.Password))
		if p.SNI != "" { sb.WriteString(fmt.Sprintf("    sni: %s\n", p.SNI)) }
		sb.WriteString("    udp: true\n")
		sb.WriteString("    skip-cert-verify: false\n")
		if cfg.EnableUTLS {
			sb.WriteString(fmt.Sprintf("    client-fingerprint: %s\n", cfg.UTLSFingerprint))
		}
		if p.Network == "ws" || p.Network == "grpc" {
			sb.WriteString(fmt.Sprintf("    network: %s\n", p.Network))
			if p.Network == "ws" {
				sb.WriteString(fmt.Sprintf("    ws-opts:\n      path: %s\n", p.Path))
			}
		}
	}

	// Mux
	if cfg.EnableMux {
		sb.WriteString(fmt.Sprintf("    smux:\n      enabled: true\n      protocol: %s\n      max-streams: %d\n",
			cfg.MuxProtocol, cfg.MuxConcurrency))
	}

	// Proxy Groups
	sb.WriteString("\nproxy-groups:\n")
	sb.WriteString(fmt.Sprintf("  - name: \"🚀 Proxy\"\n    type: select\n    proxies:\n      - \"%s\"\n      - DIRECT\n\n", proxyName))
	sb.WriteString("  - name: \"🌐 Auto\"\n    type: url-test\n    proxies:\n")
	sb.WriteString(fmt.Sprintf("      - \"%s\"\n", proxyName))
	sb.WriteString("    url: http://www.gstatic.com/generate_204\n    interval: 300\n\n")

	// Rules
	sb.WriteString("rules:\n")

	if cfg.EnableBlock && len(cfg.BlockList) > 0 {
		for _, d := range cfg.BlockList {
			d = strings.TrimSpace(d)
			if d != "" { sb.WriteString(fmt.Sprintf("  - DOMAIN-SUFFIX,%s,REJECT\n", d)) }
		}
	}

	if cfg.EnableBypass && len(cfg.BypassList) > 0 {
		for _, d := range cfg.BypassList {
			d = strings.TrimSpace(d)
			if d != "" { sb.WriteString(fmt.Sprintf("  - DOMAIN-SUFFIX,%s,DIRECT\n", d)) }
		}
	}

	// Default Iran bypass
	sb.WriteString("  - GEOIP,IR,DIRECT\n")
	sb.WriteString("  - GEOIP,private,DIRECT\n")
	sb.WriteString("  - MATCH,🚀 Proxy\n")

	// Device ID comment
	if cfg.DeviceID != "" {
		header := fmt.Sprintf("\n# device-id: %s\n# device-name: %s\n# user-agent: %s\n",
			cfg.DeviceID, cfg.DeviceName, cfg.UserAgent)
		return header + sb.String()
	}

	return sb.String()
}

func generateClashConfig(cfg NapsterConfig) string {
	// Same as napster but with clash header
	return "# Clash Config\n" + generateNapsterConfig(cfg)
}

// ─── HTTP Handlers ──────────────────────────────────────────────────────────

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, indexHTML)
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	link := r.URL.Query().Get("link")
	if link == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "no link provided"})
		return
	}

	cfg, err := parseV2RayLink(link)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(cfg)
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var cfg NapsterConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		json.NewEncoder(w).Encode(GenerateResponse{Error: "invalid JSON: " + err.Error()})
		return
	}

	// Parse the v2ray link
	if cfg.V2RayLink != "" {
		parsed, err := parseV2RayLink(cfg.V2RayLink)
		if err != nil {
			json.NewEncoder(w).Encode(GenerateResponse{Error: "link parse error: " + err.Error()})
			return
		}
		cfg.ParsedProxy = parsed
	}

	// Auto-generate device ID if empty
	if cfg.DeviceID == "" {
		cfg.DeviceID = uuid.New().String()
	}

	var result string
	switch cfg.OutputFormat {
	case "clash":
		result = generateClashConfig(cfg)
	default:
		result = generateNapsterConfig(cfg)
	}

	json.NewEncoder(w).Encode(GenerateResponse{Success: true, Config: result})
}

func handleNewDeviceID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"deviceId": uuid.New().String()})
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func toInt(v interface{}) int {
	switch val := v.(type) {
	case float64: return int(val)
	case string:
		var i int
		fmt.Sscanf(val, "%d", &i)
		return i
	}
	return 0
}

func toIntStr(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

// ─── Main ────────────────────────────────────────────────────────────────────

func main() {
	port := "8080"
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/parse", handleParse)
	mux.HandleFunc("/api/generate", handleGenerate)
	mux.HandleFunc("/api/device-id", handleNewDeviceID)

	addr := ":" + port
	log.Printf("🚀 Napster Config Generator running at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// ─── Embedded HTML ───────────────────────────────────────────────────────────

const indexHTML = `<!DOCTYPE html>
<html lang="fa" dir="rtl">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Napster Config Generator</title>
<style>
  @import url('https://fonts.googleapis.com/css2?family=Vazirmatn:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;600&display=swap');

  :root {
    --bg: #0a0e1a;
    --surface: #111827;
    --surface2: #1a2236;
    --border: #1e2d4a;
    --accent: #3b82f6;
    --accent2: #6366f1;
    --green: #10b981;
    --red: #ef4444;
    --yellow: #f59e0b;
    --text: #e2e8f0;
    --text2: #94a3b8;
    --glow: 0 0 20px rgba(59,130,246,0.3);
  }

  * { margin:0; padding:0; box-sizing:border-box; }

  body {
    background: var(--bg);
    color: var(--text);
    font-family: 'Vazirmatn', sans-serif;
    min-height: 100vh;
    background-image:
      radial-gradient(ellipse at 20% 50%, rgba(59,130,246,0.06) 0%, transparent 50%),
      radial-gradient(ellipse at 80% 20%, rgba(99,102,241,0.06) 0%, transparent 50%);
  }

  header {
    padding: 20px 40px;
    border-bottom: 1px solid var(--border);
    display: flex;
    align-items: center;
    gap: 16px;
    background: rgba(17,24,39,0.8);
    backdrop-filter: blur(12px);
    position: sticky;
    top: 0;
    z-index: 100;
  }

  .logo {
    width: 40px; height: 40px;
    background: linear-gradient(135deg, var(--accent), var(--accent2));
    border-radius: 10px;
    display: flex; align-items: center; justify-content: center;
    font-size: 20px;
    box-shadow: var(--glow);
  }

  header h1 { font-size: 20px; font-weight: 700; letter-spacing: 0.5px; }
  header p { font-size: 13px; color: var(--text2); margin-top: 2px; }

  .badge {
    margin-right: auto;
    background: rgba(59,130,246,0.15);
    border: 1px solid rgba(59,130,246,0.3);
    color: var(--accent);
    padding: 4px 12px;
    border-radius: 20px;
    font-size: 12px;
    font-weight: 600;
  }

  main {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 24px;
    padding: 32px 40px;
    max-width: 1600px;
    margin: 0 auto;
  }

  @media (max-width: 1024px) { main { grid-template-columns: 1fr; } }

  .panel {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 16px;
    overflow: hidden;
  }

  .panel-header {
    padding: 16px 24px;
    border-bottom: 1px solid var(--border);
    display: flex;
    align-items: center;
    gap: 10px;
    font-weight: 600;
    font-size: 15px;
    background: var(--surface2);
  }

  .panel-icon { font-size: 18px; }

  .panel-body { padding: 24px; display: flex; flex-direction: column; gap: 20px; }

  .section-title {
    font-size: 12px;
    font-weight: 600;
    color: var(--accent);
    text-transform: uppercase;
    letter-spacing: 1px;
    margin-bottom: 12px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .section-title::after {
    content: '';
    flex: 1;
    height: 1px;
    background: var(--border);
  }

  .field { display: flex; flex-direction: column; gap: 6px; }
  .field label { font-size: 13px; color: var(--text2); font-weight: 500; }

  input[type=text], input[type=number], select, textarea {
    background: var(--bg);
    border: 1px solid var(--border);
    color: var(--text);
    padding: 10px 14px;
    border-radius: 10px;
    font-family: inherit;
    font-size: 14px;
    transition: border-color 0.2s, box-shadow 0.2s;
    width: 100%;
    direction: ltr;
  }

  input:focus, select:focus, textarea:focus {
    outline: none;
    border-color: var(--accent);
    box-shadow: 0 0 0 3px rgba(59,130,246,0.15);
  }

  textarea { resize: vertical; min-height: 80px; }

  .input-with-btn { display: flex; gap: 8px; }
  .input-with-btn input { flex: 1; }

  .grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
  .grid-3 { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 14px; }

  .toggle-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 14px;
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 10px;
  }

  .toggle-label { font-size: 13px; font-weight: 500; }
  .toggle-desc { font-size: 11px; color: var(--text2); margin-top: 2px; }

  .toggle {
    position: relative;
    width: 44px; height: 24px;
    flex-shrink: 0;
  }

  .toggle input { opacity: 0; width: 0; height: 0; }

  .toggle-slider {
    position: absolute; inset: 0;
    background: var(--border);
    border-radius: 24px;
    cursor: pointer;
    transition: 0.3s;
  }

  .toggle-slider::before {
    content: '';
    position: absolute;
    width: 18px; height: 18px;
    left: 3px; top: 3px;
    background: white;
    border-radius: 50%;
    transition: 0.3s;
  }

  .toggle input:checked + .toggle-slider { background: var(--accent); }
  .toggle input:checked + .toggle-slider::before { transform: translateX(20px); }

  .btn {
    padding: 10px 20px;
    border-radius: 10px;
    border: none;
    cursor: pointer;
    font-family: inherit;
    font-size: 14px;
    font-weight: 600;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 8px;
    justify-content: center;
  }

  .btn-primary {
    background: linear-gradient(135deg, var(--accent), var(--accent2));
    color: white;
    width: 100%;
    padding: 14px;
    font-size: 15px;
    box-shadow: 0 4px 15px rgba(59,130,246,0.3);
  }

  .btn-primary:hover { transform: translateY(-1px); box-shadow: 0 6px 20px rgba(59,130,246,0.4); }
  .btn-primary:active { transform: translateY(0); }

  .btn-sm {
    background: var(--surface2);
    border: 1px solid var(--border);
    color: var(--text2);
    padding: 8px 14px;
    font-size: 12px;
    flex-shrink: 0;
  }

  .btn-sm:hover { border-color: var(--accent); color: var(--accent); }

  .btn-copy {
    background: rgba(16,185,129,0.15);
    border: 1px solid rgba(16,185,129,0.3);
    color: var(--green);
    padding: 8px 16px;
    font-size: 13px;
  }

  .btn-copy:hover { background: rgba(16,185,129,0.25); }

  .btn-dl {
    background: rgba(99,102,241,0.15);
    border: 1px solid rgba(99,102,241,0.3);
    color: var(--accent2);
    padding: 8px 16px;
    font-size: 13px;
  }

  .parsed-info {
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 14px;
    font-family: 'JetBrains Mono', monospace;
    font-size: 12px;
    color: var(--text2);
    display: none;
  }

  .parsed-info.show { display: block; }
  .parsed-info .proto-badge {
    display: inline-block;
    background: rgba(59,130,246,0.2);
    color: var(--accent);
    border-radius: 6px;
    padding: 2px 8px;
    font-size: 11px;
    margin-bottom: 8px;
    font-weight: 600;
  }

  .output-area {
    position: relative;
    flex: 1;
  }

  .output-textarea {
    width: 100%;
    height: 500px;
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 16px;
    font-family: 'JetBrains Mono', monospace;
    font-size: 12px;
    color: #a8ff78;
    line-height: 1.7;
    resize: none;
    direction: ltr;
  }

  .output-actions {
    display: flex;
    gap: 8px;
    margin-top: 12px;
  }

  .status-bar {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 14px;
    border-radius: 10px;
    font-size: 13px;
    font-weight: 500;
  }

  .status-bar.success { background: rgba(16,185,129,0.1); border: 1px solid rgba(16,185,129,0.3); color: var(--green); }
  .status-bar.error { background: rgba(239,68,68,0.1); border: 1px solid rgba(239,68,68,0.3); color: var(--red); }
  .status-bar.hidden { display: none; }

  .dot { width: 8px; height: 8px; border-radius: 50%; background: currentColor; }

  ::-webkit-scrollbar { width: 6px; height: 6px; }
  ::-webkit-scrollbar-track { background: var(--bg); }
  ::-webkit-scrollbar-thumb { background: var(--border); border-radius: 3px; }
</style>
</head>
<body>
<header>
  <div class="logo">⚡</div>
  <div>
    <h1>Napster Config Generator</h1>
    <p>تولید کانفیگ حرفه‌ای با تمام تنظیمات امنیتی</p>
  </div>
  <div class="badge">v2.0.0</div>
</header>

<main>
  <!-- Left Panel: Settings -->
  <div style="display:flex;flex-direction:column;gap:20px;">

    <!-- V2Ray Link -->
    <div class="panel">
      <div class="panel-header"><span class="panel-icon">🔗</span> لینک V2Ray</div>
      <div class="panel-body">
        <div class="field">
          <label>لینک vmess / vless / trojan</label>
          <div class="input-with-btn">
            <input type="text" id="v2rayLink" placeholder="vmess://... یا vless://... یا trojan://..." oninput="onLinkChange()">
            <button class="btn btn-sm" onclick="parseLink()">📡 پارس</button>
          </div>
        </div>
        <div class="parsed-info" id="parsedInfo"></div>
      </div>
    </div>

    <!-- Device & Identity -->
    <div class="panel">
      <div class="panel-header"><span class="panel-icon">📱</span> هویت دستگاه</div>
      <div class="panel-body">
        <div class="grid-2">
          <div class="field">
            <label>Device ID</label>
            <div class="input-with-btn">
              <input type="text" id="deviceId" placeholder="خودکار تولید می‌شود">
              <button class="btn btn-sm" onclick="genDeviceId()">🔄</button>
            </div>
          </div>
          <div class="field">
            <label>Device Name</label>
            <input type="text" id="deviceName" placeholder="مثال: iPhone 15 Pro" value="Napster Client">
          </div>
        </div>
        <div class="field">
          <label>User Agent</label>
          <input type="text" id="userAgent" placeholder="User-Agent سفارشی" value="Napster/2.0 (iOS 17.0)">
        </div>
      </div>
    </div>

    <!-- Security Settings -->
    <div class="panel">
      <div class="panel-header"><span class="panel-icon">🔒</span> تنظیمات امنیتی</div>
      <div class="panel-body">
        <div class="section-title">TLS & Fingerprint</div>

        <div class="toggle-row">
          <div>
            <div class="toggle-label">فعال‌سازی TLS 1.3</div>
            <div class="toggle-desc">استفاده از TLS 1.3 برای امنیت بیشتر</div>
          </div>
          <label class="toggle"><input type="checkbox" id="enableTls13" checked><span class="toggle-slider"></span></label>
        </div>

        <div class="toggle-row">
          <div>
            <div class="toggle-label">uTLS / Browser Fingerprint</div>
            <div class="toggle-desc">جعل fingerprint مرورگر برای دور زدن DPI</div>
          </div>
          <label class="toggle"><input type="checkbox" id="enableUtls" checked onchange="toggleUtls()"><span class="toggle-slider"></span></label>
        </div>

        <div class="field" id="utlsField">
          <label>نوع Fingerprint</label>
          <select id="utlsFingerprint">
            <option value="chrome">Chrome (توصیه شده)</option>
            <option value="firefox">Firefox</option>
            <option value="safari">Safari</option>
            <option value="ios">iOS Safari</option>
            <option value="android">Android Chrome</option>
            <option value="edge">Microsoft Edge</option>
            <option value="360">360 Browser</option>
            <option value="qq">QQ Browser</option>
            <option value="random">Random (تصادفی)</option>
          </select>
        </div>

        <div class="section-title" style="margin-top:8px;">Mux / Multiplexing</div>

        <div class="toggle-row">
          <div>
            <div class="toggle-label">فعال‌سازی Mux</div>
            <div class="toggle-desc">ادغام چند اتصال در یک تانل</div>
          </div>
          <label class="toggle"><input type="checkbox" id="enableMux" onchange="toggleMux()"><span class="toggle-slider"></span></label>
        </div>

        <div id="muxFields" style="display:none;" class="grid-2">
          <div class="field">
            <label>Mux Protocol</label>
            <select id="muxProtocol">
              <option value="smux">smux</option>
              <option value="yamux">yamux</option>
              <option value="h2mux">h2mux</option>
            </select>
          </div>
          <div class="field">
            <label>Max Concurrency</label>
            <input type="number" id="muxConcurrency" value="8" min="1" max="64">
          </div>
        </div>
      </div>
    </div>

    <!-- DNS Settings -->
    <div class="panel">
      <div class="panel-header"><span class="panel-icon">🌐</span> تنظیمات DNS</div>
      <div class="panel-body">
        <div class="grid-2">
          <div class="field">
            <label>DNS Mode</label>
            <select id="dnsMode">
              <option value="fake-ip">Fake-IP (توصیه شده)</option>
              <option value="redir-host">Redir-Host</option>
              <option value="normal">Normal</option>
            </select>
          </div>
          <div class="field">
            <label>DNS Server اصلی</label>
            <select id="dnsServer">
              <option value="1.1.1.1">Cloudflare (1.1.1.1)</option>
              <option value="8.8.8.8">Google (8.8.8.8)</option>
              <option value="9.9.9.9">Quad9 (9.9.9.9)</option>
              <option value="208.67.222.222">OpenDNS</option>
              <option value="custom">سفارشی...</option>
            </select>
          </div>
        </div>

        <div class="field">
          <label>DNS Fallback</label>
          <select id="fallbackDns">
            <option value="8.8.8.8">Google (8.8.8.8)</option>
            <option value="1.1.1.1">Cloudflare (1.1.1.1)</option>
            <option value="208.67.220.220">OpenDNS Secondary</option>
          </select>
        </div>

        <div class="toggle-row">
          <div>
            <div class="toggle-label">DNS over TLS</div>
            <div class="toggle-desc">رمزگذاری درخواست‌های DNS</div>
          </div>
          <label class="toggle"><input type="checkbox" id="enableDnsoTls" onchange="toggleDoT()"><span class="toggle-slider"></span></label>
        </div>

        <div class="field" id="dotField" style="display:none;">
          <label>DoT Server</label>
          <select id="dnsotlsServer">
            <option value="cloudflare-dns.com">cloudflare-dns.com</option>
            <option value="dns.google">dns.google</option>
            <option value="dns.quad9.net">dns.quad9.net</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Routing Settings -->
    <div class="panel">
      <div class="panel-header"><span class="panel-icon">🗺️</span> مسیریابی</div>
      <div class="panel-body">
        <div class="grid-2">
          <div class="field">
            <label>Proxy Mode</label>
            <select id="proxyMode">
              <option value="rule">Rule (هوشمند)</option>
              <option value="global">Global (همه از پروکسی)</option>
              <option value="direct">Direct (همه مستقیم)</option>
            </select>
          </div>
          <div class="field">
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

        <div class="toggle-row">
          <div>
            <div class="toggle-label">Bypass (دور زدن) برای ایران</div>
            <div class="toggle-desc">سایت‌های ایرانی مستقیم متصل شوند</div>
          </div>
          <label class="toggle"><input type="checkbox" id="enableBypass" checked onchange="toggleBypass()"><span class="toggle-slider"></span></label>
        </div>

        <div class="field" id="bypassField">
          <label>لیست Bypass (هر خط یک دامنه)</label>
          <textarea id="bypassList" rows="4" placeholder="ir&#10;digikala.com&#10;aparat.com">ir
shaparak.ir
digikala.com
aparat.com
snapp.ir</textarea>
        </div>

        <div class="toggle-row">
          <div>
            <div class="toggle-label">Block (مسدودسازی)</div>
            <div class="toggle-desc">بلاک کردن دامنه‌های خاص</div>
          </div>
          <label class="toggle"><input type="checkbox" id="enableBlock" onchange="toggleBlock()"><span class="toggle-slider"></span></label>
        </div>

        <div class="field" id="blockField" style="display:none;">
          <label>لیست Block</label>
          <textarea id="blockList" rows="3" placeholder="ads.example.com&#10;tracker.com"></textarea>
        </div>
      </div>
    </div>

    <!-- Advanced -->
    <div class="panel">
      <div class="panel-header"><span class="panel-icon">⚙️</span> تنظیمات پیشرفته</div>
      <div class="panel-body">
        <div class="grid-3">
          <div class="field">
            <label>MTU</label>
            <input type="number" id="mtu" value="1500" min="576" max="9000">
          </div>
          <div class="field">
            <label>UDP Timeout (ثانیه)</label>
            <input type="number" id="udpTimeout" value="60" min="10" max="300">
          </div>
          <div class="field">
            <label>فرمت خروجی</label>
            <select id="outputFormat">
              <option value="napster">Napster</option>
              <option value="clash">Clash</option>
            </select>
          </div>
        </div>

        <div class="grid-2">
          <div class="toggle-row">
            <div>
              <div class="toggle-label">IPv6</div>
            </div>
            <label class="toggle"><input type="checkbox" id="enableIpv6"><span class="toggle-slider"></span></label>
          </div>
          <div class="toggle-row">
            <div>
              <div class="toggle-label">Sniffing</div>
            </div>
            <label class="toggle"><input type="checkbox" id="enableSniffing" checked><span class="toggle-slider"></span></label>
          </div>
        </div>

        <button class="btn btn-primary" onclick="generateConfig()">
          ⚡ تولید کانفیگ
        </button>
      </div>
    </div>

  </div>

  <!-- Right Panel: Output -->
  <div class="panel" style="position:sticky;top:90px;align-self:start;">
    <div class="panel-header"><span class="panel-icon">📄</span> کانفیگ تولید شده</div>
    <div class="panel-body">
      <div id="statusBar" class="status-bar hidden"></div>
      <div class="output-area">
        <textarea class="output-textarea" id="output" readonly placeholder="# کانفیگ اینجا نمایش داده می‌شود...&#10;# ابتدا لینک V2Ray وارد کنید و روی «تولید کانفیگ» کلیک کنید"></textarea>
      </div>
      <div class="output-actions">
        <button class="btn btn-copy" onclick="copyOutput()">📋 کپی</button>
        <button class="btn btn-dl" onclick="downloadOutput()">💾 دانلود</button>
        <button class="btn btn-sm" onclick="clearOutput()">🗑️ پاک کردن</button>
      </div>
    </div>
  </div>
</main>

<script>
function onLinkChange() {
  const link = document.getElementById('v2rayLink').value.trim();
  if (!link) { document.getElementById('parsedInfo').classList.remove('show'); return; }
  if (link.length > 20) parseLink();
}

async function parseLink() {
  const link = document.getElementById('v2rayLink').value.trim();
  if (!link) return;
  try {
    const res = await fetch('/api/parse?link=' + encodeURIComponent(link));
    const data = await res.json();
    if (data.error) { showParsed(null, data.error); return; }
    showParsed(data);
  } catch(e) { showParsed(null, e.message); }
}

function showParsed(data, error) {
  const el = document.getElementById('parsedInfo');
  if (error) {
    el.innerHTML = '<span style="color:var(--red)">❌ ' + error + '</span>';
  } else {
    const proto = data.protocol.toUpperCase();
    const server = data.address + ':' + data.port;
    const id = data.uuid || data.password || '-';
    const net = data.network || 'tcp';
    const tls = data.tls || 'none';
    el.innerHTML = '<span class="proto-badge">' + proto + '</span><br>' +
      '🖥️ Server: <b>' + server + '</b><br>' +
      '🔑 ID: <b>' + id.substring(0,20) + (id.length>20?'...':'') + '</b><br>' +
      '🌐 Network: <b>' + net + '</b> &nbsp; 🔒 TLS: <b>' + tls + '</b>' +
      (data.remarks ? '<br>📝 ' + data.remarks : '');
  }
  el.classList.add('show');
}

async function genDeviceId() {
  try {
    const res = await fetch('/api/device-id');
    const data = await res.json();
    document.getElementById('deviceId').value = data.deviceId;
  } catch(e) {}
}

function toggleUtls() {
  document.getElementById('utlsField').style.display =
    document.getElementById('enableUtls').checked ? 'flex' : 'none';
}

function toggleMux() {
  document.getElementById('muxFields').style.display =
    document.getElementById('enableMux').checked ? 'grid' : 'none';
}

function toggleBypass() {
  document.getElementById('bypassField').style.display =
    document.getElementById('enableBypass').checked ? 'flex' : 'none';
}

function toggleBlock() {
  document.getElementById('blockField').style.display =
    document.getElementById('enableBlock').checked ? 'flex' : 'none';
}

function toggleDoT() {
  document.getElementById('dotField').style.display =
    document.getElementById('enableDnsoTls').checked ? 'flex' : 'none';
}

function g(id) { return document.getElementById(id); }
function gv(id) { return g(id).value; }
function gc(id) { return g(id).checked; }

async function generateConfig() {
  const link = gv('v2rayLink').trim();
  if (!link) { showStatus('error', 'لطفاً ابتدا لینک V2Ray وارد کنید'); return; }

  const bypassRaw = gv('bypassList').split('\n').map(s=>s.trim()).filter(Boolean);
  const blockRaw  = gv('blockList').split('\n').map(s=>s.trim()).filter(Boolean);

  const payload = {
    v2rayLink:       link,
    deviceId:        gv('deviceId'),
    deviceName:      gv('deviceName'),
    userAgent:       gv('userAgent'),
    enableTls13:     gc('enableTls13'),
    enableUtls:      gc('enableUtls'),
    utlsFingerprint: gv('utlsFingerprint'),
    enableMux:       gc('enableMux'),
    muxConcurrency:  parseInt(gv('muxConcurrency')) || 8,
    muxProtocol:     gv('muxProtocol'),
    dnsMode:         gv('dnsMode'),
    dnsServer:       gv('dnsServer'),
    fallbackDNS:     gv('fallbackDns'),
    enableDnsoTls:   gc('enableDnsoTls'),
    dnsotlsServer:   gv('dnsotlsServer'),
    proxyMode:       gv('proxyMode'),
    logLevel:        gv('logLevel'),
    enableBypass:    gc('enableBypass'),
    bypassList:      bypassRaw,
    enableBlock:     gc('enableBlock'),
    blockList:       blockRaw,
    mtu:             parseInt(gv('mtu')) || 1500,
    udpTimeout:      parseInt(gv('udpTimeout')) || 60,
    enableIpv6:      gc('enableIpv6'),
    enableSniffing:  gc('enableSniffing'),
    sniffingDomains: gc('enableSniffing'),
    outputFormat:    gv('outputFormat'),
  };

  try {
    const res = await fetch('/api/generate', {
      method: 'POST',
      headers: {'Content-Type':'application/json'},
      body: JSON.stringify(payload)
    });
    const data = await res.json();
    if (data.error) { showStatus('error', data.error); return; }
    g('output').value = data.config;
    showStatus('success', 'کانفیگ با موفقیت تولید شد ✓');
  } catch(e) {
    showStatus('error', 'خطا در ارتباط با سرور: ' + e.message);
  }
}

function showStatus(type, msg) {
  const el = g('statusBar');
  el.className = 'status-bar ' + type;
  el.innerHTML = '<div class="dot"></div>' + msg;
}

function copyOutput() {
  const txt = g('output').value;
  if (!txt) return;
  navigator.clipboard.writeText(txt).then(() => showStatus('success', 'کپی شد ✓'));
}

function downloadOutput() {
  const txt = g('output').value;
  if (!txt) return;
  const fmt = gv('outputFormat');
  const blob = new Blob([txt], {type:'text/plain'});
  const a = document.createElement('a');
  a.href = URL.createObjectURL(blob);
  a.download = fmt + '-config-' + Date.now() + '.yaml';
  a.click();
}

function clearOutput() {
  g('output').value = '';
  g('statusBar').className = 'status-bar hidden';
  g('parsedInfo').classList.remove('show');
  g('v2rayLink').value = '';
}

// Init: generate device ID on load
genDeviceId();
</script>
</body>
</html>
`
