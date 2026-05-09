class Komyzi < Formula
  desc "AI Agent Configuration Manager - Save, copy and manage configurations for your AI agents"
  homepage "https://github.com/komyzi/komyzi"
  version "0.2.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/komyzi/komyzi/releases/download/v0.2.0/komyzi-darwin-amd64"
      sha256 "TODO: add sha256 after first release"
    end

    if Hardware::CPU.arm?
      url "https://github.com/komyzi/komyzi/releases/download/v0.2.0/komyzi-darwin-arm64"
      sha256 "TODO: add sha256 after first release"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/komyzi/komyzi/releases/download/v0.2.0/komyzi-linux-amd64"
      sha256 "TODO: add sha256 after first release"
    end

    if Hardware::CPU.arm64?
      url "https://github.com/komyzi/komyzi/releases/download/v0.2.0/komyzi-linux-arm64"
      sha256 "TODO: add sha256 after first release"
    end
  end

  def install
    bin.install buildpath/resolved_path
  end

  test do
    system "#{bin}/komyzi", "--version"
  end
end
