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
)

// ─── Structs ──────────────────────────────────────────────────────────────────

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

// NapsterSettings — فقط تنظیمات ضروری
type NapsterSettings struct {
	EnableDeviceLock bool     `json:"enableDeviceLock"`
	DeviceID         string   `json:"deviceId"`
	UserAgent        string   `json:"userAgent"`
	EnablePassword   bool     `json:"enablePassword"`
	Password         string   `json:"password"`
	EnableBypass     bool     `json:"enableBypass"`
	BypassDomains    []string `json:"bypassDomains"`
	BypassIPs        []string `json:"bypassIPs"`
	ProxyMode        string   `json:"proxyMode"`
	LogLevel         string   `json:"logLevel"`
	EnableIPv6       bool     `json:"enableIpv6"`
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

type AppData struct {
	Profiles     []Profile       `json:"profiles"`
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
	NpvtB64   string `json:"npvtB64,omitempty"`
	ProfileID string `json:"profileId,omitempty"`
	Error     string `json:"error,omitempty"`
}

type NpvtFile struct {
	Version     int             `json:"version"`
	CreatedAt   string          `json:"createdAt"`
	DeviceLock  string          `json:"deviceLock,omitempty"`
	UserAgent   string          `json:"userAgent,omitempty"`
	V2RayLink   string          `json:"v2rayLink"`
	ParsedProxy V2RayConfig     `json:"proxy"`
	Settings    NapsterSettings `json:"settings"`
	Config      string          `json:"config"`
	Encrypted   bool            `json:"encrypted"`
}

// ─── Persistence ──────────────────────────────────────────────────────────────

func dataPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "napster-data.json"
	}
	return filepath.Join(filepath.Dir(exe), "napster-data.json")
}

func loadData() AppData {
	b, err := os.ReadFile(dataPath())
	if err != nil {
		return AppData{Profiles: []Profile{}, LastSettings: defaultSettings()}
	}
	var d AppData
	if err := json.Unmarshal(b, &d); err != nil {
		return AppData{Profiles: []Profile{}, LastSettings: defaultSettings()}
	}
	if d.Profiles == nil {
		d.Profiles = []Profile{}
	}
	return d
}

func saveData(d AppData) {
	b, _ := json.MarshalIndent(d, "", "  ")
	os.WriteFile(dataPath(), b, 0644)
}

func defaultSettings() NapsterSettings {
	return NapsterSettings{
		UserAgent:     "Napster/2.0",
		EnableBypass:  true,
		BypassDomains: []string{"ir", "shaparak.ir", "digikala.com", "aparat.com", "snapp.ir", "divar.ir"},
		BypassIPs:     []string{"192.168.0.0/16", "10.0.0.0/8", "172.16.0.0/12"},
		ProxyMode:     "rule",
		LogLevel:      "warning",
	}
}

// ─── Key Derivation ───────────────────────────────────────────────────────────

func deriveKey(password string, salt []byte) []byte {
	key := make([]byte, 32)
	h := sha256.New()
	for i := 0; i < 10000; i++ {
		h.Write([]byte(password))
		h.Write(salt)
		h.Write(key)
		copy(key, h.Sum(nil))
		h.Reset()
	}
	return key
}

// ─── AES-GCM ─────────────────────────────────────────────────────────────────

// gcmDecrypt: ورودی = nonce(12) + ciphertext + tag(16)
func gcmDecrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ns := gcm.NonceSize() // 12
	if len(data) < ns+gcm.Overhead() {
		return nil, fmt.Errorf("داده ناقص")
	}
	plain, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return nil, fmt.Errorf("رمز اشتباه یا داده خراب")
	}
	return plain, nil
}

// encrypt برای فرمت داخلی سرور
func encrypt(plain, password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := deriveKey(password, salt)
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
	ct := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return "NAPSTER_ENC:" + base64.StdEncoding.EncodeToString(append(salt, ct...)), nil
}

func decrypt(enc, password string) (string, error) {
	if !strings.HasPrefix(enc, "NAPSTER_ENC:") {
		return "", fmt.Errorf("not encrypted")
	}
	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(enc, "NAPSTER_ENC:"))
	if err != nil {
		return "", fmt.Errorf("base64 نامعتبر")
	}
	if len(data) < 16 {
		return "", fmt.Errorf("داده کوتاه است")
	}
	key := deriveKey(password, data[:16])
	plain, err := gcmDecrypt(data[16:], key)
	if err != nil {
		return "", fmt.Errorf("رمز اشتباه است")
	}
	return string(plain), nil
}

// ─── رمزگشایی فرمت NPVT1 ─────────────────────────────────────────────────────
//
// فرمت: NPVT1 <salt_b64>,<enc_config_b64>,<enc_v2link_b64>
//
// هر بخش رمزشده = nonce(12 بایت) + ciphertext + tag(16 بایت)
// کلید = deriveKey(password, salt_bytes)
//
func decryptNPVT1(content, password string) (configYAML, v2link string, err error) {
	if password == "" {
		err = fmt.Errorf("NEEDS_PASSWORD")
		return
	}

	// حذف پیشوند "NPVT1"
	body := strings.TrimSpace(content[5:]) // بعد از "NPVT1"

	parts := strings.SplitN(body, ",", 3)
	if len(parts) < 2 {
		err = fmt.Errorf("فرمت NPVT1 نامعتبر")
		return
	}

	saltBytes, e := base64.StdEncoding.DecodeString(strings.TrimSpace(parts[0]))
	if e != nil {
		err = fmt.Errorf("salt نامعتبر در NPVT1")
		return
	}

	key := deriveKey(password, saltBytes)

	// رمزگشایی config (part2)
	ctBytes, e := base64.StdEncoding.DecodeString(strings.TrimSpace(parts[1]))
	if e != nil {
		err = fmt.Errorf("ciphertext نامعتبر در NPVT1")
		return
	}
	plainConfig, e := gcmDecrypt(ctBytes, key)
	if e != nil {
		err = fmt.Errorf("رمز اشتباه است")
		return
	}
	configYAML = string(plainConfig)

	// رمزگشایی v2ray link (part3) اگر وجود داشت
	if len(parts) == 3 && strings.TrimSpace(parts[2]) != "" {
		v2Bytes, e := base64.StdEncoding.DecodeString(strings.TrimSpace(parts[2]))
		if e == nil {
			plainV2, e2 := gcmDecrypt(v2Bytes, key)
			if e2 == nil {
				v2link = strings.TrimSpace(string(plainV2))
			}
		}
	}
	return
}

// ─── V2Ray Parsers ────────────────────────────────────────────────────────────

func tryDecodeBase64(s string) ([]byte, error) {
	for _, enc := range []base64.Encoding{
		*base64.StdEncoding, *base64.RawStdEncoding,
		*base64.URLEncoding, *base64.RawURLEncoding,
	} {
		e := enc
		if b, err := e.DecodeString(s); err == nil {
			return b, nil
		}
	}
	return nil, fmt.Errorf("base64 نامعتبر")
}

func parseVMess(link string) (V2RayConfig, error) {
	decoded, err := tryDecodeBase64(strings.TrimPrefix(link, "vmess://"))
	if err != nil {
		return V2RayConfig{}, fmt.Errorf("VMess base64 نامعتبر")
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(decoded, &raw); err != nil {
		return V2RayConfig{}, fmt.Errorf("VMess JSON نامعتبر")
	}
	cfg := V2RayConfig{Protocol: "vmess"}
	gs := func(k string) string {
		if v, ok := raw[k].(string); ok { return v }
		return ""
	}
	gi := func(k string) int {
		if v, ok := raw[k]; ok {
			switch x := v.(type) {
			case float64: return int(x)
			case string:
				var n int; fmt.Sscanf(x, "%d", &n); return n
			}
		}
		return 0
	}
	cfg.Address = gs("add"); cfg.Port = gi("port"); cfg.UUID = gs("id")
	cfg.AlterId = gi("aid"); cfg.Security = gs("scy"); cfg.Network = gs("net")
	cfg.Path = gs("path"); cfg.Host = gs("host"); cfg.TLS = gs("tls")
	cfg.SNI = gs("sni"); cfg.ALPN = gs("alpn"); cfg.Remarks = gs("ps")
	cfg.Fingerprint = gs("fp")
	if cfg.Security == "" { cfg.Security = "auto" }
	if cfg.Network == "" { cfg.Network = "tcp" }
	return cfg, nil
}

func parseVLess(link string) (V2RayConfig, error) {
	u, err := url.Parse(link)
	if err != nil { return V2RayConfig{}, err }
	cfg := V2RayConfig{Protocol: "vless"}
	cfg.UUID = u.User.Username()
	cfg.Address = u.Hostname()
	fmt.Sscanf(u.Port(), "%d", &cfg.Port)
	cfg.Remarks, _ = url.QueryUnescape(u.Fragment)
	q := u.Query()
	cfg.Network = q.Get("type"); cfg.TLS = q.Get("security")
	cfg.SNI = q.Get("sni"); cfg.ALPN = q.Get("alpn")
	cfg.Path = q.Get("path"); cfg.Host = q.Get("host")
	cfg.Flow = q.Get("flow"); cfg.Fingerprint = q.Get("fp")
	cfg.PublicKey = q.Get("pbk"); cfg.ShortID = q.Get("sid")
	if cfg.Network == "" { cfg.Network = "tcp" }
	return cfg, nil
}

func parseTrojan(link string) (V2RayConfig, error) {
	u, err := url.Parse(link)
	if err != nil { return V2RayConfig{}, err }
	cfg := V2RayConfig{Protocol: "trojan"}
	cfg.Password = u.User.Username()
	cfg.Address = u.Hostname()
	fmt.Sscanf(u.Port(), "%d", &cfg.Port)
	cfg.Remarks, _ = url.QueryUnescape(u.Fragment)
	q := u.Query()
	cfg.Network = q.Get("type"); cfg.TLS = q.Get("security")
	cfg.SNI = q.Get("sni"); cfg.ALPN = q.Get("alpn")
	cfg.Path = q.Get("path"); cfg.Host = q.Get("host")
	cfg.Fingerprint = q.Get("fp")
	if cfg.Network == "" { cfg.Network = "tcp" }
	if cfg.TLS == "" { cfg.TLS = "tls" }
	return cfg, nil
}

func parseV2RayLink(link string) (V2RayConfig, error) {
	link = strings.TrimSpace(link)
	switch {
	case strings.HasPrefix(link, "vmess://"): return parseVMess(link)
	case strings.HasPrefix(link, "vless://"): return parseVLess(link)
	case strings.HasPrefix(link, "trojan://"): return parseTrojan(link)
	default: return V2RayConfig{}, fmt.Errorf("پروتکل پشتیبانی نمیشه (vmess / vless / trojan)")
	}
}

func toV2RayLink(cfg V2RayConfig) string {
	switch cfg.Protocol {
	case "vmess":
		raw := map[string]interface{}{
			"v": "2", "ps": cfg.Remarks, "add": cfg.Address, "port": cfg.Port,
			"id": cfg.UUID, "aid": cfg.AlterId, "scy": cfg.Security,
			"net": cfg.Network, "path": cfg.Path, "host": cfg.Host,
			"tls": cfg.TLS, "sni": cfg.SNI, "alpn": cfg.ALPN, "fp": cfg.Fingerprint,
		}
		b, _ := json.Marshal(raw)
		return "vmess://" + base64.StdEncoding.EncodeToString(b)
	case "vless":
		u := url.URL{Scheme: "vless", User: url.User(cfg.UUID),
			Host: fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), Fragment: cfg.Remarks}
		q := url.Values{}
		if cfg.Network != "" { q.Set("type", cfg.Network) }
		if cfg.TLS != "" { q.Set("security", cfg.TLS) }
		if cfg.SNI != "" { q.Set("sni", cfg.SNI) }
		if cfg.ALPN != "" { q.Set("alpn", cfg.ALPN) }
		if cfg.Path != "" { q.Set("path", cfg.Path) }
		if cfg.Host != "" { q.Set("host", cfg.Host) }
		if cfg.Flow != "" { q.Set("flow", cfg.Flow) }
		if cfg.Fingerprint != "" { q.Set("fp", cfg.Fingerprint) }
		if cfg.PublicKey != "" { q.Set("pbk", cfg.PublicKey) }
		if cfg.ShortID != "" { q.Set("sid", cfg.ShortID) }
		u.RawQuery = q.Encode()
		return u.String()
	case "trojan":
		u := url.URL{Scheme: "trojan", User: url.User(cfg.Password),
			Host: fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), Fragment: cfg.Remarks}
		q := url.Values{}
		if cfg.Network != "" { q.Set("type", cfg.Network) }
		if cfg.TLS != "" { q.Set("security", cfg.TLS) }
		if cfg.SNI != "" { q.Set("sni", cfg.SNI) }
		if cfg.ALPN != "" { q.Set("alpn", cfg.ALPN) }
		if cfg.Path != "" { q.Set("path", cfg.Path) }
		if cfg.Host != "" { q.Set("host", cfg.Host) }
		if cfg.Fingerprint != "" { q.Set("fp", cfg.Fingerprint) }
		u.RawQuery = q.Encode()
		return u.String()
	}
	return ""
}

// ─── Config Builder ───────────────────────────────────────────────────────────

func buildNapsterConfig(p V2RayConfig, s NapsterSettings) string {
	var b strings.Builder
	w := func(f string, a ...interface{}) { fmt.Fprintf(&b, f, a...) }

	mode := s.ProxyMode
	if mode == "" { mode = "rule" }
	logLevel := s.LogLevel
	if logLevel == "" { logLevel = "warning" }

	w("# Napster Config Generator v3\n")
	w("# Generated: %s\n", time.Now().Format("2006-01-02 15:04"))
	w("\nmixed-port: 7890\nsocks-port: 7891\nport: 7892\nredir-port: 7893\n\n")
	w("mode: %s\nallow-lan: false\nlog-level: %s\nipv6: %v\n\n", mode, logLevel, s.EnableIPv6)
	w("external-controller: 127.0.0.1:9090\nsecret: \"\"\n\n")

	name := p.Remarks
	if name == "" {
		name = fmt.Sprintf("%s-%s-%d", strings.ToUpper(p.Protocol), p.Address, p.Port)
	}

	w("dns:\n  enable: true\n  enhanced-mode: fake-ip\n  listen: 0.0.0.0:53\n")
	w("  use-hosts: true\n  respect-rules: true\n")
	w("  fake-ip-range: 198.18.0.0/15\n  fake-ip-filter:\n    - '.lan'\n    - '.local'\n")
	w("  nameserver:\n    - 1.1.1.1\n")
	w("  fallback:\n    - 8.8.8.8\n  fallback-filter:\n    geoip: true\n    geoip-code: IR\n    ipcidr:\n      - 240.0.0.0/4\n\n")

	w("sniffer:\n  enable: true\n  sniff:\n    HTTP:\n      ports: [80, 8080-8880]\n    TLS:\n      ports: [443, 8443]\n    QUIC:\n      ports: [443]\n\n")

	w("proxies:\n")
	fp := p.Fingerprint
	if fp == "" { fp = "chrome" }

	switch p.Protocol {
	case "vmess":
		w("  - name: \"%s\"\n    type: vmess\n    server: %s\n    port: %d\n", name, p.Address, p.Port)
		w("    uuid: %s\n    alterId: %d\n    cipher: %s\n    udp: true\n", p.UUID, p.AlterId, p.Security)
		if p.TLS == "tls" {
			w("    tls: true\n    skip-cert-verify: false\n")
			if p.SNI != "" { w("    servername: %s\n", p.SNI) }
		}
		w("    client-fingerprint: %s\n", fp)
		writeNetworkOpts(&b, p)
	case "vless":
		w("  - name: \"%s\"\n    type: vless\n    server: %s\n    port: %d\n", name, p.Address, p.Port)
		w("    uuid: %s\n    udp: true\n", p.UUID)
		if p.Flow != "" { w("    flow: %s\n", p.Flow) }
		if p.TLS == "tls" || p.TLS == "reality" {
			w("    tls: true\n    skip-cert-verify: false\n")
			if p.SNI != "" { w("    servername: %s\n", p.SNI) }
		}
		if p.TLS == "reality" {
			w("    reality-opts:\n      public-key: %s\n", p.PublicKey)
			if p.ShortID != "" { w("      short-id: %s\n", p.ShortID) }
		}
		w("    client-fingerprint: %s\n", fp)
		writeNetworkOpts(&b, p)
	case "trojan":
		w("  - name: \"%s\"\n    type: trojan\n    server: %s\n    port: %d\n", name, p.Address, p.Port)
		w("    password: %s\n    udp: true\n    tls: true\n    skip-cert-verify: false\n", p.Password)
		if p.SNI != "" { w("    sni: %s\n", p.SNI) }
		w("    client-fingerprint: %s\n", fp)
		writeNetworkOpts(&b, p)
	}

	w("\nproxy-groups:\n")
	w("  - name: \"Proxy\"\n    type: select\n    proxies:\n      - \"%s\"\n      - DIRECT\n\n", name)
	w("  - name: \"Auto\"\n    type: url-test\n    proxies:\n      - \"%s\"\n    url: http://www.gstatic.com/generate_204\n    interval: 300\n    tolerance: 50\n\n", name)

	w("rules:\n")
	for _, ip := range []string{"127.0.0.0/8", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} {
		w("  - IP-CIDR,%s,DIRECT,no-resolve\n", ip)
	}
	if s.EnableBypass {
		for _, d := range s.BypassDomains {
			if d = strings.TrimSpace(d); d != "" {
				w("  - DOMAIN-SUFFIX,%s,DIRECT\n", d)
			}
		}
		for _, ip := range s.BypassIPs {
			if ip = strings.TrimSpace(ip); ip != "" {
				if net.ParseIP(ip) != nil {
					w("  - IP-CIDR,%s/32,DIRECT,no-resolve\n", ip)
				} else {
					w("  - IP-CIDR,%s,DIRECT,no-resolve\n", ip)
				}
			}
		}
	}
	w("  - GEOIP,IR,DIRECT,no-resolve\n")
	w("  - GEOIP,private,DIRECT,no-resolve\n")
	w("  - MATCH,Proxy\n")
	return b.String()
}

func writeNetworkOpts(b *strings.Builder, p V2RayConfig) {
	if p.Network == "" || p.Network == "tcp" { return }
	fmt.Fprintf(b, "    network: %s\n", p.Network)
	switch p.Network {
	case "ws":
		fmt.Fprintf(b, "    ws-opts:\n      path: \"%s\"\n", p.Path)
		if p.Host != "" { fmt.Fprintf(b, "      headers:\n        Host: \"%s\"\n", p.Host) }
	case "grpc":
		fmt.Fprintf(b, "    grpc-opts:\n      grpc-service-name: \"%s\"\n", p.Path)
	case "h2":
		fmt.Fprintf(b, "    h2-opts:\n      path: \"%s\"\n", p.Path)
		if p.Host != "" { fmt.Fprintf(b, "      host:\n        - \"%s\"\n", p.Host) }
	case "httpupgrade":
		fmt.Fprintf(b, "    httpupgrade-opts:\n      path: \"%s\"\n", p.Path)
		if p.Host != "" { fmt.Fprintf(b, "      host: \"%s\"\n", p.Host) }
	}
}

// ─── ساخت NPVT ───────────────────────────────────────────────────────────────

func buildNpvt(p V2RayConfig, s NapsterSettings, cfgText string) (string, error) {
	npvt := NpvtFile{
		Version:     3,
		CreatedAt:   time.Now().Format("2006-01-02T15:04:05Z"),
		UserAgent:   s.UserAgent,
		V2RayLink:   toV2RayLink(p),
		ParsedProxy: p,
		Settings:    s,
		Config:      cfgText,
		Encrypted:   false,
	}
	if s.EnableDeviceLock && s.DeviceID != "" {
		npvt.DeviceLock = s.DeviceID
	}

	jsonBytes, err := json.MarshalIndent(npvt, "", "  ")
	if err != nil { return "", err }

	if s.EnablePassword && s.Password != "" {
		encStr, err := encrypt(string(jsonBytes), s.Password)
		if err != nil { return "", err }
		wrapper := map[string]interface{}{
			"version": 3, "encrypted": true, "data": encStr,
		}
		wb, _ := json.Marshal(wrapper)
		return base64.StdEncoding.EncodeToString(wb), nil
	}

	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

// ─── Parse NPVT (هر دو فرمت) ─────────────────────────────────────────────────

func parseNpvt(content, password string) (*NpvtFile, string, error) {
	content = strings.TrimSpace(content)

	// فرمت NPVT1 — از اپلیکیشن Napster
	if strings.HasPrefix(content, "NPVT1") {
		cfgYAML, v2link, err := decryptNPVT1(content, password)
		if err != nil { return nil, "", err }
		synthetic := &NpvtFile{
			Version:   1,
			CreatedAt: time.Now().Format("2006-01-02T15:04:05Z"),
			Config:    cfgYAML,
			V2RayLink: v2link,
			Encrypted: true,
		}
		if v2link != "" {
			if cfg, err := parseV2RayLink(v2link); err == nil {
				synthetic.ParsedProxy = cfg
			}
		}
		return synthetic, v2link, nil
	}

	// فرمت base64 JSON — از این سرور
	raw, err := tryDecodeBase64(content)
	if err != nil {
		return nil, "", fmt.Errorf("فایل نامعتبر — نه NPVT1 است نه base64")
	}

	var wrapper map[string]interface{}
	if err := json.Unmarshal(raw, &wrapper); err != nil {
		return nil, "", fmt.Errorf("JSON نامعتبر در فایل")
	}

	var jsonStr string
	if enc, ok := wrapper["encrypted"].(bool); ok && enc {
		if password == "" { return nil, "", fmt.Errorf("NEEDS_PASSWORD") }
		dataStr, ok := wrapper["data"].(string)
		if !ok { return nil, "", fmt.Errorf("ساختار فایل رمزشده نامعتبر") }
		decrypted, err := decrypt(dataStr, password)
		if err != nil { return nil, "", err }
		jsonStr = decrypted
	} else {
		jsonStr = string(raw)
	}

	var npvt NpvtFile
	if err := json.Unmarshal([]byte(jsonStr), &npvt); err != nil {
		return nil, "", fmt.Errorf("ساختار فایل npvt نامعتبر")
	}

	link := npvt.V2RayLink
	if link == "" { link = toV2RayLink(npvt.ParsedProxy) }
	return &npvt, link, nil
}

// ─── Server Info ──────────────────────────────────────────────────────────────

func getServerInfo(address string) map[string]string {
	info := map[string]string{
		"address": address, "status": "offline", "ping": "N/A", "country": "N/A", "org": "N/A",
	}
	for _, port := range []string{"443", "80"} {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", address+":"+port, 3*time.Second)
		if err == nil {
			conn.Close()
			info["ping"] = fmt.Sprintf("%dms", time.Since(start).Milliseconds())
			info["status"] = "online"
			break
		}
	}
	cl := &http.Client{Timeout: 4 * time.Second}
	if resp, err := cl.Get("http://ip-api.com/json/" + address + "?fields=country,countryCode,org"); err == nil {
		defer resp.Body.Close()
		var r map[string]interface{}
		if json.NewDecoder(resp.Body).Decode(&r) == nil {
			if v, ok := r["country"].(string); ok { info["country"] = v }
			if v, ok := r["countryCode"].(string); ok { info["countryCode"] = v }
			if v, ok := r["org"].(string); ok { info["org"] = v }
		}
	}
	return info
}

// ─── HTTP Handlers ────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexHTML)
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	cfg, err := parseV2RayLink(r.URL.Query().Get("link"))
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, cfg)
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, GenerateResponse{Error: "JSON نامعتبر"})
		return
	}
	p, err := parseV2RayLink(req.V2RayLink)
	if err != nil {
		writeJSON(w, GenerateResponse{Error: err.Error()})
		return
	}

	cfgText := buildNapsterConfig(p, req.Settings)
	resp := GenerateResponse{Success: true, Config: cfgText}

	npvtB64, err := buildNpvt(p, req.Settings, cfgText)
	if err == nil { resp.NpvtB64 = npvtB64 }

	appData := loadData()
	prof := Profile{
		ID: uuid.New().String(), Name: req.ProfileName,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		V2RayLink: req.V2RayLink, ParsedProxy: p,
		Settings: req.Settings, Config: cfgText,
	}
	if prof.Name == "" {
		prof.Name = fmt.Sprintf("پروفایل %d", len(appData.Profiles)+1)
	}
	appData.Profiles = append(appData.Profiles, prof)
	appData.LastSettings = req.Settings
	saveData(appData)
	resp.ProfileID = prof.ID
	writeJSON(w, resp)
}

func handleProfiles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		d := loadData()
		writeJSON(w, map[string]interface{}{"profiles": d.Profiles, "lastSettings": d.LastSettings})
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		d := loadData()
		np := []Profile{}
		for _, p := range d.Profiles {
			if p.ID != id { np = append(np, p) }
		}
		d.Profiles = np
		saveData(d)
		writeJSON(w, map[string]bool{"ok": true})
	default:
		w.WriteHeader(405)
	}
}

func handleImportNpvt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var req struct {
		Content  string `json:"content"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	npvt, v2link, err := parseNpvt(strings.TrimSpace(req.Content), req.Password)
	if err != nil {
		if err.Error() == "NEEDS_PASSWORD" {
			writeJSON(w, map[string]string{"error": "NEEDS_PASSWORD"})
			return
		}
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]interface{}{
		"success": true, "v2link": v2link,
		"config": npvt.Config, "proxy": npvt.ParsedProxy, "settings": npvt.Settings,
		"meta": map[string]interface{}{
			"version": npvt.Version, "createdAt": npvt.CreatedAt,
			"deviceLock": npvt.DeviceLock, "userAgent": npvt.UserAgent, "encrypted": npvt.Encrypted,
		},
	})
}

func handleServerInfo(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("address")
	if addr == "" {
		writeJSON(w, map[string]string{"error": "no address"})
		return
	}
	writeJSON(w, getServerInfo(addr))
}

func handleDeviceID(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"deviceId": uuid.New().String()})
}

// ─── Main ─────────────────────────────────────────────────────────────────────

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/parse", handleParse)
	mux.HandleFunc("/api/generate", handleGenerate)
	mux.HandleFunc("/api/profiles", handleProfiles)
	mux.HandleFunc("/api/import-npvt", handleImportNpvt)
	mux.HandleFunc("/api/server-info", handleServerInfo)
	mux.HandleFunc("/api/device-id", handleDeviceID)

	port := "8080"
	if p := os.Getenv("PORT"); p != "" { port = p }
	log.Printf("Napster Config Generator -> http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
