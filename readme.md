# ⚡ Napster Config Generator

تولیدکننده کانفیگ Napster/Clash از لینک‌های V2Ray با رابط وب زیبا

## ✨ امکانات

- **پشتیبانی از پروتکل‌ها:** VMess، VLess، Trojan
- **تنظیمات امنیتی کامل:** TLS 1.3، uTLS Fingerprint، Mux/Multiplexing
- **DNS پیشرفته:** Fake-IP، DNS over TLS، Fallback DNS
- **مسیریابی هوشمند:** Bypass ایران، Block دامنه‌ها، Rule/Global/Direct
- **Device Identity:** تولید خودکار Device ID، User-Agent سفارشی
- **رابط وب زیبا** به زبان فارسی
- **خروجی:** فرمت Napster و Clash

## 🚀 اجرا

### دانلود آماده (توصیه شده)

از [Releases](../../releases) فایل مناسب سیستم خود را دانلود کنید:

| سیستم‌عامل | فایل |
|---|---|
| Windows 64-bit | `napster-config-windows-amd64.exe` |
| Windows 32-bit | `napster-config-windows-386.exe` |
| Linux 64-bit | `napster-config-linux-amd64` |
| macOS (M1/M2/M3) | `napster-config-macos-arm64` |
| macOS (Intel) | `napster-config-macos-amd64` |

سپس اجرا کنید — مرورگر را باز کنید و به آدرس زیر بروید:
```
http://localhost:8080
```

### کامپایل دستی

```bash
go mod download
go build -ldflags="-s -w" -o napster-config .
./napster-config
```

## 📋 نحوه استفاده

1. فایل را اجرا کنید
2. مرورگر را باز کنید: `http://localhost:8080`
3. لینک V2Ray خود را وارد کنید (vmess/vless/trojan)
4. تنظیمات امنیتی، DNS، و مسیریابی را تنظیم کنید
5. روی **«تولید کانفیگ»** کلیک کنید
6. کانفیگ را کپی یا دانلود کنید

## 🔧 تنظیمات قابل تغییر

### امنیت
- TLS 1.3
- uTLS Browser Fingerprint (Chrome, Firefox, Safari, iOS, Android, Edge, ...)
- Mux/Multiplexing (smux, yamux, h2mux)

### DNS
- Fake-IP / Redir-Host / Normal
- DNS over TLS
- Custom Nameserver & Fallback

### مسیریابی
- Proxy Mode: Rule / Global / Direct
- Bypass لیست سفارشی برای سایت‌های ایرانی
- Block لیست برای مسدودسازی دامنه‌ها

### پیشرفته
- MTU سفارشی
- UDP Timeout
- IPv6 Toggle
- Sniffing
- Device ID و User-Agent سفارشی

## 🏗️ GitHub Actions

با هر push به branch اصلی، فایل‌های قابل اجرا برای تمام پلتفرم‌ها به صورت خودکار کامپایل می‌شوند.

برای دریافت فایل‌ها:
1. به تب **Actions** بروید
2. آخرین run را باز کنید
3. از بخش **Artifacts** فایل‌ها را دانلود کنید

یا یک **Release** بسازید تا فایل‌ها به آن attach شوند.
