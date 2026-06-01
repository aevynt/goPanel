package apps

type AppInfo struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Image       string `json:"image"`
	DefaultPort int    `json:"default_port"`
	Template    string `json:"-"`
}

var Catalog = []AppInfo{
	{
		Name:        "Pi-hole",
		Key:         "pihole",
		Description: "A black hole for Internet advertisements",
		Category:    "Network",
		Image:       "pihole/pihole:latest",
		DefaultPort: 8088,
		Template: `version: '3.8'
services:
  pihole:
    container_name: pihole
    image: pihole/pihole:latest
    ports:
      - "${PORT}:80/tcp"
    environment:
      TZ: 'Asia/Ho_Chi_Minh'
      WEBPASSWORD: 'admin'
    restart: unless-stopped`,
	},
	{
		Name:        "Plex Media Server",
		Key:         "plex",
		Description: "Stream your personal media library anywhere",
		Category:    "Media",
		Image:       "plexinc/pms-docker:latest",
		DefaultPort: 32400,
		Template: `version: '3.8'
services:
  plex:
    container_name: plex
    image: plexinc/pms-docker:latest
    ports:
      - "${PORT}:32400/tcp"
    environment:
      TZ: 'Asia/Ho_Chi_Minh'
    restart: unless-stopped`,
	},
	{
		Name:        "Home Assistant",
		Key:         "homeassistant",
		Description: "Open source home automation that puts local control first",
		Category:    "Smart Home",
		Image:       "ghcr.io/home-assistant/home-assistant:stable",
		DefaultPort: 8123,
		Template: `version: '3.8'
services:
  homeassistant:
    container_name: homeassistant
    image: ghcr.io/home-assistant/home-assistant:stable
    ports:
      - "${PORT}:8123/tcp"
    environment:
      TZ: 'Asia/Ho_Chi_Minh'
    restart: unless-stopped`,
	},
	{
		Name:        "Vaultwarden",
		Key:         "vaultwarden",
		Description: "Unofficial Bitwarden compatible server written in Rust",
		Category:    "Security",
		Image:       "vaultwarden/server:latest",
		DefaultPort: 8080,
		Template: `version: '3.8'
services:
  vaultwarden:
    container_name: vaultwarden
    image: vaultwarden/server:latest
    ports:
      - "${PORT}:80/tcp"
    environment:
      WEBSOCKET_ENABLED: 'true'
    restart: unless-stopped`,
	},
	{
		Name:        "AdGuard Home",
		Key:         "adguardhome",
		Description: "Network-wide software for blocking ads & tracking",
		Category:    "Network",
		Image:       "adguard/adguardhome:latest",
		DefaultPort: 3000,
		Template: `version: '3.8'
services:
  adguardhome:
    container_name: adguardhome
    image: adguard/adguardhome:latest
    ports:
      - "${PORT}:3000/tcp"
    restart: unless-stopped`,
	},
	{
		Name:        "qBittorrent",
		Key:         "qbittorrent",
		Description: "Free and open-source BitTorrent client",
		Category:    "Utilities",
		Image:       "linuxserver/qbittorrent:latest",
		DefaultPort: 8090,
		Template: `version: '3.8'
services:
  qbittorrent:
    container_name: qbittorrent
    image: linuxserver/qbittorrent:latest
    ports:
      - "${PORT}:8080/tcp"
    environment:
      PUID: '1000'
      PGID: '1000'
      TZ: 'Asia/Ho_Chi_Minh'
      WEBUI_PORT: '8080'
    restart: unless-stopped`,
	},
}
