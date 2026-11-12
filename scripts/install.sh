#!/usr/bin/env sh
# Halo 一键安装脚本。
# 用法：
#   curl -fsSL https://raw.githubusercontent.com/zfxzfxzfx852/halo/main/scripts/install.sh | sh
# 可通过环境变量覆盖：
#   HALO_VERSION=v0.9.0   指定版本（默认 latest）
#   HALO_INSTALL_DIR=...  安装目录（默认 /usr/local/bin）
set -e

REPO="zfxzfxzfx852/halo"
VERSION="${HALO_VERSION:-latest}"
INSTALL_DIR="${HALO_INSTALL_DIR:-/usr/local/bin}"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
  linux*) os=linux ;;
  darwin*) os=darwin ;;
  *) echo "unsupported OS: $os" >&2; exit 1 ;;
esac

arch=$(uname -m)
case "$arch" in
  x86_64 | amd64) arch=amd64 ;;
  arm64 | aarch64) arch=arm64 ;;
  *) echo "unsupported architecture: $arch" >&2; exit 1 ;;
esac

if [ "$VERSION" = "latest" ]; then
  VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | grep '"tag_name"' | head -1 | sed 's/.*: *"//; s/".*//')
fi

asset="halo-${VERSION}-${os}-${arch}.tar.gz"
url="https://github.com/${REPO}/releases/download/${VERSION}/${asset}"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading Halo ${VERSION} (${os}/${arch})…"
curl -fsSL "$url" -o "$tmp/halo.tar.gz"
tar -xzf "$tmp/halo.tar.gz" -C "$tmp"
bin=$(find "$tmp" -name halo -type f | head -1)

if [ ! -w "$INSTALL_DIR" ]; then
  echo "需要写入 $INSTALL_DIR，可能提示输入 sudo 密码。"
  SUDO=sudo
else
  SUDO=""
fi
$SUDO mkdir -p "$INSTALL_DIR"
$SUDO mv "$bin" "$INSTALL_DIR/halo"
$SUDO chmod +x "$INSTALL_DIR/halo"

echo "Halo 已安装到 $INSTALL_DIR/halo"
echo "运行 'halo init' 生成配置，或直接 'halo' 启动后访问 http://127.0.0.1:8787/console/"
