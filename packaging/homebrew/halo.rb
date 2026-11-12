# Halo Homebrew Formula
# 发布时：运行 `make release checksums`，把各平台 tar.gz 的 sha256 填入下方，
# 然后同步到 homebrew-tap 仓库，用户即可 `brew install zfxzfxzfx852/tap/halo`。
class Halo < Formula
  desc "Local-first observability and security control plane for AI agents"
  homepage "https://github.com/zfxzfxzfx852/halo"
  version "0.9.0"
  license "MIT"
  head "https://github.com/zfxzfxzfx852/halo.git", branch: "main"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/zfxzfxzfx852/halo/releases/download/v0.9.0/halo-v0.9.0-darwin-arm64.tar.gz"
      sha256 "REPLACE_WITH_DARWIN_ARM64_SHA256"
    else
      url "https://github.com/zfxzfxzfx852/halo/releases/download/v0.9.0/halo-v0.9.0-darwin-amd64.tar.gz"
      sha256 "REPLACE_WITH_DARWIN_AMD64_SHA256"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/zfxzfxzfx852/halo/releases/download/v0.9.0/halo-v0.9.0-linux-arm64.tar.gz"
      sha256 "REPLACE_WITH_LINUX_ARM64_SHA256"
    else
      url "https://github.com/zfxzfxzfx852/halo/releases/download/v0.9.0/halo-v0.9.0-linux-amd64.tar.gz"
      sha256 "REPLACE_WITH_LINUX_AMD64_SHA256"
    end
  end

  def install
    bin.install "halo"
  end

  test do
    assert_predicate bin/"halo", :exist?
  end
end
