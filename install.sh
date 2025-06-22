#!/bin/sh
# Usage: curl -fsSL https://github.com/austin-weeks/tv-control-server/install.sh | sh

REPO="austin-weeks/tv-control-server"
NAME="tv-control"

GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS
OS=$(uname | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux*) OS=linux ;;
  darwin*) OS=darwin ;;
  msys*|mingw*|cygwin*|windows*) OS=windows ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect ARCH
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH=amd64 ;;
  arm64|aarch64) ARCH=arm64 ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Set output filename
if [ "$OS" = "windows" ]; then
  OUT="$NAME.exe"
else
  OUT="$NAME"
fi

# Get latest release tag
TAG=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | head -1 | sed -E 's/.*: "([^"]+)".*/\1/
')
if [ -z "$TAG" ]; then
  echo "Could not find latest release tag."; exit 1
fi

# Download binary
URL="https://github.com/$REPO/releases/download/$TAG/$NAME-$OS-$ARCH${OUT##$NAME}"
printf "${CYAN}Downloading tag $TAG at $URL ...${NC}\n\n"
curl -fL --progress-bar "$URL" -o "$OUT"
if [ $? -ne 0 ]; then
  echo "Download failed. Please check the URL or your network connection."
  exit 1
fi
chmod +x "$OUT"
printf "\n${GREEN}Installed $OUT to $(pwd)${NC}\n"

# Create default config.json if it doesn't exist
CONFIG_FILE="config.json"
if [ ! -f "$CONFIG_FILE" ]; then
  cat > "$CONFIG_FILE" <<EOF
{
    "tv_ip": "",
    "app_name": "Gopher Remote",
    "app_port": "1234",
    "token_file": ".tv_token",
    "tv_port": "8002",
    "client_password": "",
    "brightness_location": 3,
    "initial_delay_ms": 2000
}
EOF
  printf "\n${GREEN}Created $CONFIG_FILE with default settings.${NC}\n"
fi

# Check if tv_ip is set in config.json
TV_IP="$(awk -F '"tv_ip"[[:space:]]*:[[:space:]]*"' '{if (NF>1) print $2}' "$CONFIG_FILE" | awk -F'"' '{print $1}' | head -1)"
if [ -z "$TV_IP" ]; then
  printf "\n${YELLOW}Please add your TV's IP address to config.json (tv_ip) before running the server.${NC}\n"
fi

