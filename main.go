package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
)

// ─────────────────────────────────────────────────────────────────────────────
// STRUCTS
// ─────────────────────────────────────────────────────────────────────────────

type V2RayConfig struct {
	Protocol    string `json:"protocol"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	UUID        string `json:"uuid"`
	Remarks     string `json:"remarks"`
	AlterId     int    `json:"alterId"`
	Security    string `json:"security"`
	Network     string `json:"network"`
	Path        string `json:"path"`
	Host        string `json:"host"`
	TLS         string `json:"tls"`
	SNI         string `json:"sni"`
	ALPN        string `json:"alpn"`
	Flow        string `json:"flow"`
	Fingerprint string `json:"fingerprint"`
	PublicKey   string `json:"publicKey"`
	ShortID     string `json:"shortId"`
	Password    string `json:"password"`
}

type Profile struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	CreatedAt   string          `json:"createdAt"`
	V2RayLink   string          `json:"v2rayLink"`
	ParsedProxy V2RayConfig     `json:"parsedProxy"`
	Settings    NapsterSettings `json:"settings"`
	Config      string          `json:"config"`
}

type NapsterSettings struct {
	EnableDeviceLock bool   `json:"enableDeviceLock"`
	DeviceID         string `json:"deviceId"`
	UserAgent        string `json:"userAgent"`
	EnablePassword   bool   `json:"enablePassword"`
	Password         string `json:"password"`
	EnableTLS13      bool   `json:"enableTls13"`
	EnableUTLS       bool   `json:"enableUtls"`
	UTLSFingerprint  string `json:"utlsFingerprint"`
	EnableMux        bool   `json:"enableMux"`
	MuxConcurrency   int    `json:"muxConcurrency"`
	MuxProtocol      string `json:"muxProtocol"`
	DNSMode          string `json:"dnsMode"`
	DNSServer        string `json:"dnsServer"`
	FallbackDNS      string `json:"fallbackDNS"`
	EnableDNSoTLS    bool   `json:"enableDnsoTls"`
	DNSoTLSServer    string `json:"dnsotlsServer"`
	EnableBypass     bool   `json:"enableBypass"`
	BypassDomains    []string `json:"bypassDomains"`
	BypassIPs        []string `json:"bypassIPs"`
	EnableBlock      bool   `json:"enableBlock"`
	BlockList        []string `json:"blockList"`
	ProxyMode        string `json:"proxyMode"`
	LogLevel         string `json:"logLevel"`
	MTU              int    `json:"mtu"`
	EnableIPv6       bool   `json:"enableIpv6"`
	EnableSniffing   bool   `json:"enableSniffing"`
	OutputFormat     string `json:"outputFormat"`
	Theme            string `json:"theme"`
}

type AppData struct {
	Profiles []Profile      `json:"profiles"`
	LastSettings NapsterSettings `json:"lastSettings"`
}

type GenerateRequest struct {
	V2RayLink   string          `json:"v2rayLink"`
	ProfileName string          `json:"profileName"`
	Settings    NapsterSettings `json:"settings"`
}

type GenerateResponse struct {
	Success   bool   `json:"success"`
	Config    string `json:"config"`
	Encrypted string `json:"encrypted,omitempty"`
	ProfileID string `json:"profileId,omitempty"`
	Error     string `json:"error,omitempty"`
}

type ImportResponse struct {
	Success   bool        `json:"success"`
	Decrypted string      `json:"decrypted,omitempty"`
	Info      interface{} `json:"info,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// ─────────────────────────────────────────────────────────────────────────────
// DATA PERSISTENCE
// ─────────────────────────────────────────────────────────────────────────────

func dataFilePath() string {
	exe, err := os.Executable()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(filepath.Dir(exe), "config.json")
}

func loadData() AppData {
	data, err := os.ReadFile(dataFilePath())
	if err != nil {
		return AppData{Profiles: []Profile{}, LastSettings: defaultSettings()}
	}
	var appData AppData
	if err := json.Unmarshal(data, &appData); err != nil {
		return AppData{Profiles: []Profile{}, LastSettings: defaultSettings()}
	}
	if appData.Profiles == nil {
		appData.Profiles = []Profile{}
	}
	return appData
}

func saveData(appData AppData) error {
	data, err := json.MarshalIndent(appData, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFilePath(), data, 0644)
}

func defaultSettings() NapsterSettings {
	return NapsterSettings{
		EnableTLS13:     true,
		EnableUTLS:      true,
		UTLSFingerprint: "chrome",
		MuxConcurrency:  8,
		MuxProtocol:     "smux",
		DNSMode:         "fake-ip",
		DNSServer:       "1.1.1.1",
		FallbackDNS:     "8.8.8.8",
		EnableBypass:    true,
		BypassDomains: []string{
			"ir", "shaparak.ir", "digikala.com", "aparat.com",
			"snapp.ir", "divar.ir", "torob.com", "fidibo.com",
		},
		BypassIPs: []string{
			"192.168.0.0/16",
			"10.0.0.0/8",
			"172.16.0.0/12",
			"127.0.0.0/8",
		},
		ProxyMode:      "rule",
		LogLevel:       "warning",
		MTU:            1500,
		EnableSniffing: true,
		OutputFormat:   "napster",
		Theme:          "dark",
		UserAgent:      "Napster/2.0",
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// ENCRYPTION
// ─────────────────────────────────────────────────────────────────────────────

func encryptConfig(plaintext, password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	combined := append(salt, ciphertext...)
	return "NAPSTER_ENC:" + base64.StdEncoding.EncodeToString(combined), nil
}

func decryptConfig(encrypted, password string) (string, error) {
	if !strings.HasPrefix(encrypted, "NAPSTER_ENC:") {
		return "", fmt.Errorf("not an encrypted napster config")
	}
	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(encrypted, "NAPSTER_ENC:"))
	if err != nil {
		return "", fmt.Errorf("invalid base64")
	}
	if len(data) < 16 {
		return "", fmt.Errorf("data too short")
	}
	salt := data[:16]
	ciphertext := data[16:]
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("wrong password or corrupted data")
	}
	return string(plaintext), nil
}

// ─────────────────────────────────────────────────────────────────────────────
// V2RAY PARSERS
// ─────────────────────────────────────────────────────────────────────────────

func parseVMessLink(link string) (V2RayConfig, error) {
	b64 := strings.TrimPrefix(link, "vmess://")
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(b64)
		if err != nil {
			// try url-safe
			decoded, err = base64.URLEncoding.DecodeString(b64)
			if err != nil {
				decoded, err = base64.RawURLEncoding.DecodeString(b64)
				if err != nil {
					return V2RayConfig{}, fmt.Errorf("invalid vmess base64")
				}
			}
		}
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(decoded, &raw); err != nil {
		return V2RayConfig{}, fmt.Errorf("invalid vmess JSON: %v", err)
	}
	cfg := V2RayConfig{Protocol: "vmess"}
	if v, ok := raw["add"].(string); ok  { cfg.Address = v }
	if v, ok := raw["port"]; ok           { cfg.Port = toInt(v) }
	if v, ok := raw["id"].(string); ok    { cfg.UUID = v }
	if v, ok := raw["aid"]; ok            { cfg.AlterId = toInt(v) }
	if v, ok := raw["scy"].(string); ok   { cfg.Security = v }
	if v, ok := raw["net"].(string); ok   { cfg.Network = v }
	if v, ok := raw["path"].(string); ok  { cfg.Path = v }
	if v, ok := raw["host"].(string); ok  { cfg.Host = v }
	if v, ok := raw["tls"].(string); ok   { cfg.TLS = v }
	if v, ok := raw["sni"].(string); ok   { cfg.SNI = v }
	if v, ok := raw["alpn"].(string); ok  { cfg.ALPN = v }
	if v, ok := raw["ps"].(string); ok    { cfg.Remarks = v }
	if v, ok := raw["fp"].(string); ok    { cfg.Fingerprint = v }
	if cfg.Security == ""                  { cfg.Security = "auto" }
	if cfg.Network == ""                   { cfg.Network = "tcp" }
	return cfg, nil
}

func parseVLessLink(link string) (V2RayConfig, error) {
	u, err := url.Parse(link)
	if err != nil { return V2RayConfig{}, err }
	cfg := V2RayConfig{Protocol: "vless"}
	cfg.UUID    = u.User.Username()
	cfg.Address = u.Hostname()
	cfg.Port    = toIntStr(u.Port())
	cfg.Remarks, _ = url.QueryUnescape(u.Fragment)
	q := u.Query()
	cfg.Network     = q.Get("type")
	cfg.TLS         = q.Get("security")
	cfg.SNI         = q.Get("sni")
	cfg.ALPN        = q.Get("alpn")
	cfg.Path        = q.Get("path")
	cfg.Host        = q.Get("host")
	cfg.Flow        = q.Get("flow")
	cfg.Fingerprint = q.Get("fp")
	cfg.PublicKey   = q.Get("pbk")
	cfg.ShortID     = q.Get("sid")
	if cfg.Network == "" { cfg.Network = "tcp" }
	return cfg, nil
}

func parseTrojanLink(link string) (V2RayConfig, error) {
	u, err := url.Parse(link)
	if err != nil { return V2RayConfig{}, err }
	cfg := V2RayConfig{Protocol: "trojan"}
	cfg.Password = u.User.Username()
	cfg.Address  = u.Hostname()
	cfg.Port     = toIntStr(u.Port())
	cfg.Remarks, _ = url.QueryUnescape(u.Fragment)
	q := u.Query()
	cfg.Network     = q.Get("type")
	cfg.TLS         = q.Get("security")
	cfg.SNI         = q.Get("sni")
	cfg.ALPN        = q.Get("alpn")
	cfg.Path        = q.Get("path")
	cfg.Host        = q.Get("host")
	cfg.Fingerprint = q.Get("fp")
	if cfg.Network == "" { cfg.Network = "tcp" }
	if cfg.TLS == ""     { cfg.TLS = "tls" }
	return cfg, nil
}

func parseV2RayLink(link string) (V2RayConfig, error) {
	link = strings.TrimSpace(link)
	switch {
	case strings.HasPrefix(link, "vmess://"):  return parseVMessLink(link)
	case strings.HasPrefix(link, "vless://"):  return parseVLessLink(link)
	case strings.HasPrefix(link, "trojan://"): return parseTrojanLink(link)
	default:
		return V2RayConfig{}, fmt.Errorf("پروتکل پشتیبانی نمیشه (فقط vmess/vless/trojan)")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// CONFIG GENERATOR
// ─────────────────────────────────────────────────────────────────────────────

func generateConfig(p V2RayConfig, s NapsterSettings) string {
	var sb strings.Builder
	now := time.Now().Format("2006-01-02 15:04")

	sb.WriteString("# ╔══════════════════════════════════════════╗\n")
	sb.WriteString("# ║      Napster Config Generator v3         ║\n")
	sb.WriteString(fmt.Sprintf("# ║      %s                   ║\n", now))
	sb.WriteString("# ╚══════════════════════════════════════════╝\n")

	if s.EnableDeviceLock && s.DeviceID != "" {
		sb.WriteString(fmt.Sprintf("# device-lock: %s\n", s.DeviceID))
	}
	if s.UserAgent != "" {
		sb.WriteString(fmt.Sprintf("# user-agent: %s\n", s.UserAgent))
	}
	sb.WriteString("\n")

	// Ports
	sb.WriteString("mixed-port: 7890\n")
	sb.WriteString("socks-port: 7891\n")
	sb.WriteString("port: 7892\n")
	sb.WriteString("redir-port: 7893\n")
	sb.WriteString("tproxy-port: 7894\n\n")

	// Mode
	sb.WriteString(fmt.Sprintf("mode: %s\n", s.ProxyMode))
	sb.WriteString("allow-lan: false\n")
	sb.WriteString(fmt.Sprintf("log-level: %s\n", s.LogLevel))
	sb.WriteString(fmt.Sprintf("ipv6: %v\n\n", s.EnableIPv6))

	// External controller
	sb.WriteString("external-controller: 127.0.0.1:9090\n")
	sb.WriteString("secret: \"\"\n\n")

	// Profile name
	proxyName := p.Remarks
	if proxyName == "" {
		proxyName = fmt.Sprintf("%s-%s:%d", strings.ToUpper(p.Protocol), p.Address, p.Port)
	}

	// DNS
	sb.WriteString("dns:\n")
	sb.WriteString("  enable: true\n")
	sb.WriteString(fmt.Sprintf("  enhanced-mode: %s\n", s.DNSMode))
	sb.WriteString("  listen: 0.0.0.0:53\n")
	sb.WriteString("  use-hosts: true\n")
	sb.WriteString("  respect-rules: true\n")
	if s.DNSMode == "fake-ip" {
		sb.WriteString("  fake-ip-range: 198.18.0.0/15\n")
		sb.WriteString("  fake-ip-filter:\n")
		sb.WriteString("    - '*.lan'\n")
		sb.WriteString("    - 'localhost.ptlogin2.qq.com'\n")
	}
	if s.EnableDNSoTLS {
		sb.WriteString(fmt.Sprintf("  default-nameserver:\n    - %s\n", s.DNSServer))
		sb.WriteString(fmt.Sprintf("  nameserver:\n    - tls://%s\n", s.DNSoTLSServer))
	} else {
		sb.WriteString(fmt.Sprintf("  nameserver:\n    - %s\n", s.DNSServer))
	}
	sb.WriteString(fmt.Sprintf("  fallback:\n    - %s\n", s.FallbackDNS))
	sb.WriteString("  fallback-filter:\n")
	sb.WriteString("    geoip: true\n")
	sb.WriteString("    geoip-code: IR\n")
	sb.WriteString("    ipcidr:\n")
	sb.WriteString("      - 240.0.0.0/4\n\n")

	// Sniffing
	if s.EnableSniffing {
		sb.WriteString("sniffer:\n")
		sb.WriteString("  enable: true\n")
		sb.WriteString("  sniff:\n")
		sb.WriteString("    HTTP:\n      ports: [80, 8080-8880]\n")
		sb.WriteString("    TLS:\n      ports: [443, 8443]\n")
		sb.WriteString("    QUIC:\n      ports: [443]\n\n")
	}

	// Proxy
	sb.WriteString("proxies:\n")

	fp := p.Fingerprint
	if fp == "" && s.EnableUTLS { fp = s.UTLSFingerprint }

	switch p.Protocol {
	case "vmess":
		sb.WriteString(fmt.Sprintf("  - name: \"%s\"\n    type: vmess\n", proxyName))
		sb.WriteString(fmt.Sprintf("    server: %s\n    port: %d\n", p.Address, p.Port))
		sb.WriteString(fmt.Sprintf("    uuid: %s\n    alterId: %d\n    cipher: %s\n",
			p.UUID, p.AlterId, p.Security))
		sb.WriteString("    udp: true\n")
		if p.TLS == "tls" {
			sb.WriteString("    tls: true\n    skip-cert-verify: false\n")
			if p.SNI != "" { sb.WriteString(fmt.Sprintf("    servername: %s\n", p.SNI)) }
			if s.EnableTLS13 { sb.WriteString("    # tls13 enabled via client\n") }
		}
		if p.Network != "tcp" && p.Network != "" {
			sb.WriteString(fmt.Sprintf("    network: %s\n", p.Network))
			writeNetworkOpts(&sb, p)
		}
		if fp != "" { sb.WriteString(fmt.Sprintf("    client-fingerprint: %s\n", fp)) }

	case "vless":
		sb.WriteString(fmt.Sprintf("  - name: \"%s\"\n    type: vless\n", proxyName))
		sb.WriteString(fmt.Sprintf("    server: %s\n    port: %d\n", p.Address, p.Port))
		sb.WriteString(fmt.Sprintf("    uuid: %s\n", p.UUID))
		sb.WriteString("    udp: true\n")
		if p.Flow != "" { sb.WriteString(fmt.Sprintf("    flow: %s\n", p.Flow)) }
		if p.TLS == "tls" || p.TLS == "reality" {
			sb.WriteString("    tls: true\n    skip-cert-verify: false\n")
			if p.SNI != "" { sb.WriteString(fmt.Sprintf("    servername: %s\n", p.SNI)) }
		}
		if p.TLS == "reality" {
			sb.WriteString("    reality-opts:\n")
			sb.WriteString(fmt.Sprintf("      public-key: %s\n", p.PublicKey))
			if p.ShortID != "" { sb.WriteString(fmt.Sprintf("      short-id: %s\n", p.ShortID)) }
		}
		if fp != "" { sb.WriteString(fmt.Sprintf("    client-fingerprint: %s\n", fp)) }
		if p.Network != "tcp" && p.Network != "" {
			sb.WriteString(fmt.Sprintf("    network: %s\n", p.Network))
			writeNetworkOpts(&sb, p)
		}

	case "trojan":
		sb.WriteString(fmt.Sprintf("  - name: \"%s\"\n    type: trojan\n", proxyName))
		sb.WriteString(fmt.Sprintf("    server: %s\n    port: %d\n", p.Address, p.Port))
		sb.WriteString(fmt.Sprintf("    password: %s\n", p.Password))
		sb.WriteString("    udp: true\n    tls: true\n    skip-cert-verify: false\n")
		if p.SNI != "" { sb.WriteString(fmt.Sprintf("    sni: %s\n", p.SNI)) }
		if fp != "" { sb.WriteString(fmt.Sprintf("    client-fingerprint: %s\n", fp)) }
		if p.Network != "tcp" && p.Network != "" {
			sb.WriteString(fmt.Sprintf("    network: %s\n", p.Network))
			writeNetworkOpts(&sb, p)
		}
	}

	// Mux
	if s.EnableMux {
		sb.WriteString(fmt.Sprintf("    smux:\n      enabled: true\n      protocol: %s\n      max-streams: %d\n",
			s.MuxProtocol, s.MuxConcurrency))
	}

	// Proxy Groups
	sb.WriteString("\nproxy-groups:\n")
	sb.WriteString(fmt.Sprintf("  - name: \"🚀 Proxy\"\n    type: select\n    proxies:\n      - \"%s\"\n      - DIRECT\n\n", proxyName))
	sb.WriteString(fmt.Sprintf("  - name: \"♻️ Auto\"\n    type: url-test\n    proxies:\n      - \"%s\"\n    url: http://www.gstatic.com/generate_204\n    interval: 300\n    tolerance: 50\n\n", proxyName))
	sb.WriteString("  - name: \"🎯 Direct\"\n    type: select\n    proxies:\n      - DIRECT\n      - 🚀 Proxy\n\n")

	// Rules
	sb.WriteString("rules:\n")

	// Block list first
	if s.EnableBlock {
		for _, d := range s.BlockList {
			d = strings.TrimSpace(d)
			if d != "" { sb.WriteString(fmt.Sprintf("  - DOMAIN-SUFFIX,%s,REJECT\n", d)) }
		}
	}

	// Bypass IPs (private + custom)
	// Always bypass private/local IPs
	privateIPs := []string{
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}
	for _, ip := range privateIPs {
		sb.WriteString(fmt.Sprintf("  - IP-CIDR,%s,DIRECT,no-resolve\n", ip))
	}

	// Custom bypass domains
	if s.EnableBypass {
		for _, d := range s.BypassDomains {
			d = strings.TrimSpace(d)
			if d == "" { continue }
			if strings.Contains(d, "/") {
				// It's a CIDR
				sb.WriteString(fmt.Sprintf("  - IP-CIDR,%s,DIRECT,no-resolve\n", d))
			} else if net.ParseIP(d) != nil {
				sb.WriteString(fmt.Sprintf("  - IP-CIDR,%s/32,DIRECT,no-resolve\n", d))
			} else {
				sb.WriteString(fmt.Sprintf("  - DOMAIN-SUFFIX,%s,DIRECT\n", d))
			}
		}
		// Custom bypass IPs
		for _, ip := range s.BypassIPs {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				// avoid duplicating private IPs
				alreadyAdded := false
				for _, p := range privateIPs {
					if p == ip { alreadyAdded = true; break }
				}
				if !alreadyAdded {
					sb.WriteString(fmt.Sprintf("  - IP-CIDR,%s,DIRECT,no-resolve\n", ip))
				}
			}
		}
	}

	// Iran GeoIP — always direct
	sb.WriteString("  - GEOIP,IR,DIRECT,no-resolve\n")
	sb.WriteString("  - GEOIP,private,DIRECT,no-resolve\n")

	// Final
	sb.WriteString("  - MATCH,🚀 Proxy\n")

	return sb.String()
}

func writeNetworkOpts(sb *strings.Builder, p V2RayConfig) {
	switch p.Network {
	case "ws":
		sb.WriteString(fmt.Sprintf("    ws-opts:\n      path: \"%s\"\n", p.Path))
		if p.Host != "" {
			sb.WriteString(fmt.Sprintf("      headers:\n        Host: \"%s\"\n", p.Host))
		}
	case "grpc":
		sb.WriteString(fmt.Sprintf("    grpc-opts:\n      grpc-service-name: \"%s\"\n", p.Path))
	case "h2":
		sb.WriteString(fmt.Sprintf("    h2-opts:\n      path: \"%s\"\n", p.Path))
		if p.Host != "" {
			sb.WriteString(fmt.Sprintf("      host:\n        - \"%s\"\n", p.Host))
		}
	case "httpupgrade":
		sb.WriteString(fmt.Sprintf("    httpupgrade-opts:\n      path: \"%s\"\n", p.Path))
		if p.Host != "" {
			sb.WriteString(fmt.Sprintf("      host: \"%s\"\n", p.Host))
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// SERVER INFO
// ─────────────────────────────────────────────────────────────────────────────

func getServerInfo(address string) map[string]string {
	info := map[string]string{
		"address": address,
		"status":  "unknown",
		"ping":    "N/A",
		"country": "N/A",
		"org":     "N/A",
	}

	// Ping via TCP
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address+":443", 3*time.Second)
	if err == nil {
		conn.Close()
		info["ping"] = fmt.Sprintf("%dms", time.Since(start).Milliseconds())
		info["status"] = "online"
	} else {
		// try port 80
		start = time.Now()
		conn, err = net.DialTimeout("tcp", address+":80", 3*time.Second)
		if err == nil {
			conn.Close()
			info["ping"] = fmt.Sprintf("%dms", time.Since(start).Milliseconds())
			info["status"] = "online"
		} else {
			info["status"] = "offline"
		}
	}

	// IP info
	client := &http.Client{Timeout: 4 * time.Second}
	resp, err := client.Get("http://ip-api.com/json/" + address + "?fields=country,countryCode,org,isp")
	if err == nil {
		defer resp.Body.Close()
		var result map[string]interface{}
		if json.NewDecoder(resp.Body).Decode(&result) == nil {
			if c, ok := result["country"].(string); ok  { info["country"] = c }
			if cc, ok := result["countryCode"].(string); ok { info["countryCode"] = cc }
			if o, ok := result["org"].(string); ok      { info["org"] = o }
		}
	}

	return info
}

// ─────────────────────────────────────────────────────────────────────────────
// HTTP HANDLERS
// ─────────────────────────────────────────────────────────────────────────────

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, indexHTML)
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	link := r.URL.Query().Get("link")
	cfg, err := parseV2RayLink(link)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(cfg)
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost { w.WriteHeader(405); return }

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(GenerateResponse{Error: "invalid JSON"})
		return
	}

	parsed, err := parseV2RayLink(req.V2RayLink)
	if err != nil {
		json.NewEncoder(w).Encode(GenerateResponse{Error: err.Error()})
		return
	}

	configText := generateConfig(parsed, req.Settings)
	resp := GenerateResponse{Success: true, Config: configText}

	// Encrypt if requested
	if req.Settings.EnablePassword && req.Settings.Password != "" {
		enc, err := encryptConfig(configText, req.Settings.Password)
		if err == nil {
			resp.Encrypted = enc
		}
	}

	// Save profile
	appData := loadData()
	profile := Profile{
		ID:          uuid.New().String(),
		Name:        req.ProfileName,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		V2RayLink:   req.V2RayLink,
		ParsedProxy: parsed,
		Settings:    req.Settings,
		Config:      configText,
	}
	if profile.Name == "" {
		profile.Name = fmt.Sprintf("Profile %d", len(appData.Profiles)+1)
	}
	appData.Profiles = append(appData.Profiles, profile)
	appData.LastSettings = req.Settings
	saveData(appData)
	resp.ProfileID = profile.ID

	json.NewEncoder(w).Encode(resp)
}

func handleProfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		appData := loadData()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"profiles":     appData.Profiles,
			"lastSettings": appData.LastSettings,
		})
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		appData := loadData()
		newProfiles := []Profile{}
		for _, p := range appData.Profiles {
			if p.ID != id { newProfiles = append(newProfiles, p) }
		}
		appData.Profiles = newProfiles
		saveData(appData)
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

func handleDecrypt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost { w.WriteHeader(405); return }

	var req struct {
		Content  string `json:"content"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	content := strings.TrimSpace(req.Content)

	// If not encrypted, just parse it
	if !strings.HasPrefix(content, "NAPSTER_ENC:") {
		// Try to extract info from plain YAML config
		info := parseYAMLConfig(content)
		json.NewEncoder(w).Encode(ImportResponse{
			Success: true,
			Decrypted: content,
			Info: info,
		})
		return
	}

	// Encrypted — need password
	if req.Password == "" {
		json.NewEncoder(w).Encode(ImportResponse{
			Error: "NEEDS_PASSWORD",
		})
		return
	}

	decrypted, err := decryptConfig(content, req.Password)
	if err != nil {
		json.NewEncoder(w).Encode(ImportResponse{Error: err.Error()})
		return
	}

	info := parseYAMLConfig(decrypted)
	json.NewEncoder(w).Encode(ImportResponse{
		Success:   true,
		Decrypted: decrypted,
		Info:      info,
	})
}

func parseYAMLConfig(content string) map[string]string {
	info := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# device-lock:") {
			info["deviceLock"] = strings.TrimSpace(strings.TrimPrefix(line, "# device-lock:"))
		}
		if strings.HasPrefix(line, "# user-agent:") {
			info["userAgent"] = strings.TrimSpace(strings.TrimPrefix(line, "# user-agent:"))
		}
		if strings.HasPrefix(line, "mode:") {
			info["mode"] = strings.TrimSpace(strings.TrimPrefix(line, "mode:"))
		}
		if strings.HasPrefix(line, "    server:") {
			info["server"] = strings.TrimSpace(strings.TrimPrefix(line, "    server:"))
		}
		if strings.HasPrefix(line, "    type:") && info["protocol"] == "" {
			info["protocol"] = strings.TrimSpace(strings.TrimPrefix(line, "    type:"))
		}
		if strings.HasPrefix(line, "    name:") && info["name"] == "" {
			info["name"] = strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, "    name:")), "\"")
		}
	}
	return info
}

func handleServerInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	addr := r.URL.Query().Get("address")
	if addr == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "no address"})
		return
	}
	info := getServerInfo(addr)
	json.NewEncoder(w).Encode(info)
}

func handleNewDeviceID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"deviceId": uuid.New().String()})
}

// ─────────────────────────────────────────────────────────────────────────────
// HELPERS
// ─────────────────────────────────────────────────────────────────────────────

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
func toIntStr(s string) int { var i int; fmt.Sscanf(s, "%d", &i); return i }

// ─────────────────────────────────────────────────────────────────────────────
// MAIN
// ─────────────────────────────────────────────────────────────────────────────

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/parse",       handleParse)
	mux.HandleFunc("/api/generate",    handleGenerate)
	mux.HandleFunc("/api/profiles",    handleProfiles)
	mux.HandleFunc("/api/decrypt",     handleDecrypt)
	mux.HandleFunc("/api/server-info", handleServerInfo)
	mux.HandleFunc("/api/device-id",   handleNewDeviceID)

	port := "8080"
	if p := os.Getenv("PORT"); p != "" { port = p }

	log.Printf("🚀 Napster Config Generator → http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// ─────────────────────────────────────────────────────────────────────────────
// EMBEDDED HTML
// ─────────────────────────────────────────────────────────────────────────────

const indexHTML = `<!DOCTYPE html>
<html lang="fa" dir="rtl">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Napster Config Generator</title>
<style>
@import url('https://fonts.googleapis.com/css2?family=Vazirmatn:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;600&display=swap');
:root {
  --bg:#0a0e1a;--surface:#111827;--surface2:#1a2236;--border:#1e2d4a;
  --accent:#3b82f6;--accent2:#6366f1;--green:#10b981;--red:#ef4444;
  --yellow:#f59e0b;--text:#e2e8f0;--text2:#94a3b8;
}
body.light {
  --bg:#f1f5f9;--surface:#ffffff;--surface2:#f8fafc;--border:#e2e8f0;
  --text:#0f172a;--text2:#64748b;
}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--text);font-family:'Vazirmatn',sans-serif;min-height:100vh;
  background-image:radial-gradient(ellipse at 20% 50%,rgba(59,130,246,.06) 0%,transparent 50%),
  radial-gradient(ellipse at 80% 20%,rgba(99,102,241,.06) 0%,transparent 50%)}
header{padding:16px 32px;border-bottom:1px solid var(--border);display:flex;align-items:center;gap:12px;
  background:rgba(17,24,39,.8);backdrop-filter:blur(12px);position:sticky;top:0;z-index:100}
.logo{width:38px;height:38px;background:linear-gradient(135deg,var(--accent),var(--accent2));
  border-radius:10px;display:flex;align-items:center;justify-content:center;font-size:18px}
header h1{font-size:18px;font-weight:700}
header p{font-size:12px;color:var(--text2);margin-top:1px}
.header-actions{margin-right:auto;display:flex;gap:8px;align-items:center}
.badge{background:rgba(59,130,246,.15);border:1px solid rgba(59,130,246,.3);color:var(--accent);
  padding:3px 10px;border-radius:20px;font-size:11px;font-weight:600}
.tabs{display:flex;gap:0;border-bottom:1px solid var(--border);background:var(--surface2);padding:0 32px}
.tab{padding:12px 20px;font-size:13px;font-weight:500;cursor:pointer;border-bottom:2px solid transparent;
  color:var(--text2);transition:.2s;font-family:inherit;background:none;border-top:none;border-left:none;border-right:none}
.tab.active{color:var(--accent);border-bottom-color:var(--accent)}
.tab:hover{color:var(--text)}
.page{display:none;padding:28px 32px;max-width:1600px;margin:0 auto}
.page.active{display:block}
.grid2{display:grid;grid-template-columns:1fr 1fr;gap:20px}
.grid3{display:grid;grid-template-columns:1fr 1fr 1fr;gap:14px}
@media(max-width:900px){.grid2{grid-template-columns:1fr}.grid3{grid-template-columns:1fr 1fr}}
.panel{background:var(--surface);border:1px solid var(--border);border-radius:14px;overflow:hidden;margin-bottom:18px}
.ph{padding:14px 20px;border-bottom:1px solid var(--border);display:flex;align-items:center;gap:8px;
  font-weight:600;font-size:14px;background:var(--surface2)}
.pb{padding:20px;display:flex;flex-direction:column;gap:14px}
.sec{font-size:11px;font-weight:600;color:var(--accent);text-transform:uppercase;letter-spacing:1px;
  margin-bottom:4px;display:flex;align-items:center;gap:8px}
.sec::after{content:'';flex:1;height:1px;background:var(--border)}
.field{display:flex;flex-direction:column;gap:5px}
.field label{font-size:12px;color:var(--text2);font-weight:500}
input[type=text],input[type=number],input[type=password],select,textarea{
  background:var(--bg);border:1px solid var(--border);color:var(--text);padding:9px 12px;
  border-radius:9px;font-family:inherit;font-size:13px;transition:border-color .2s,box-shadow .2s;
  width:100%;direction:ltr}
input:focus,select:focus,textarea:focus{outline:none;border-color:var(--accent);
  box-shadow:0 0 0 3px rgba(59,130,246,.15)}
textarea{resize:vertical;min-height:80px}
.ibtn{display:flex;gap:7px}
.ibtn input,.ibtn select{flex:1}
.trow{display:flex;align-items:center;justify-content:space-between;padding:10px 12px;
  background:var(--bg);border:1px solid var(--border);border-radius:9px}
.tlabel{font-size:13px;font-weight:500}
.tdesc{font-size:11px;color:var(--text2);margin-top:1px}
.toggle{position:relative;width:42px;height:23px;flex-shrink:0}
.toggle input{opacity:0;width:0;height:0}
.ts{position:absolute;inset:0;background:var(--border);border-radius:23px;cursor:pointer;transition:.3s}
.ts::before{content:'';position:absolute;width:17px;height:17px;left:3px;top:3px;
  background:white;border-radius:50%;transition:.3s}
.toggle input:checked+.ts{background:var(--accent)}
.toggle input:checked+.ts::before{transform:translateX(19px)}
.btn{padding:9px 18px;border-radius:9px;border:none;cursor:pointer;font-family:inherit;
  font-size:13px;font-weight:600;transition:all .2s;display:flex;align-items:center;gap:7px;justify-content:center}
.btn-primary{background:linear-gradient(135deg,var(--accent),var(--accent2));color:#fff;
  width:100%;padding:13px;font-size:14px;box-shadow:0 4px 15px rgba(59,130,246,.3)}
.btn-primary:hover{transform:translateY(-1px);box-shadow:0 6px 20px rgba(59,130,246,.4)}
.btn-sm{background:var(--surface2);border:1px solid var(--border);color:var(--text2);
  padding:7px 12px;font-size:12px;flex-shrink:0}
.btn-sm:hover{border-color:var(--accent);color:var(--accent)}
.btn-green{background:rgba(16,185,129,.15);border:1px solid rgba(16,185,129,.3);color:var(--green)}
.btn-green:hover{background:rgba(16,185,129,.25)}
.btn-purple{background:rgba(99,102,241,.15);border:1px solid rgba(99,102,241,.3);color:var(--accent2)}
.btn-purple:hover{background:rgba(99,102,241,.25)}
.btn-red{background:rgba(239,68,68,.1);border:1px solid rgba(239,68,68,.3);color:var(--red)}
.btn-red:hover{background:rgba(239,68,68,.2)}
.parsed-box{background:var(--bg);border:1px solid var(--border);border-radius:9px;padding:12px;
  font-family:'JetBrains Mono',monospace;font-size:12px;color:var(--text2);display:none}
.parsed-box.show{display:block}
.pbadge{display:inline-block;background:rgba(59,130,246,.2);color:var(--accent);
  border-radius:6px;padding:2px 8px;font-size:11px;margin-bottom:7px;font-weight:600}
.out-ta{width:100%;height:420px;background:var(--bg);border:1px solid var(--border);
  border-radius:11px;padding:14px;font-family:'JetBrains Mono',monospace;font-size:12px;
  color:#86efac;line-height:1.7;resize:none;direction:ltr}
.status{display:flex;align-items:center;gap:7px;padding:9px 12px;border-radius:9px;font-size:13px;font-weight:500}
.status.ok{background:rgba(16,185,129,.1);border:1px solid rgba(16,185,129,.3);color:var(--green)}
.status.err{background:rgba(239,68,68,.1);border:1px solid rgba(239,68,68,.3);color:var(--red)}
.status.warn{background:rgba(245,158,11,.1);border:1px solid rgba(245,158,11,.3);color:var(--yellow)}
.status.hide{display:none}
.dot{width:7px;height:7px;border-radius:50%;background:currentColor}
.prof-card{background:var(--bg);border:1px solid var(--border);border-radius:11px;padding:14px;
  display:flex;flex-direction:column;gap:8px}
.prof-name{font-weight:600;font-size:14px}
.prof-meta{font-size:11px;color:var(--text2);display:flex;gap:12px;flex-wrap:wrap}
.prof-actions{display:flex;gap:7px;margin-top:4px}
.ping-badge{display:inline-flex;align-items:center;gap:5px;font-size:12px;padding:3px 9px;
  border-radius:20px;font-weight:600}
.ping-ok{background:rgba(16,185,129,.15);color:var(--green)}
.ping-err{background:rgba(239,68,68,.15);color:var(--red)}
.ping-wait{background:rgba(245,158,11,.15);color:var(--yellow)}
.import-area{width:100%;min-height:120px;background:var(--bg);border:1px solid var(--border);
  border-radius:9px;padding:12px;font-family:'JetBrains Mono',monospace;font-size:12px;
  color:var(--text);resize:vertical;direction:ltr}
.qr-canvas{border-radius:10px;border:3px solid var(--border)}
.server-info{display:grid;grid-template-columns:1fr 1fr;gap:8px}
.si-item{background:var(--bg);border:1px solid var(--border);border-radius:8px;
  padding:10px 12px;font-size:12px}
.si-label{color:var(--text2);margin-bottom:3px;font-size:11px}
.si-val{font-weight:600;font-family:'JetBrains Mono',monospace}
::-webkit-scrollbar{width:5px;height:5px}
::-webkit-scrollbar-track{background:var(--bg)}
::-webkit-scrollbar-thumb{background:var(--border);border-radius:3px}
</style>
</head>
<body>
<header>
  <div class="logo">⚡</div>
  <div>
    <h1>Napster Config Generator</h1>
    <p>ساخت کانفیگ حرفه‌ای V2Ray</p>
  </div>
  <div class="header-actions">
    <div class="badge">v3.0</div>
    <button class="btn btn-sm" onclick="toggleTheme()" id="themeBtn">🌙 تاریک</button>
  </div>
</header>

<div class="tabs">
  <button class="tab active" onclick="showTab('generate')">⚡ ساخت کانفیگ</button>
  <button class="tab" onclick="showTab('profiles')">📋 پروفایل‌ها</button>
  <button class="tab" onclick="showTab('import')">📂 Import کانفیگ</button>
</div>

<!-- ─── TAB: GENERATE ─── -->
<div class="page active" id="tab-generate">
<div class="grid2">
<div>

  <!-- Link -->
  <div class="panel">
    <div class="ph">🔗 لینک V2Ray</div>
    <div class="pb">
      <div class="field">
        <label>لینک vmess / vless / trojan</label>
        <div class="ibtn">
          <input type="text" id="v2rayLink" placeholder="vmess://... یا vless://... یا trojan://...">
          <button class="btn btn-sm" onclick="parseLink()">📡 پارس</button>
        </div>
      </div>
      <div class="parsed-box" id="parsedBox"></div>
      <div id="serverInfoBox" style="display:none">
        <div class="sec">اطلاعات سرور</div>
        <div class="server-info" id="serverInfoGrid"></div>
      </div>
      <div class="field">
        <label>نام پروفایل</label>
        <input type="text" id="profileName" placeholder="مثال: سرور آلمان">
      </div>
    </div>
  </div>

  <!-- Device ID -->
  <div class="panel">
    <div class="ph">📱 Device ID</div>
    <div class="pb">
      <div class="trow">
        <div><div class="tlabel">فعال کردن Device Lock</div>
        <div class="tdesc">کانفیگ فقط روی این Device ID کار کند</div></div>
        <label class="toggle"><input type="checkbox" id="enableDeviceLock" onchange="toggleField('deviceLockFields',this.checked)">
        <span class="ts"></span></label>
      </div>
      <div id="deviceLockFields" style="display:none">
        <div class="field">
          <label>Device ID (UUID دستگاه — از اپلیکیشن نپستر کپی کن)</label>
          <div class="ibtn">
            <input type="text" id="deviceId" placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" dir="ltr">
            <button class="btn btn-sm" onclick="genDeviceId()">🔄 جدید</button>
          </div>
        </div>
        <div class="field">
          <label>User Agent (اختیاری)</label>
          <input type="text" id="userAgent" value="Napster/2.0" dir="ltr">
        </div>
      </div>
      <div id="deviceLockOff" style="font-size:12px;color:var(--text2)">
        ℹ️ بدون Device Lock — کانفیگ روی هر دستگاهی کار می‌کند
      </div>
    </div>
  </div>

  <!-- Security -->
  <div class="panel">
    <div class="ph">🔒 تنظیمات امنیتی</div>
    <div class="pb">
      <div class="sec">رمزگذاری خروجی</div>
      <div class="trow">
        <div><div class="tlabel">رمزگذاری فایل خروجی</div>
        <div class="tdesc">فایل با AES-256 رمز می‌شود</div></div>
        <label class="toggle"><input type="checkbox" id="enablePassword" onchange="toggleField('passwordField',this.checked)">
        <span class="ts"></span></label>
      </div>
      <div id="passwordField" style="display:none">
        <div class="field">
          <label>رمز عبور</label>
          <input type="password" id="configPassword" placeholder="رمز قوی وارد کنید" dir="ltr">
        </div>
      </div>

      <div class="sec">TLS & Fingerprint</div>
      <div class="trow">
        <div><div class="tlabel">TLS 1.3</div></div>
        <label class="toggle"><input type="checkbox" id="enableTls13" checked><span class="ts"></span></label>
      </div>
      <div class="trow">
        <div><div class="tlabel">uTLS Browser Fingerprint</div>
        <div class="tdesc">جعل fingerprint مرورگر برای دور زدن DPI</div></div>
        <label class="toggle"><input type="checkbox" id="enableUtls" checked onchange="toggleField('utlsField',this.checked)">
        <span class="ts"></span></label>
      </div>
      <div id="utlsField" class="field">
        <label>نوع Fingerprint</label>
        <select id="utlsFingerprint">
          <option value="chrome">Chrome (توصیه)</option>
          <option value="firefox">Firefox</option>
          <option value="safari">Safari</option>
          <option value="ios">iOS Safari</option>
          <option value="android">Android Chrome</option>
          <option value="edge">Microsoft Edge</option>
          <option value="random">Random</option>
        </select>
      </div>

      <div class="sec">Mux / Multiplexing</div>
      <div class="trow">
        <div><div class="tlabel">فعال‌سازی Mux</div>
        <div class="tdesc">ادغام چند اتصال در یک تانل</div></div>
        <label class="toggle"><input type="checkbox" id="enableMux" onchange="toggleField('muxFields',this.checked)">
        <span class="ts"></span></label>
      </div>
      <div id="muxFields" style="display:none" class="grid2">
        <div class="field"><label>Mux Protocol</label>
          <select id="muxProtocol">
            <option value="smux">smux</option>
            <option value="yamux">yamux</option>
            <option value="h2mux">h2mux</option>
          </select>
        </div>
        <div class="field"><label>Concurrency</label>
          <input type="number" id="muxConcurrency" value="8" min="1" max="64">
        </div>
      </div>
    </div>
  </div>

  <!-- DNS -->
  <div class="panel">
    <div class="ph">🌐 تنظیمات DNS</div>
    <div class="pb">
      <div class="grid2">
        <div class="field"><label>DNS Mode</label>
          <select id="dnsMode">
            <option value="fake-ip">Fake-IP (توصیه)</option>
            <option value="redir-host">Redir-Host</option>
          </select>
        </div>
        <div class="field"><label>DNS اصلی</label>
          <select id="dnsServer">
            <option value="1.1.1.1">Cloudflare 1.1.1.1</option>
            <option value="8.8.8.8">Google 8.8.8.8</option>
            <option value="9.9.9.9">Quad9</option>
          </select>
        </div>
      </div>
      <div class="field"><label>DNS Fallback</label>
        <select id="fallbackDns">
          <option value="8.8.8.8">Google 8.8.8.8</option>
          <option value="1.1.1.1">Cloudflare 1.1.1.1</option>
          <option value="208.67.220.220">OpenDNS</option>
        </select>
      </div>
      <div class="trow">
        <div><div class="tlabel">DNS over TLS</div></div>
        <label class="toggle"><input type="checkbox" id="enableDnsoTls" onchange="toggleField('dotField',this.checked)">
        <span class="ts"></span></label>
      </div>
      <div id="dotField" style="display:none">
        <div class="field"><label>DoT Server</label>
          <select id="dnsotlsServer">
            <option value="cloudflare-dns.com">cloudflare-dns.com</option>
            <option value="dns.google">dns.google</option>
            <option value="dns.quad9.net">dns.quad9.net</option>
          </select>
        </div>
      </div>
    </div>
  </div>

  <!-- Routing -->
  <div class="panel">
    <div class="ph">🗺️ مسیریابی</div>
    <div class="pb">
      <div class="grid2">
        <div class="field"><label>Proxy Mode</label>
          <select id="proxyMode">
            <option value="rule">Rule (هوشمند)</option>
            <option value="global">Global</option>
            <option value="direct">Direct</option>
          </select>
        </div>
        <div class="field"><label>Log Level</label>
          <select id="logLevel">
            <option value="warning">Warning</option>
            <option value="info">Info</option>
            <option value="debug">Debug</option>
            <option value="silent">Silent</option>
          </select>
        </div>
      </div>

      <div class="trow">
        <div><div class="tlabel">Bypass ایران</div>
        <div class="tdesc">سایت‌های ایرانی مستقیم — بدون پروکسی</div></div>
        <label class="toggle"><input type="checkbox" id="enableBypass" checked onchange="toggleField('bypassFields',this.checked)">
        <span class="ts"></span></label>
      </div>

      <div id="bypassFields">
        <div class="field"><label>دامنه‌های Bypass (هر خط یک دامنه)</label>
          <textarea id="bypassDomains" rows="5">ir
shaparak.ir
digikala.com
aparat.com
snapp.ir
divar.ir</textarea>
        </div>
        <div class="field"><label>IP های Bypass — CIDR (هر خط یک آیپی/رنج)</label>
          <textarea id="bypassIPs" rows="3" dir="ltr">192.168.0.0/16
10.0.0.0/8
172.16.0.0/12</textarea>
        </div>
      </div>

      <div class="trow">
        <div><div class="tlabel">Block دامنه‌ها</div></div>
        <label class="toggle"><input type="checkbox" id="enableBlock" onchange="toggleField('blockField',this.checked)">
        <span class="ts"></span></label>
      </div>
      <div id="blockField" style="display:none">
        <div class="field"><label>لیست Block</label>
          <textarea id="blockList" rows="3" placeholder="ads.example.com"></textarea>
        </div>
      </div>
    </div>
  </div>

  <!-- Advanced -->
  <div class="panel">
    <div class="ph">⚙️ پیشرفته</div>
    <div class="pb">
      <div class="grid3">
        <div class="field"><label>MTU</label><input type="number" id="mtu" value="1500" min="576" max="9000"></div>
        <div class="field"><label>فرمت خروجی</label>
          <select id="outputFormat">
            <option value="napster">Napster</option>
            <option value="clash">Clash</option>
          </select>
        </div>
        <div class="field"><label>IPv6</label>
          <div class="trow" style="margin-top:4px">
            <span style="font-size:12px">فعال</span>
            <label class="toggle"><input type="checkbox" id="enableIpv6"><span class="ts"></span></label>
          </div>
        </div>
      </div>
      <div class="trow">
        <div><div class="tlabel">Sniffing</div></div>
        <label class="toggle"><input type="checkbox" id="enableSniffing" checked><span class="ts"></span></label>
      </div>
    </div>
  </div>

  <button class="btn btn-primary" onclick="generate()">⚡ ساخت کانفیگ</button>
</div>

<!-- Output -->
<div>
  <div class="panel" style="position:sticky;top:80px">
    <div class="ph">📄 خروجی</div>
    <div class="pb">
      <div id="statusBar" class="status hide"></div>
      <textarea class="out-ta" id="output" readonly placeholder="# کانفیگ اینجا نمایش داده می‌شود..."></textarea>
      <div style="display:flex;gap:8px;flex-wrap:wrap;margin-top:4px">
        <button class="btn btn-green" onclick="copyOut()">📋 کپی</button>
        <button class="btn btn-purple" onclick="downloadOut()">💾 دانلود .yaml</button>
        <button class="btn btn-sm" onclick="downloadEnc()">🔐 دانلود رمزشده</button>
        <button class="btn btn-sm" onclick="showQR()">📤 QR Code</button>
        <button class="btn btn-red" onclick="clearAll()">🗑️</button>
      </div>
      <canvas id="qrCanvas" class="qr-canvas" style="display:none;margin-top:12px;max-width:220px"></canvas>
    </div>
  </div>
</div>
</div>
</div>

<!-- ─── TAB: PROFILES ─── -->
<div class="page" id="tab-profiles">
  <div style="max-width:900px">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px">
      <h2 style="font-size:18px;font-weight:700">📋 پروفایل‌های ذخیره‌شده</h2>
      <button class="btn btn-sm" onclick="loadProfiles()">🔄 بارگذاری</button>
    </div>
    <div id="profilesContainer" style="display:flex;flex-direction:column;gap:12px">
      <div style="color:var(--text2);font-size:13px;text-align:center;padding:40px">در حال بارگذاری...</div>
    </div>
  </div>
</div>

<!-- ─── TAB: IMPORT ─── -->
<div class="page" id="tab-import">
  <div style="max-width:900px">
    <div class="panel">
      <div class="ph">📂 Import و رمزگشایی کانفیگ نپستر</div>
      <div class="pb">
        <div class="field">
          <label>محتوای فایل کانفیگ را اینجا paste کنید</label>
          <textarea class="import-area" id="importContent" placeholder="محتوای فایل .yaml یا کانفیگ رمزشده (NAPSTER_ENC:...) را paste کنید"></textarea>
        </div>
        <div id="importPwdWrap" style="display:none">
          <div class="field">
            <label>🔐 این کانفیگ رمز دارد — رمز را وارد کنید</label>
            <input type="password" id="importPassword" placeholder="رمز عبور" dir="ltr">
          </div>
        </div>
        <div style="display:flex;gap:8px">
          <button class="btn btn-primary" style="flex:1" onclick="importConfig()">🔍 آنالیز و رمزگشایی</button>
        </div>
        <div id="importStatus" class="status hide"></div>
        <div id="importResult" style="display:none">
          <div class="sec">اطلاعات استخراج‌شده</div>
          <div id="importInfo" style="background:var(--bg);border:1px solid var(--border);
            border-radius:9px;padding:14px;font-size:13px;margin-bottom:12px"></div>
          <div class="sec">محتوای کانفیگ</div>
          <textarea class="out-ta" id="importOutput" readonly style="height:350px"></textarea>
          <div style="display:flex;gap:8px;margin-top:8px">
            <button class="btn btn-green" onclick="copyImport()">📋 کپی</button>
            <button class="btn btn-purple" onclick="downloadImport()">💾 دانلود</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<script>
// ── Theme ──
let isDark = true;
function toggleTheme(){
  isDark=!isDark;
  document.body.classList.toggle('light',!isDark);
  document.getElementById('themeBtn').textContent=isDark?'🌙 تاریک':'☀️ روشن';
}

// ── Tabs ──
function showTab(name){
  document.querySelectorAll('.tab').forEach((t,i)=>t.classList.toggle('active',['generate','profiles','import'][i]===name));
  document.querySelectorAll('.page').forEach(p=>p.classList.remove('active'));
  document.getElementById('tab-'+name).classList.add('active');
  if(name==='profiles') loadProfiles();
}

// ── Toggle helper ──
function toggleField(id, show){
  const el=document.getElementById(id);
  if(el) el.style.display=show?'':'none';
  if(id==='deviceLockFields'){
    document.getElementById('deviceLockOff').style.display=show?'none':'';
  }
}

// ── Parse Link ──
async function parseLink(){
  const link=document.getElementById('v2rayLink').value.trim();
  if(!link) return;
  try{
    const r=await fetch('/api/parse?link='+encodeURIComponent(link));
    const d=await r.json();
    if(d.error){showParsed(null,d.error);return;}
    showParsed(d);
    // Auto-fill profile name
    if(d.remarks && !document.getElementById('profileName').value)
      document.getElementById('profileName').value=d.remarks;
    // Load server info
    fetchServerInfo(d.address);
  }catch(e){showParsed(null,e.message);}
}

function showParsed(d,err){
  const el=document.getElementById('parsedBox');
  if(err){el.innerHTML='<span style="color:var(--red)">❌ '+err+'</span>';el.classList.add('show');return;}
  const proto=d.protocol.toUpperCase();
  const id=d.uuid||d.password||'-';
  let html=`<span class="pbadge">${proto}</span> `;
  if(d.remarks) html+=`<b>${d.remarks}</b><br>`;
  html+=`🖥️ <b>${d.address}:${d.port}</b> &nbsp; 🌐 ${d.network||'tcp'} &nbsp; 🔒 ${d.tls||'none'}<br>`;
  html+=`🔑 ${id.substring(0,28)}${id.length>28?'...':''}`;
  if(d.flow) html+=`<br>⚡ flow: ${d.flow}`;
  if(d.publicKey) html+=`<br>🔑 Reality PubKey: ${d.publicKey.substring(0,20)}...`;
  if(d.sni) html+=`<br>🌍 SNI: ${d.sni}`;
  if(d.fingerprint) html+=`<br>🖱️ fp: ${d.fingerprint}`;
  el.innerHTML=html; el.classList.add('show');
}

async function fetchServerInfo(addr){
  if(!addr) return;
  const box=document.getElementById('serverInfoBox');
  const grid=document.getElementById('serverInfoGrid');
  box.style.display='block';
  grid.innerHTML='<div style="color:var(--text2);font-size:12px;padding:8px">در حال بررسی سرور...</div>';
  try{
    const r=await fetch('/api/server-info?address='+encodeURIComponent(addr));
    const d=await r.json();
    const pingClass=d.status==='online'?'ping-ok':'ping-err';
    grid.innerHTML=`
      <div class="si-item"><div class="si-label">وضعیت</div>
        <div class="si-val"><span class="ping-badge ${pingClass}">${d.status==='online'?'🟢 آنلاین':'🔴 آفلاین'}</span></div></div>
      <div class="si-item"><div class="si-label">پینگ</div><div class="si-val">${d.ping}</div></div>
      <div class="si-item"><div class="si-label">کشور</div><div class="si-val">${d.countryCode||''} ${d.country||'N/A'}</div></div>
      <div class="si-item"><div class="si-label">ISP/Org</div><div class="si-val" style="font-size:11px">${d.org||'N/A'}</div></div>
    `;
  }catch(e){grid.innerHTML='<div style="color:var(--text2);font-size:12px">خطا در دریافت اطلاعات سرور</div>';}
}

// ── Generate ──
let lastConfig='', lastEncrypted='';

async function generate(){
  const link=document.getElementById('v2rayLink').value.trim();
  if(!link){showStatus('err','لینک V2Ray را وارد کنید');return;}

  const s={
    enableDeviceLock: gc('enableDeviceLock'),
    deviceId:         gv('deviceId'),
    userAgent:        gv('userAgent'),
    enablePassword:   gc('enablePassword'),
    password:         gv('configPassword'),
    enableTls13:      gc('enableTls13'),
    enableUtls:       gc('enableUtls'),
    utlsFingerprint:  gv('utlsFingerprint'),
    enableMux:        gc('enableMux'),
    muxConcurrency:   parseInt(gv('muxConcurrency'))||8,
    muxProtocol:      gv('muxProtocol'),
    dnsMode:          gv('dnsMode'),
    dnsServer:        gv('dnsServer'),
    fallbackDNS:      gv('fallbackDns'),
    enableDnsoTls:    gc('enableDnsoTls'),
    dnsotlsServer:    gv('dnsotlsServer'),
    proxyMode:        gv('proxyMode'),
    logLevel:         gv('logLevel'),
    enableBypass:     gc('enableBypass'),
    bypassDomains:    gv('bypassDomains').split('\n').map(s=>s.trim()).filter(Boolean),
    bypassIPs:        gv('bypassIPs').split('\n').map(s=>s.trim()).filter(Boolean),
    enableBlock:      gc('enableBlock'),
    blockList:        gv('blockList').split('\n').map(s=>s.trim()).filter(Boolean),
    mtu:              parseInt(gv('mtu'))||1500,
    enableIpv6:       gc('enableIpv6'),
    enableSniffing:   gc('enableSniffing'),
    outputFormat:     gv('outputFormat'),
  };

  try{
    const r=await fetch('/api/generate',{method:'POST',
      headers:{'Content-Type':'application/json'},
      body:JSON.stringify({v2rayLink:link,profileName:gv('profileName'),settings:s})});
    const d=await r.json();
    if(!d.success){showStatus('err',d.error);return;}
    lastConfig=d.config;
    lastEncrypted=d.encrypted||'';
    document.getElementById('output').value=d.config;
    showStatus('ok','✅ کانفیگ با موفقیت ساخته شد — ذخیره شد در config.json');
  }catch(e){showStatus('err','خطا: '+e.message);}
}

function showStatus(type,msg){
  const el=document.getElementById('statusBar');
  el.className='status '+type;
  el.innerHTML='<div class="dot"></div>'+msg;
}

// ── Output actions ──
function copyOut(){
  if(!lastConfig) return;
  navigator.clipboard.writeText(lastConfig).then(()=>showStatus('ok','کپی شد ✓'));
}
function downloadOut(){
  if(!lastConfig) return;
  dl(lastConfig, gv('profileName')||'config', '.yaml');
}
function downloadEnc(){
  if(!lastEncrypted){showStatus('warn','ابتدا رمزگذاری را فعال کنید');return;}
  dl(lastEncrypted, (gv('profileName')||'config')+'-encrypted', '.txt');
}
function dl(content, name, ext){
  const a=document.createElement('a');
  a.href=URL.createObjectURL(new Blob([content],{type:'text/plain'}));
  a.download=name+ext; a.click();
}
function clearAll(){
  lastConfig='';lastEncrypted='';
  document.getElementById('output').value='';
  document.getElementById('statusBar').className='status hide';
  document.getElementById('parsedBox').classList.remove('show');
  document.getElementById('v2rayLink').value='';
  document.getElementById('serverInfoBox').style.display='none';
  document.getElementById('qrCanvas').style.display='none';
}

// ── QR Code (simple text QR via API) ──
async function showQR(){
  if(!lastConfig){showStatus('warn','ابتدا کانفیگ بسازید');return;}
  const canvas=document.getElementById('qrCanvas');
  // Use a simple QR library via CDN
  if(!window.QRCode){
    const s=document.createElement('script');
    s.src='https://cdnjs.cloudflare.com/ajax/libs/qrcodejs/1.0.0/qrcode.min.js';
    s.onload=()=>renderQR();
    document.head.appendChild(s);
  } else { renderQR(); }
}
function renderQR(){
  const canvas=document.getElementById('qrCanvas');
  canvas.style.display='block';
  canvas.width=220; canvas.height=220;
  try{
    new QRCode(canvas,{text:lastConfig.substring(0,500),width:220,height:220,
      colorDark:'#3b82f6',colorLight:'#0a0e1a'});
  }catch(e){}
}

// ── Device ID ──
async function genDeviceId(){
  try{const r=await fetch('/api/device-id');const d=await r.json();
    document.getElementById('deviceId').value=d.deviceId;}catch(e){}
}

// ── Profiles ──
async function loadProfiles(){
  try{
    const r=await fetch('/api/profiles');
    const d=await r.json();
    const c=document.getElementById('profilesContainer');
    if(!d.profiles||d.profiles.length===0){
      c.innerHTML='<div style="color:var(--text2);font-size:13px;text-align:center;padding:40px">هنوز پروفایلی ذخیره نشده</div>';
      return;
    }
    c.innerHTML=d.profiles.slice().reverse().map(p=>`
      <div class="prof-card">
        <div class="prof-name">${p.name||'بدون نام'}</div>
        <div class="prof-meta">
          <span>📅 ${p.createdAt}</span>
          <span class="pbadge" style="font-size:10px">${(p.parsedProxy.protocol||'?').toUpperCase()}</span>
          <span>🖥️ ${p.parsedProxy.address||'?'}:${p.parsedProxy.port||'?'}</span>
          ${p.settings.enableDeviceLock?'<span>🔒 Device Locked</span>':''}
          ${p.settings.enablePassword?'<span>🔐 رمزگذاری شده</span>':''}
        </div>
        <div class="prof-actions">
          <button class="btn btn-sm" onclick="viewProfile('${p.id}')">👁️ مشاهده</button>
          <button class="btn btn-sm" onclick="loadProfile('${p.id}')">✏️ بارگذاری</button>
          <button class="btn btn-red" style="font-size:12px;padding:7px 12px" onclick="deleteProfile('${p.id}')">🗑️</button>
        </div>
      </div>`).join('');
  }catch(e){}
}

let _profiles=[];
async function getProfiles(){
  if(_profiles.length) return _profiles;
  const r=await fetch('/api/profiles');
  const d=await r.json();
  _profiles=d.profiles||[];
  return _profiles;
}

async function viewProfile(id){
  const profiles=await getProfiles();
  const p=profiles.find(x=>x.id===id);
  if(!p) return;
  showTab('generate');
  document.getElementById('output').value=p.config;
  lastConfig=p.config;
  showStatus('ok','پروفایل: '+p.name);
}

async function loadProfile(id){
  const profiles=await getProfiles();
  const p=profiles.find(x=>x.id===id);
  if(!p) return;
  showTab('generate');
  document.getElementById('v2rayLink').value=p.v2rayLink;
  document.getElementById('profileName').value=p.name;
  const s=p.settings;
  setChk('enableDeviceLock',s.enableDeviceLock); toggleField('deviceLockFields',s.enableDeviceLock);
  sv('deviceId',s.deviceId||''); sv('userAgent',s.userAgent||'');
  setChk('enablePassword',s.enablePassword); toggleField('passwordField',s.enablePassword);
  setChk('enableTls13',s.enableTls13); setChk('enableUtls',s.enableUtls);
  toggleField('utlsField',s.enableUtls);
  sv('utlsFingerprint',s.utlsFingerprint||'chrome');
  setChk('enableMux',s.enableMux); toggleField('muxFields',s.enableMux);
  sv('muxProtocol',s.muxProtocol||'smux'); sv('muxConcurrency',s.muxConcurrency||8);
  sv('dnsMode',s.dnsMode||'fake-ip'); sv('dnsServer',s.dnsServer||'1.1.1.1');
  sv('fallbackDns',s.fallbackDNS||'8.8.8.8');
  setChk('enableDnsoTls',s.enableDnsoTls); toggleField('dotField',s.enableDnsoTls);
  sv('proxyMode',s.proxyMode||'rule'); sv('logLevel',s.logLevel||'warning');
  setChk('enableBypass',s.enableBypass); toggleField('bypassFields',s.enableBypass);
  sv('bypassDomains',(s.bypassDomains||[]).join('\n'));
  sv('bypassIPs',(s.bypassIPs||[]).join('\n'));
  setChk('enableBlock',s.enableBlock); toggleField('blockField',s.enableBlock);
  sv('blockList',(s.blockList||[]).join('\n'));
  sv('mtu',s.mtu||1500); setChk('enableIpv6',s.enableIpv6);
  setChk('enableSniffing',s.enableSniffing!==false); sv('outputFormat',s.outputFormat||'napster');
  parseLink();
  showStatus('ok','پروفایل «'+p.name+'» بارگذاری شد');
}

async function deleteProfile(id){
  if(!confirm('حذف شود؟')) return;
  await fetch('/api/profiles?id='+id,{method:'DELETE'});
  _profiles=[];
  loadProfiles();
}

// ── Import ──
async function importConfig(){
  const content=document.getElementById('importContent').value.trim();
  const pwd=document.getElementById('importPassword').value;
  if(!content){return;}
  try{
    const r=await fetch('/api/decrypt',{method:'POST',
      headers:{'Content-Type':'application/json'},
      body:JSON.stringify({content,password:pwd})});
    const d=await r.json();
    const wrap=document.getElementById('importPwdWrap');
    const status=document.getElementById('importStatus');
    const result=document.getElementById('importResult');

    if(d.error==='NEEDS_PASSWORD'){
      wrap.style.display='';
      status.className='status warn';
      status.innerHTML='<div class="dot"></div>🔐 این کانفیگ رمزگذاری شده — رمز را وارد کنید';
      result.style.display='none';
      return;
    }
    if(!d.success){
      status.className='status err';
      status.innerHTML='<div class="dot"></div>❌ '+d.error;
      result.style.display='none';
      return;
    }

    status.className='status ok';
    status.innerHTML='<div class="dot"></div>✅ کانفیگ خوانده شد';
    document.getElementById('importOutput').value=d.decrypted;
    const info=d.info||{};
    let infoHtml='';
    if(info.name) infoHtml+=`<div>📝 نام پروکسی: <b>${info.name}</b></div>`;
    if(info.protocol) infoHtml+=`<div>🔌 پروتکل: <b>${info.protocol}</b></div>`;
    if(info.server) infoHtml+=`<div>🖥️ سرور: <b>${info.server}</b></div>`;
    if(info.mode) infoHtml+=`<div>🗺️ Mode: <b>${info.mode}</b></div>`;
    if(info.deviceLock) infoHtml+=`<div>🔒 Device Lock: <b>${info.deviceLock}</b></div>`;
    if(info.userAgent) infoHtml+=`<div>📱 User-Agent: <b>${info.userAgent}</b></div>`;
    if(!infoHtml) infoHtml='<span style="color:var(--text2)">اطلاعات اضافی‌ای استخراج نشد</span>';
    document.getElementById('importInfo').innerHTML=infoHtml;
    result.style.display='block';
  }catch(e){
    document.getElementById('importStatus').className='status err';
    document.getElementById('importStatus').innerHTML='<div class="dot"></div>خطا: '+e.message;
  }
}

function copyImport(){
  navigator.clipboard.writeText(document.getElementById('importOutput').value);
}
function downloadImport(){
  dl(document.getElementById('importOutput').value,'imported-config','.yaml');
}

// ── Helpers ──
function gv(id){return document.getElementById(id)?.value||'';}
function gc(id){return document.getElementById(id)?.checked||false;}
function sv(id,v){const el=document.getElementById(id);if(el)el.value=v;}
function setChk(id,v){const el=document.getElementById(id);if(el)el.checked=!!v;}

// Init
window.onload=()=>{
  toggleField('deviceLockFields',false);
};
</script>
</body>
</html>
`
