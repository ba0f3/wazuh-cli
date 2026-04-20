#!/usr/bin/env sh
# install.sh — one-liner installer for wazuh-cli
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/ba0f3/wazuh-cli/main/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/ba0f3/wazuh-cli/main/install.sh | sh -s -- --version v1.2.3
#   curl -fsSL https://raw.githubusercontent.com/ba0f3/wazuh-cli/main/install.sh | sh -s -- --dir ~/.local/bin
set -e

REPO="ba0f3/wazuh-cli"
BINARY="wazuh-cli"
INSTALL_DIR=""
REQUESTED_VERSION=""

# ── parse flags ────────────────────────────────────────────────────────────────
while [ $# -gt 0 ]; do
  case "$1" in
    --version|-v) REQUESTED_VERSION="$2"; shift 2 ;;
    --dir|-d)     INSTALL_DIR="$2";        shift 2 ;;
    --help|-h)
      echo "Usage: install.sh [--version <tag>] [--dir <install-dir>]"
      exit 0 ;;
    *) echo "Unknown option: $1" >&2; exit 1 ;;
  esac
done

# ── colour helpers ─────────────────────────────────────────────────────────────
red()   { printf '\033[31m%s\033[0m\n' "$*"; }
green() { printf '\033[32m%s\033[0m\n' "$*"; }
bold()  { printf '\033[1m%s\033[0m\n'  "$*"; }

# ── detect OS ─────────────────────────────────────────────────────────────────
OS="$(uname -s)"
case "$OS" in
  Linux)   OS_NAME="Linux"  ;;
  Darwin)  OS_NAME="Darwin" ;;
  MINGW*|MSYS*|CYGWIN*) OS_NAME="Windows" ;;
  *)
    red "Unsupported OS: $OS"
    echo "Please download manually from https://github.com/$REPO/releases" >&2
    exit 1
    ;;
esac

# ── detect architecture ───────────────────────────────────────────────────────
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64)         ARCH_NAME="x86_64"  ;;
  aarch64|arm64)        ARCH_NAME="arm64"   ;;
  armv7l|armv6l|armhf) ARCH_NAME="armv6"   ;;
  *)
    red "Unsupported architecture: $ARCH"
    echo "Please download manually from https://github.com/$REPO/releases" >&2
    exit 1
    ;;
esac

# ── resolve version ───────────────────────────────────────────────────────────
if [ -z "$REQUESTED_VERSION" ]; then
  bold "Fetching latest release version..."
  if command -v curl >/dev/null 2>&1; then
    LATEST=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
      | grep '"tag_name"' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
  elif command -v wget >/dev/null 2>&1; then
    LATEST=$(wget -qO- "https://api.github.com/repos/$REPO/releases/latest" \
      | grep '"tag_name"' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
  else
    red "Neither curl nor wget is available."
    exit 1
  fi
  VERSION="$LATEST"
else
  VERSION="$REQUESTED_VERSION"
fi

if [ -z "$VERSION" ]; then
  red "Could not determine version. Use --version <tag> to specify one."
  exit 1
fi

bold "Installing $BINARY $VERSION for $OS_NAME/$ARCH_NAME..."

# ── build download URLs ───────────────────────────────────────────────────────
EXT="tar.gz"
[ "$OS_NAME" = "Windows" ] && EXT="zip"

ARCHIVE="${BINARY}_${OS_NAME}_${ARCH_NAME}.${EXT}"
BASE_URL="https://github.com/$REPO/releases/download/$VERSION"
ARCHIVE_URL="$BASE_URL/$ARCHIVE"
CHECKSUM_URL="$BASE_URL/checksums.txt"

# ── download ──────────────────────────────────────────────────────────────────
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

download() {
  url="$1"; dest="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$dest"
  else
    wget -qO "$dest" "$url"
  fi
}

bold "Downloading $ARCHIVE..."
download "$ARCHIVE_URL"   "$TMP/$ARCHIVE"
download "$CHECKSUM_URL"  "$TMP/checksums.txt"

# ── verify checksum ───────────────────────────────────────────────────────────
bold "Verifying checksum..."
if command -v sha256sum >/dev/null 2>&1; then
  EXPECTED=$(grep "$ARCHIVE" "$TMP/checksums.txt" | awk '{print $1}')
  ACTUAL=$(sha256sum "$TMP/$ARCHIVE" | awk '{print $1}')
elif command -v shasum >/dev/null 2>&1; then
  EXPECTED=$(grep "$ARCHIVE" "$TMP/checksums.txt" | awk '{print $1}')
  ACTUAL=$(shasum -a 256 "$TMP/$ARCHIVE" | awk '{print $1}')
else
  echo "Warning: sha256sum/shasum not found, skipping checksum verification." >&2
  EXPECTED=""; ACTUAL=""
fi

if [ -n "$EXPECTED" ] && [ "$EXPECTED" != "$ACTUAL" ]; then
  red "Checksum mismatch!"
  echo "  expected: $EXPECTED" >&2
  echo "  actual:   $ACTUAL"   >&2
  exit 1
fi

# ── extract ───────────────────────────────────────────────────────────────────
bold "Extracting..."
case "$EXT" in
  tar.gz) tar -xzf "$TMP/$ARCHIVE" -C "$TMP" ;;
  zip)    unzip -q  "$TMP/$ARCHIVE" -d "$TMP" ;;
esac

BIN_SRC="$TMP/$BINARY"
[ "$OS_NAME" = "Windows" ] && BIN_SRC="$TMP/${BINARY}.exe"

# ── choose install dir ────────────────────────────────────────────────────────
if [ -z "$INSTALL_DIR" ]; then
  if [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
  elif [ -d "$HOME/.local/bin" ]; then
    INSTALL_DIR="$HOME/.local/bin"
  else
    INSTALL_DIR="$HOME/bin"
    mkdir -p "$INSTALL_DIR"
  fi
fi

# ── install ───────────────────────────────────────────────────────────────────
DEST="$INSTALL_DIR/$BINARY"
[ "$OS_NAME" = "Windows" ] && DEST="$INSTALL_DIR/${BINARY}.exe"

cp "$BIN_SRC" "$DEST"
chmod +x "$DEST"

green "✓ Installed $BINARY $VERSION → $DEST"

# ── PATH hint ─────────────────────────────────────────────────────────────────
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "  Add $INSTALL_DIR to your PATH:"
    echo "    export PATH=\"\$PATH:$INSTALL_DIR\""
    ;;
esac

echo ""
bold "Verify the installation:"
echo "  $BINARY --version"
echo "  $BINARY auth login --url https://wazuh:55000 --user admin --insecure"
