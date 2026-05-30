#!/bin/bash
set -e

PANEL_BIN="/usr/local/bin/gopanel"
PANEL_CONFIG="/etc/gopanel.json"
PANEL_DATA="/var/lib/gopanel"
PANEL_SERVICE="/etc/systemd/system/gopanel.service"

echo "==> Building goPanel..."
cd "$(dirname "$0")"

if [ ! -d "web/dist" ]; then
  echo "Building frontend..."
  cd web && npm install && npm run build && cd ..
fi

echo "Building Go binary..."
go build -ldflags="-s -w" -o gopanel .

echo "==> Installing binary..."
install -m 755 gopanel "$PANEL_BIN"

echo "==> Creating data directories..."
mkdir -p "$PANEL_DATA/binaries"
mkdir -p "$PANEL_DATA/logs"

echo "==> Creating systemd service..."
cat > "$PANEL_SERVICE" << 'EOF'
[Unit]
Description=goPanel - Server Management Panel
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/gopanel --config /etc/gopanel.json
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload

echo "==> Starting goPanel..."
systemctl enable gopanel
systemctl start gopanel

echo "==> Install complete!"
echo "goPanel is running on http://localhost:8080"
echo "Default login: admin / admin"
echo ""
echo "You can change the config at: $PANEL_CONFIG"
