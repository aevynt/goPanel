# goPanel 🚀

> A premium, lightweight, and modern glassmorphic Server Control Panel designed specifically for Home Labs. Manage your services, containers, networks, and configurations in one beautiful dashboard.

---

## Key Features

- 🖥️ **Systemd & Windows Service Manager**: Easily view, start, stop, enable, or disable system services and custom binaries.
- 🐳 **Docker & Compose Management**: Monitor container states, view live container logs, and deploy multi-container projects directly via `docker-compose.yml` sheets.
- 🌐 **Dynamic Caddy Reverse Proxy & SSL**: Configure sites dynamically using Caddy’s REST API. Auto-provision HTTPS certificates securely via Let's Encrypt with clear visual shield checkmark badges.
- 📂 **Secure File Manager**: Access your files, upload logs, and perform contextual **Zip Folder** or **Extract Archive** actions natively with path-traversal (Zip Slip) security.
- ⚡ **Curated 1-Click App Store**: Asynchronously install popular homelab apps (Pi-hole, Plex, Home Assistant, Vaultwarden, AdGuard Home, qBittorrent) on free host ports with dynamic reverse proxies mapped in seconds.
- 📊 **Resource Telemetry & CPU Temp**: Real-time CPU, RAM, Disk, and Load telemetry, including Linux CPU temperature sensors (`/sys/class/thermal`) styled in modern glassmorphic progress dials.
- 🔔 **Discord & Telegram Alerts**: Run background health supervisors checking resource utilization and temperature every 60s, dispatching alerts to webhooks with a built-in 30-minute anti-spam cooldown.
- 🔒 **Enhanced Security**: Secure authorization, instant Admin Password Changer, and standard Time-based Two-Factor Authentication (TOTP 2FA) built entirely using the Go standard library for zero-dependency binary optimization.

---

## Installation & Getting Started

### Quick Linux Installer (Recommended)
You can install and run `goPanel` on any modern Linux system using our automated setup script:
```bash
curl -fsSL https://raw.githubusercontent.com/aevynt/goPanel/master/install.sh | bash
```

### Build from Source
If you want to compile `goPanel` locally, follow these steps:

#### 1. Prerequisites
- **Go** (version 1.22+)
- **Node.js** & **npm** (for the frontend assets compilation)

#### 2. Compile the Vue Frontend
```bash
cd web
npm install
npm run build
cd ..
```

#### 3. Compile the Go Backend
```bash
# This embeds the built dist files and compiles the backend binary
go build -o gopanel .
```

#### 4. Run the Panel
```bash
./gopanel -config ./gopanel.json
```
Open `http://localhost:3636` in your browser and sign in using the default credentials:
- **Username**: `admin`
- **Password**: `admin`

---

## Configuration Settings

You can customize `goPanel` using the `gopanel.json` file. Here is a list of parameters:

| Parameter | Type | Default | Description |
|---|---|---|---|
| `port` | number | `3636` | The port the control panel listens on for HTTP traffic. |
| `data_dir` | string | `/var/lib/gopanel` | Directory where databases, shares, and compose projects are saved. |
| `binaries_dir` | string | `/var/lib/gopanel/binaries` | Folder used to store custom uploaded binaries. |
| `jwt_secret` | string | *Auto-generated* | Cryptographic token signature key. |
| `caddy_admin_url` | string | `http://localhost:2019` | The admin REST API endpoint of your Caddy daemon. |
| `panel_domain` | string | `""` | Global domain used to access your panel. |
| `discord_webhook` | string | `""` | Discord channel webhook URL for system alert notifications. |
| `telegram_token` | string | `""` | Telegram Bot API Token for alerting. |
| `telegram_chat_id` | string | `""` | Telegram target chat or group ID. |
| `alert_temp_threshold` | number | `80` | CPU temperature limit (°C) that triggers alert notifications. |

---

## Development & Contribution

Contributions are welcome! Please read our [CONTRIBUTING.md](CONTRIBUTING.md) to understand our coding standards, branch setups, and pull request flows.

## License

This project is open-source and licensed under the permissive terms of the [MIT License](LICENSE).
